/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 15:46:39
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-03 21:16:03
 * @FilePath: /recommend/internal/recommend/video.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package recommend

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

type VideoMeta struct {
	Id         string             `json:"video_id" `
	UserId     string             `json:"user_id" bson:"userId"`
	UserName   string             `json:"user_name" bson:"userName"`
	Title      string             `json:"title,omitempty" bson:"title"`
	Views      uint64             `json:"views,omitempty" bson:"views"`
	UploadTime time.Time          `json:"uploadTime,omitempty" bson:"uploadTime"`
	Duration   string             `json:"duration,omitempty" bson:"duration"`
	Png        string             `json:"png,omitempty" bson:"png"`
	M_id       primitive.ObjectID `json:"-" bson:"_id"`
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
