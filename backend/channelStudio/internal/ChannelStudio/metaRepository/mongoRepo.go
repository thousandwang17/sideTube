/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 16:07:51
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-30 16:47:09
 * @FilePath: /ChannelStudio/internal/ChannelStudio/metaRepository/mongoRepo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package metaRepository

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"sideTube/ChannelStudio/internal/ChannelStudio"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongolRepo struct {
	client   *mongo.Client
	videoCol *mongo.Collection
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
	InvalidID   = errors.New("invalid video id")
	ErrDatabase = errors.New("dataBase error!!")
)

type mongoErr struct {
	E error
}

func (m mongoErr) Error() string {
	return m.E.Error()
}

func (m mongoErr) StatusCode() int {
	return http.StatusBadRequest
}

func NewMongoRepo(db *mongo.Client) ChannelStudio.MetaRepository {
	return &mongolRepo{
		client:   db,
		videoCol: db.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_VIDEO_COLLECTION")),
	}
}

type VideoID struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

func (m mongolRepo) UpdateVideoMeta(ctx context.Context, userId string, data ChannelStudio.VideoEditMeta, pictureName string) error {

	objectId, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		log.Println("Invalid video id", data.Id)
		return InvalidID
	}

	setdata := bson.D{
		{Key: "title", Value: data.Title},
		{Key: "desc", Value: data.Description},
		{Key: "lastUpdatetime", Value: primitive.Timestamp{T: uint32(time.Now().Unix())}},
	}

	// if pictureName upload , update png address
	if pictureName != "" {
		setdata = append(setdata, bson.E{Key: "png", Value: pictureName})
	}

	res := m.videoCol.FindOneAndUpdate(
		ctx,
		bson.D{
			{Key: "_id", Value: objectId},
			{Key: "userId", Value: userId},
		},
		bson.D{
			{Key: "$set",
				Value: setdata,
			},
		})
	return res.Err()
}

func (m mongolRepo) GetVideoById(ctx context.Context, videoId, userId string) (ChannelStudio.VideoMeta, error) {
	objectId, err := primitive.ObjectIDFromHex(videoId)
	videoMeta := ChannelStudio.VideoMeta{}
	if err != nil {
		log.Println("Invalid video id", videoId)
		return videoMeta, InvalidID
	}

	res := m.videoCol.FindOne(
		ctx,
		bson.D{
			{Key: "_id", Value: objectId},
			{Key: "userId", Value: userId},
		})

	if res.Err() != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return videoMeta, mongoErr{errors.New("video is not exist")}
		}
		return videoMeta, ErrDatabase
	}

	if err = res.Decode(&videoMeta); err != nil {
		log.Println("GetVideoById Decode", err)
		return videoMeta, ErrDatabase
	}

	return videoMeta, nil
}

func (m mongolRepo) UpdateVideoState(ctx context.Context, videoId, userId string, state int8) error {

	objectId, err := primitive.ObjectIDFromHex(videoId)
	if err != nil {
		log.Println("Invalid video id", videoId)
		return InvalidID
	}

	filter := bson.D{
		{Key: "_id", Value: objectId},
		{Key: "userId", Value: userId},
		{Key: "title", Value: bson.M{"$exists": true}},
		{Key: "desc", Value: bson.M{"$exists": true}},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "permission", Value: state},
			{Key: "lastUpdatetime", Value: primitive.Timestamp{T: uint32(time.Now().Unix())}},
		}},
	}

	opts := options.Update().SetUpsert(false)
	result, err := m.videoCol.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		log.Println(videoId, userId, err)
		return ErrDatabase
	}

	if result.MatchedCount == 0 {
		return mongoErr{errors.New("title and desc are empty")}
	}
	return nil
}

func (m mongolRepo) GetVideoList(ctx context.Context, userId string, skip int64, length int64, public bool) ([]ChannelStudio.VideoMeta, error) {

	filter := bson.D{{Key: "userId", Value: userId}}
	if public {
		filter = append(filter, bson.E{Key: "permission", Value: 1})
	}
	sort := bson.D{{Key: "uploadTime", Value: -1}}

	opts := options.Find().SetSort(sort).SetLimit(length).SetSkip(skip)

	cursor, err := m.videoCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var results []ChannelStudio.VideoMeta
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	for i := range results {
		results[i].Id = results[i].M_id.Hex()
	}

	return results, nil
}

func (m mongolRepo) GetPublicVideoCount(ctx context.Context, userId string) (count int64, err error) {

	filter := bson.D{{Key: "userId", Value: userId}, {Key: "mpd", Value: bson.M{"$exists": true}}, {Key: "permission", Value: 1}}
	count, err = m.videoCol.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m mongolRepo) GetVideoCount(ctx context.Context, userId string) (count int64, err error) {

	filter := bson.D{{Key: "userId", Value: userId}}

	count, err = m.videoCol.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}
