/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-16 11:59:10
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-18 22:59:21
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
	CompleteMultipartUpload(c context.Context, videoID string) error
}

// in produntion project, suggest to use etcd for distribute locking system
type LockSystem interface {
	Lock(c context.Context, videoID string, ttl time.Duration) error
	UnLock(c context.Context, videoID string) error
}

// use kafka or rabbitmq
type Queue interface {
	Consume() (<-chan amqp091.Delivery, error)
	// delay : 1000 = 1s
	ReQueue(ctx context.Context, body []byte, delay int64) error
	PublishDectedVideo(ctx context.Context, body []byte) error
}

type Mission struct {
	VideoId string    `json:"video_id"`
	UserId  string    `json:"user_id"`
	Time    time.Time `json:"time"`
	Retries int       `json:"retries"`
}

type DecteddMission struct {
	VideoId string    `json:"video_id"`
	UserId  string    `json:"user_id"`
	Time    time.Time `json:"time"`
}
