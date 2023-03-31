/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-16 12:09:49
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-13 14:39:38
 * @FilePath: /detectVideo/internal/worker/handler/handler.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package handler

import (
	"context"
	"detectVideo/internal/worker"
	"encoding/json"
	"fmt"
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
	notifyClose chan bool
	shutdown    chan bool
	state       state
	wg          *sync.WaitGroup
}

func New(v worker.VideoRepository, lock worker.LockSystem, queue worker.Queue) worker.Handler {
	return &handle{
		videoRepo:   v,
		lockRepo:    lock,
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
	if err = h.lockRepo.Lock(ctx, mission.VideoId, TimeOut); err != nil {
		// if connect lost or key alerdy exist, just re-queue the message
		log.Println("lock err : ", err)

		if err := h.queueRepo.ReQueue(ctx, b, delay); err != nil {
			d.Nack(false, true)
		}
		AckTrue(d)
		return
	}

	defer h.lockRepo.UnLock(ctx, mission.VideoId)

	// set id for each encode mission of the video
	missionID := 0

	// trasnform each support format to next encoding mission
	var encodeVideoMissions [][]byte
	var encodeAudioMissions [][]byte
	if vs, as, vtc, err := h.videoRepo.DetectVideo(ctx, mission.VideoId); err != nil {
		log.Print(err)
		if err := h.queueRepo.ReQueue(ctx, b, delay); err != nil {
			log.Println("ReQueue", err)
			// if err, let message retry
			d.Nack(false, true)
		}
		AckTrue(d)
		return
	} else {
		log.Println(len(vs), vtc)
		encodeVideoMissions = make([][]byte, 0, len(vs)*vtc)

		// videos missions
		for v := range vs {
			// encode

			for chunckNum := 0; chunckNum < vtc; chunckNum++ {
				missionByte, err := json.Marshal(worker.EncodeVideoMission{
					UserId:       mission.UserId,
					VideoId:      mission.VideoId,
					Time:         time.Now(),
					VideoFormat:  vs[v],
					MissionID:    missionID,
					SubMissionID: chunckNum,
					TotalChunk:   vtc,
					MissionType:  worker.VideoType,
				})

				if err != nil {
					log.Println("json transform failed ", mission, err)
					d.Nack(false, false)
					return
				}

				encodeVideoMissions = append(encodeVideoMissions,
					missionByte,
				)
				log.Printf("%s-%d %d => %v", mission.VideoId, missionID, chunckNum, vs[v])
			}

			missionID++
		}

		encodeAudioMissions = make([][]byte, 0, len(as))

		// aduio missions
		for a := range as {
			// encode
			missionByte, err := json.Marshal(worker.EncodeVideoMission{
				UserId:      mission.UserId,
				VideoId:     mission.VideoId,
				Time:        time.Now(),
				AudioFormat: as[a],
				MissionID:   missionID,
				MissionType: worker.AudioType,
			})

			if err != nil {
				log.Println("json transform failed ", mission, err)
				d.Nack(false, false)
				return
			}

			encodeAudioMissions = append(encodeAudioMissions,
				missionByte,
			)

			log.Printf("%s %d => %v", mission.VideoId, missionID, as[a])
			missionID++
		}
	}

	if err = h.lockRepo.SetMissionMap(ctx, mission.VideoId, missionID-1); err != nil {
		log.Print(err)
		if err := h.queueRepo.ReQueue(ctx, b, delay); err != nil {
			// if err, let message retry
			d.Nack(false, true)
		}
		AckTrue(d)
		return
	}

	// publish encoding mission
	if err = h.queueRepo.PublishEncode(ctx, encodeVideoMissions, encodeAudioMissions); err != nil {
		log.Print(err)
		if err := h.queueRepo.ReQueue(ctx, b, delay); err != nil {
			// if err, let message retry
			d.Nack(false, true)
		}
		AckTrue(d)
		return
	}
	AckTrue(d)
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
