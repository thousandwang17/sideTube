/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-16 12:09:49
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-16 20:39:37
 * @FilePath: /generateMPD/internal/worker/handler/handler.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"generateMPD/internal/worker"
	"log"
	"os"
	"strings"
	"sync"
	"time"

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
	videoRepo   worker.VideoRepository
	lockRepo    worker.LockSystem
	queueRepo   worker.Queue
	metaRepo    worker.MetaRepo
	notifyClose chan bool
	shutdown    chan bool
	state       state
	wg          *sync.WaitGroup
}

func New(v worker.VideoRepository, lock worker.LockSystem, queue worker.Queue, meta worker.MetaRepo) worker.Handler {
	return &handle{
		videoRepo:   v,
		lockRepo:    lock,
		queueRepo:   queue,
		metaRepo:    meta,
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

L:
	for {
		select {
		// nums of routine control by channel.Qos of rabbitmq
		case d, ok := <-msgs:
			if ok {
				h.wg.Add(1)
				go h.helper(d)
			} else {
				break L
			}
		case <-h.notifyClose:
			break L
		}
	}

	down := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut+5*time.Second)
	defer cancel()

	go func() {
		h.wg.Wait()
		down <- struct{}{}
	}()

	select {
	case <-down:
		log.Println("mission done")
		select {
		case _, _ = <-h.notifyClose:
			// will keep go on and  shut down the server
		default:
			// if rabbitmq connect closed, it will get there.
			os.Exit(0)

			return
		}
	case <-ctx.Done():
		log.Println("Warring : context timeout ")
	}

	close(h.shutdown)
}

func (h *handle) helper(d amqp091.Delivery) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	defer h.wg.Done()
	log.Printf("Received a message: %s", d.Body)

	// decode
	// Create a new decoder
	decoder := json.NewDecoder(strings.NewReader(string(d.Body)))

	var mission worker.Mission
	err := decoder.Decode(&mission)
	if err != nil {
		log.Println("json transform failed ", string(d.Body))
		// save to db for handling or ignored this message
		AckTrue(d)
		return
	}

	// vaild count of Retries
	if mission.Retries >= 3 {
		// if reach max limit, send to dead-letter queue.
		// this side project will not implement this part
		log.Print(err)
		AckTrue(d)
		return
	}
	mission.Retries += 1

	// encode
	b, err := json.Marshal(mission)
	if err != nil {
		log.Println("json transform failed ", mission, err)
		d.Nack(false, false)
		return
	}

	// lock this message by videoID
	videoList, audioList, err := h.lockRepo.Lock(ctx, mission.VideoId, TimeOut)
	if err != nil {
		// if connect lost or key alerdy exist, just re-queue the message
		log.Println("lock err : ", err)

		if err := h.queueRepo.ReQueue(ctx, b, delay); err != nil {
			d.Nack(false, true)
		}
		AckTrue(d)
		return
	}

	defer h.lockRepo.UnLock(ctx, mission.VideoId)

	//
	if mpdFileName, pngFileName, duration, err := h.videoRepo.GenerateMPD(ctx, mission, videoList, audioList); err != nil {
		log.Print("GenerateMPD err :", err)
		if err := h.queueRepo.ReQueue(ctx, b, delay); err != nil {
			log.Println("ReQueue", err)
			// if err, let message retry
			d.Nack(false, true)
		}
		AckTrue(d)
		return
	} else {

		// if generate file succesed  , add meta info to metaRepo
		if err := h.metaRepo.UpdateStateAndFileNames(ctx, mission, mpdFileName, pngFileName, duration); err != nil {
			log.Println("metaRepo err: ", err)
			if err := h.queueRepo.ReQueue(ctx, b, delay); err != nil {
				log.Println("ReQueue", err)
				// if err, let message retry
				d.Nack(false, true)
			}
		}

		h.lockRepo.ReleaseMissionKey(ctx, mission.VideoId)

		// merge mission
		missionByte, err := json.Marshal(worker.Mission{
			VideoId: mission.VideoId,
			UserId:  mission.UserId,
			Time:    time.Now(),
		})

		if err != nil {
			log.Println("json transform failed ", mission, err)
			d.Nack(false, false)
			return
		}

		// publish encoding mission
		if err = h.queueRepo.PublishSearchEngine(ctx, missionByte); err != nil {
			log.Print(err)
			if err := h.queueRepo.ReQueue(ctx, b, delay); err != nil {
				// if err, let message retry
				d.Nack(false, true)
			}
			AckTrue(d)
			return
		}

		log.Println(mission.VideoId, "Down!!")
		AckTrue(d)
	}
}

func (h *handle) Shutdown() {
	if h.state == Handling {
		h.state = Closed

		close(h.notifyClose)
		<-h.shutdown
	}
}

func AckTrue(d amqp091.Delivery) error {
	if err := d.Ack(false); err != nil {
		fmt.Println("d.Ack err: ", err)
		return err
	}
	return nil
}
