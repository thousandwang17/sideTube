/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 15:46:39
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-23 13:01:40
 * @FilePath: /VideoStreaming/internal/VideoStreaming/video.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package VideoStreaming

import (
	"fmt"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const (
	timeFormart = "2006-01-02 15:04:05"
	TypeMpd     = "MPD"
	TypeVideo   = "video"
)

type VideoMeta struct {
	Title       string          `json:"title,omitempty" bson:"title"`
	Description string          `json:"desc,omitempty" bson:"desc"`
	Like        uint64          `json:"like,omitempty" bson:"like"`
	DisLike     uint64          `json:"disLike,omitempty" bson:"disLike"`
	Views       uint64          `json:"views,omitempty" bson:"views"`
	Messages    uint64          `json:"messages,omitempty" bson:"messages"`
	UpdateTime  VideoUpdateTime `json:"uploadTime" bson:"uploadTime"`

	// video source list and info
	Duration     uint64      `json:"Duration,omitempty" bson:"duration"`
	VideoSources []VideoType `json:"videoSources,omitempty" bson:"videoSources"`
	AudioSources []AudioType `json:"audioSources,omitempty" bson:"audioSources"`
	Mpd          string      `json:"mpd,omitempty" bson:"mpd"`

	V_State    int8 `json:"state,omitempty" bson:"state"`
	Permission int8 `json:"permission,omitempty" bson:"permission"`
}

type VideoType struct {
	Source_ID string `json:"id" bson:"id"`
	Type      string `json:"type" bson:"type"`
	Codes     string `json:"codes,omitempty" bson:"codes"`
	Size      int64  `json:"size" bson:"size"`
	Width     int16  `json:"width" bson:"width"`
	High      int16  `json:"high" bson:"high"`
	Fps       int16  `json:"fps" bson:"fps"`
}

type AudioType struct {
	Source_ID string `json:"id" bson:"id"`
	Type      string `json:"type" bson:"type"`
	Codes     string `json:"codes,omitempty" bson:"codes"`
	Size      int64  `json:"size" bson:"size"`
	Opus      int16  `json:"opus" bson:"opus"`
}

type VideoUpdateTime struct {
	T time.Time `bson:"uploadTime"`
}

func (v VideoUpdateTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(v.T).Format(timeFormart))
	return []byte(stamp), nil
}

func (v *VideoUpdateTime) UnmarshalBSONValue(b bsontype.Type, value []byte) error {
	if b != bsontype.Timestamp {
		return fmt.Errorf("invalid bson value type '%s'", b.String())
	}
	t, _, _, ok := bsoncore.ReadTimestamp(value)
	if !ok {
		return fmt.Errorf("invalid bson Int64 value")
	}

	v.T = time.Unix(int64(t), 0)
	return nil
}

type Video struct {
	Data io.ReadCloser
}
