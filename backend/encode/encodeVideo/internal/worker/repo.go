/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-16 11:59:10
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-11 17:26:29
 * @FilePath: /encodeVideo/internal/repo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package worker

import (
	"context"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type VideoRepository interface {
	// will return format of video including video and audio
	// if len(vs) or len(as) == 0 , it will return error
	EncodeVideo(c context.Context, mission EncodeVideoMission) (encodedFileName string, err error)
}

// in produntion project, suggest to use etcd for distribute locking system
type LockSystem interface {
	Lock(c context.Context, videoID string, missionID, subMissionID int, ttl time.Duration) (done bool, err error)
	UnLock(c context.Context, videoID string, missionID, subMissionID int) error
	// allDone mean all missioon is completed
	// encodedFileName will save to list for generate mpd file
	AccomplishbSubMission(c context.Context, mission EncodeVideoMission, encodedFileName string) (allDone bool, err error)
}

// use kafka or rabbitmq
type Queue interface {
	Consume() (<-chan amqp091.Delivery, error)
	// delay : 1000 = 1s
	ReQueue(ctx context.Context, body []byte, delay int64) error
	PublishMergeEncodedVideo(ctx context.Context, missions []byte) error
}

type Mission struct {
	VideoId    string    `json:"video_id"`
	UserId     string    `json:"user_id"`
	MissionID  int       `json:"mission_id"`
	Fps        float64   `json:"fps"`
	Time       time.Time `json:"time"`
	Height     int       `json:"height"`
	TotalChunk int       `json:"total_chunk"`
}

const (
	VideoType = 1
	AudioType = 2
)

type EncodeVideoMission struct {
	VideoId      string      `json:"video_id"`
	UserId       string      `json:"user_id"`
	Time         time.Time   `json:"time"`
	VideoFormat  VideoFormat `json:"video_format,optimize"`
	Retries      int         `json:"retries"`
	MissionID    int         `json:"mission_id"`
	MissionType  int         `json:"mission_type"`
	SubMissionID int         `json:"sub_mission_id"`
	TotalChunk   int         `json:"total_chunk"`
}

type VideoFormat struct {
	Fps       float64 `json:"fps"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	OriginFps bool    `json:"OriginFps"`
}

type AudioFormat struct {
	Hz int `json:"hz"`
}
