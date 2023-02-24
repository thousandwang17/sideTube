/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 15:46:39
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-01-30 17:35:11
 * @FilePath: /ChannelStudio/internal/ChannelStudio/video.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package ChannelStudio

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const (
	timeFormart = "2006-01-02 15:04:05"
)

type VideoEditMeta struct {
	Id          string `json:"video_id" `
	Title       string `json:"title,omitempty" bson:"title"`
	Description string `json:"desc,omitempty" bson:"desc"`
}

type VideoMeta struct {
	Id          string             `json:"video_id" `
	Title       string             `json:"title,omitempty" bson:"title"`
	Description string             `json:"desc,omitempty" bson:"desc"`
	Like        uint64             `json:"like,omitempty" bson:"like"`
	DisLike     uint64             `json:"disLike,omitempty" bson:"disLike"`
	Views       uint64             `json:"views,omitempty" bson:"views"`
	Messages    uint64             `json:"messages,omitempty" bson:"messages"`
	CreateTime  time.Time          `json:"createTime,omitempty" bson:"createTime"`
	UpdateTime  VideoUpdateTime    `json:"uploadTime" bson:"uploadTime"`
	M_id        primitive.ObjectID `json:"-" bson:"_id"`
	V_State     int8               `json:"state,omitempty" bson:"state"`
	Permission  int8               `json:"permission,omitempty" bson:"permission"`
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
