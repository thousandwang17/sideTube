/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-16 12:09:49
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-10 10:26:42
 * @FilePath: /toSearchEngine/internal/worker/handler/handler.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
	"toSearchEngine/internal/worker"
	"toSearchEngine/internal/worker/searchRepo"

	"github.com/rabbitmq/amqp091-go"
)

const (
	NotStart state = iota
	Handling
	Closed

	TimeOut = time.Second * 60 * 5
	delay   = 1000 * 30 // 30s
)

type state int

type handle struct {
	metaRepo    worker.MetaRepo
	searchRepo  worker.SearchRepo
	queueRepo   worker.Queue
	notifyClose chan bool
	shutdown    chan bool
	state       state
	wg          *sync.WaitGroup
}

func New(meta worker.MetaRepo, search worker.SearchRepo, queue worker.Queue) worker.Handler {
	return &handle{
		metaRepo:    meta,
		searchRepo:  search,
		queueRepo:   queue,
		notifyClose: make(chan bool),
		shutdown:    make(chan bool),
		wg:          &sync.WaitGroup{},
	}
}

// if detect distribute locking system locking the message, we re-queue to origin queue with "x-dely" and nack this message
// be sure re-queue before nack , and if this message count of retries that in message body reach max limit ,
// we deliver this message to dead-letter for further processing.
func (h *handle) Handle() {
	msgs, err := h.queueRepo.Consume()

	if err != nil {
		log.Panicln("failed to consume from queue err : ", err)
	}

	if h.state == NotStart {
		h.state = Handling
	} else {
		log.Println("hanlder is running or closed", h.state)
		return
	}

	var batch []amqp091.Delivery
	batchCount := 10
	batchTimeout := 5 * time.Second
	timer := time.NewTimer(batchTimeout)
	defer timer.Stop()
L:
	for {
		select {
		case <-h.notifyClose:
			break L
		default:
		}

		select {
		// nums of routine control by channel.Qos of rabbitmq
		case d, ok := <-msgs:
			if ok {
				batch = append(batch, d)
				if len(batch) == batchCount {
					// Batch size has been reached, call the batch handler
					h.helper(batch)
					batch = batch[:0]
					timer.Reset(batchTimeout)
				}
			} else {
				break L
			}

		case <-timer.C:
			// Timeout has elapsed, call the batch handler with the current batch
			if len(batch) > 0 {
				h.helper(batch)
				batch = batch[:0]
			}
			timer.Reset(batchTimeout)
		case <-h.notifyClose:
			break L
		}

	}
	close(h.shutdown)
	os.Exit(0)
}

func (h *handle) helper(d []amqp091.Delivery) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	// decode
	// Create a new decoder
	missions := map[string]*worker.Mission{}
	for i := range d {
		var mission worker.Mission
		decoder := json.NewDecoder(strings.NewReader(string(d[i].Body)))
		err := decoder.Decode(&mission)
		if err != nil {
			log.Println("json transform failed ", string(d[i].Body))
			d[i].Nack(false, false)
			continue
		}
		mission.Delivery = d[i]
		missions[mission.VideoId] = &mission
		log.Println("Recived message ", string(d[i].Body))
	}

	// GetVideoMeta  will edit DBExist of mission to check videoId is exist
	info, err := h.metaRepo.GetVideoMeta(ctx, missions)
	if err != nil {
		log.Print(err)
		return
	}

	fail := 0
	for i := range missions {
		if !missions[i].DBExist {
			missions[i].Delivery.Nack(false, false)
			fail++
		}
	}

	if err = h.searchRepo.InsertMultiVideoMeta(ctx, info); err != nil {
		log.Println(err)
		if errors.Is(err, searchRepo.ErrUpsertVideo) {
			d[len(d)-1].Nack(true, false)
			log.Printf("%d messages false ", len(d))
		}
		return
	}

	AckTrue(d[len(d)-1])
	log.Printf("%d messages successed , %d faild", len(d)-fail, fail)
}

func (h *handle) Shutdown() {
	if h.state == Handling {
		h.state = Closed
		close(h.notifyClose)
		<-h.shutdown
	}
}

func AckTrue(d amqp091.Delivery) error {
	if err := d.Ack(true); err != nil {
		fmt.Println("d.Ack err: ", err)
		return err
	}
	return nil
}
