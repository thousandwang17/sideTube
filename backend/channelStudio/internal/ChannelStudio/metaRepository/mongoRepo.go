/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 16:07:51
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-01-30 17:23:13
 * @FilePath: /ChannelStudio/internal/ChannelStudio/metaRepository/mongoRepo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package metaRepository

import (
	"context"
	"errors"
	"log"
	"os"
	"sideTube/ChannelStudio/internal/ChannelStudio"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func NewMongoRepo(db *mongo.Client) ChannelStudio.MetaRepository {
	return &mongolRepo{
		db.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION")),
	}
}

type VideoID struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

func (m mongolRepo) UpdateVideoMeta(ctx context.Context, userId string, data ChannelStudio.VideoEditMeta) error {

	objectId, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		log.Println("Invalid video id", data.Id)
		return InvalidID
	}

	res := m.col.FindOneAndUpdate(
		ctx,
		bson.D{
			{Key: "_id", Value: objectId},
			{Key: "userId", Value: userId},
		},
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "title", Value: data.Title},
					{Key: "desc", Value: data.Description},
					{Key: "lastUpdatetime", Value: primitive.Timestamp{T: uint32(time.Now().Unix())}},
				},
			},
		})
	return res.Err()
}

func (m mongolRepo) UpdateVideoState(ctx context.Context, videoId, userId string, state uint8) error {

	objectId, err := primitive.ObjectIDFromHex(videoId)
	if err != nil {
		log.Println("Invalid video id", videoId)
		return InvalidID
	}

	res := m.col.FindOneAndUpdate(
		ctx,
		bson.D{
			{Key: "_id", Value: objectId},
			{Key: "userId", Value: userId},
		},
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "permission", Value: state},
					{Key: "lastUpdatetime", Value: primitive.Timestamp{T: uint32(time.Now().Unix())}},
				},
			},
		})
	return res.Err()
}

func (m mongolRepo) GetVideoList(ctx context.Context, userId string, skip int64, length int64) ([]ChannelStudio.VideoMeta, error) {

	filter := bson.D{{Key: "userId", Value: userId}}
	sort := bson.D{{Key: "uploadTime", Value: -1}}

	opts := options.Find().SetSort(sort).SetLimit(length).SetSkip(skip)

	cursor, err := m.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var results []ChannelStudio.VideoMeta
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	for i := range results {
		results[i].Id = results[i].M_id.Hex()
		log.Println(results[i])
	}

	return results, nil
}

func (m mongolRepo) GetVideoCount(ctx context.Context, userId string) (count int64, err error) {

	filter := bson.D{{Key: "userId", Value: userId}}

	count, err = m.col.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}
