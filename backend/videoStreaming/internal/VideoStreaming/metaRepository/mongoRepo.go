/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 16:07:51
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-08 20:18:08
 * @FilePath: /VideoStreaming/internal/VideoStreaming/metaRepository/mongoRepo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package metaRepository

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"sideTube/VideoStreaming/internal/VideoStreaming"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongolRepo struct {
	client     *mongo.Client
	videoCol   *mongo.Collection
	viewLogCol *mongo.Collection
}

type mongoErr struct {
	E error
}

func (m mongoErr) Error() string {
	return m.E.Error()
}

func (m mongoErr) StatusCode() int {
	return http.StatusBadRequest
}

var (
	InvalidID   = errors.New("invalid video id")
	ErrDatabase = errors.New("dataBase error!!")
)

func NewMongoRepo(db *mongo.Client) VideoStreaming.MetaRepository {
	return &mongolRepo{
		videoCol:   db.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_VIDEO_COLLECTION")),
		viewLogCol: db.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_VIEW_LOGS_COLLECTION")),
	}
}

type VideoID struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

func (m mongolRepo) GetVideoMetaById(ctx context.Context, video_id string) (VideoStreaming.VideoMeta, error) {

	objectId, err := primitive.ObjectIDFromHex(video_id)
	result := VideoStreaming.VideoMeta{}

	if err != nil {
		return VideoStreaming.VideoMeta{}, mongoErr{InvalidID}
	}

	filter := bson.D{{Key: "_id", Value: objectId}, {Key: "mpd", Value: bson.M{"$exists": true}}}

	err = m.videoCol.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, mongoErr{errors.New("video is not exist")}
		}

		return result, ErrDatabase
	}

	return result, nil
}

func (m mongolRepo) IncVideoViews(ctx context.Context, video_id, user_id string) error {

	objectId, err := primitive.ObjectIDFromHex(video_id)
	if err != nil {
		return InvalidID
	}

	rep := m.videoCol.FindOneAndUpdate(ctx,
		bson.D{{Key: "_id", Value: objectId}, {Key: "permission", Value: 1}},
		bson.D{
			{Key: "$inc",
				Value: bson.D{
					{Key: "views", Value: 1},
				},
			},
		})

	if rep.Err() != nil {
		if errors.Is(rep.Err(), mongo.ErrNoDocuments) {
			return mongoErr{errors.New("videoid is not exist")}
		}

		log.Println("IncVideoViews err: ", rep.Err())
		return ErrDatabase
	}

	// Use Upsert() to update the document if it exists or insert a new one if it doesn't
	_, err = m.viewLogCol.UpdateOne(context.Background(),
		bson.D{
			{Key: "user_id", Value: user_id},
			{Key: "video_id", Value: video_id},
		},
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "date", Value: primitive.Timestamp{T: uint32(time.Now().Unix())}},
				},
			},
		},
		options.Update().SetUpsert(true))
	if err != nil {
		return ErrDatabase
	}

	return nil
}
