/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 15:46:39
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-13 19:40:28
 * @FilePath: /videoMessage/internal/videoMessage/video.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package videoMessage

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

type VideoMessageMeta struct {
	Message    string             `json:"message,omitempty" bson:"message"`
	UserId     string             `json:"user_id" bson:"userId"`
	UserName   string             `json:"user_name" bson:"userName"`
	VideoID    string             `json:"video_id" bson:"videoId"`
	CreateTime MessageTime        `json:"create_time" bson:"createTime"`
	UpdateTime MessageTime        `json:"update_time,omitempty" bson:"lastUpdatetime"`
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Replies    int64              `json:"replies,omitempty" bson:"Replies"`
}

type VideoMessageReplyMeta struct {
	MessageID  string             `json:"message_id" bson:"messageId"`
	Message    string             `json:"message,omitempty" bson:"message"`
	UserId     string             `json:"user_id" bson:"userId"`
	UserName   string             `json:"user_name" bson:"userName"`
	CreateTime MessageTime        `json:"time" bson:"createTime"`
	UpdateTime MessageTime        `json:"update_time,omitempty" bson:"lastUpdatetime"`
	ID         primitive.ObjectID `json:"id" bson:"_id"`
}

type MessageTime struct {
	T time.Time
}

func (v MessageTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(v.T).Format(timeFormart))
	return []byte(stamp), nil
}

func (v *MessageTime) UnmarshalBSONValue(b bsontype.Type, value []byte) error {
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
