/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-16 11:59:10
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-16 20:38:34
 * @FilePath: /generateMPD/internal/repo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package worker

import (
	"context"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type VideoRepository interface {
	// generatemdp  mpd file and min png for show
	GenerateMPD(c context.Context, info Mission, videoList, audioList []string) (mpdFileName, pngfileName, duration string, err error)
}

type MetaRepo interface {
	UpdateStateAndFileNames(c context.Context, mission Mission, mpdFileName, pngfileName, duration string) error
}

// in produntion project, suggest to use etcd for distribute locking system
type LockSystem interface {
	Lock(c context.Context, videoID string, ttl time.Duration) (videoList, audioList []string, err error)
	UnLock(c context.Context, videoID string) error
	ReleaseMissionKey(c context.Context, videoID string) error
}

// use kafka or rabbitmq
type Queue interface {
	Consume() (<-chan amqp091.Delivery, error)
	// delay : 1000 = 1s
	ReQueue(ctx context.Context, body []byte, delay int64) error
	PublishSearchEngine(ctx context.Context, body []byte) error
}

type Mission struct {
	VideoId    string    `json:"video_id"`
	UserId     string    `json:"user_id"`
	Time       time.Time `json:"time"`
	TotalChunk int       `json:"total_chunk"`
	Retries    int       `json:"retries"`
}

type VideoList struct {
	Path string
}

type AudioList struct {
	Path string
}
