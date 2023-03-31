/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-16 11:59:10
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-19 19:56:24
 * @FilePath: /toSearchEngine/internal/repo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package worker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const (
	timeFormart = "2006-01-02 15:04:05"
)

// use kafka or rabbitmq
type Queue interface {
	Consume() (<-chan amqp091.Delivery, error)
}

type MetaRepo interface {
	// map[videoId]Mission
	GetVideoMeta(c context.Context, missions map[string]*Mission) ([]VideoMeta, error)
}

type SearchRepo interface {
	InsertMultiVideoMeta(c context.Context, info []VideoMeta) error
}

type Mission struct {
	VideoId  string    `json:"video_id"`
	UserId   string    `json:"user_id"`
	Time     time.Time `json:"time"`
	Delivery amqp091.Delivery
	DBExist  bool
}

type VideoMeta struct {
	M_id primitive.ObjectID `json:"-" bson:"_id"`
	Id   string             `json:"video_id" `

	UserId     string          `json:"user_id" bson:"userId"`
	UserName   string          `json:"user_name" bson:"userName"`
	Title      string          `json:"title" bson:"title"`
	Duration   string          `json:"duration" bson:"duration"`
	UpdateTime VideoUpdateTime `json:"uploadTime" bson:"uploadTime"`
	Permission int8            `json:"permission" bson:"permission"`
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
