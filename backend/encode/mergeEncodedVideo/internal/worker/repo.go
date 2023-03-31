/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-16 11:59:10
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-11 18:07:52
 * @FilePath: /mergeUploadfile/internal/repo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package worker

import (
	"context"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type VideoRepository interface {
	MergeEncodeVideo(c context.Context, info Mission) (encodeFilName string, err error)
}

// in produntion project, suggest to use etcd for distribute locking system
type LockSystem interface {
	Lock(c context.Context, videoID string, ttl time.Duration) error
	UnLock(c context.Context, videoID string) error
	// allDone mean all missioon is completed
	// encodedFileName will save to list for generate mpd file
	AccomplishbMission(c context.Context, mission Mission, encodedFileName string) (allDone bool, err error)
}

// use kafka or rabbitmq
type Queue interface {
	Consume() (<-chan amqp091.Delivery, error)
	// delay : 1000 = 1s
	ReQueue(ctx context.Context, body []byte, delay int64) error
	PublishGenerateMpd(ctx context.Context, missions []byte) error
}

type Mission struct {
	VideoId    string    `json:"video_id"`
	UserId     string    `json:"user_id"`
	Fps        float64   `json:"fps"`
	Time       time.Time `json:"time"`
	Height     int       `json:"height"`
	TotalChunk int       `json:"total_chunk"`
	Retries    int       `json:"retries"`
	MissionID  int       `json:"mission_id"`
}

type GenerateMpdMission struct {
	VideoId    string    `json:"video_id"`
	UserId     string    `json:"user_id"`
	TotalChunk int       `json:"total_chunk"`
	Time       time.Time `json:"time"`
}
