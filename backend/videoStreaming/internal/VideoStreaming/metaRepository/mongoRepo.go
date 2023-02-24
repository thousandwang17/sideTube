/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 16:07:51
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-03 14:53:49
 * @FilePath: /VideoStreaming/internal/VideoStreaming/metaRepository/mongoRepo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package metaRepository

import (
	"context"
	"errors"
	"os"
	"sideTube/VideoStreaming/internal/VideoStreaming"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongolRepo struct {
	col *mongo.Collection
}

// video state code
const (
	Encoding = iota
)

// video state code
const (
	UnPublish = iota
	Publish
)

var (
	InvalidID = errors.New("invalid video id")
)

func NewMongoRepo(db *mongo.Client) VideoStreaming.MetaRepository {
	return &mongolRepo{
		db.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION")),
	}
}

type VideoID struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

func (m mongolRepo) GetVideoMetaById(ctx context.Context, video_id string) (VideoStreaming.VideoMeta, error) {

	objectId, err := primitive.ObjectIDFromHex(video_id)
	result := VideoStreaming.VideoMeta{}

	if err != nil {
		return VideoStreaming.VideoMeta{}, InvalidID
	}

	filter := bson.D{{Key: "_id", Value: objectId}}

	err = m.col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return VideoStreaming.VideoMeta{}, err
	}

	return result, nil
}
