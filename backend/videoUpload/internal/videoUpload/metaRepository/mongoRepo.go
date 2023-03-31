/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 16:07:51
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-02 21:05:06
 * @FilePath: /videoUpload/internal/videoUpload/metaRepository/mongoRepo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package metaRepository

import (
	"context"
	"errors"
	"log"
	"os"
	"sideTube/videoUpload/internal/videoUpload"
	"sideTube/videoUpload/internal/videoUpload/metaRepository/cacheRepository"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongolRepo struct {
	cache      cacheRepository.Repository
	db         *mongo.Client
	database   string
	collection string
}

type meta struct {
	State int64 `bson:"state"`
}

// video state code
const (
	Uploading = iota
	Encoding
	UnPublish
	Publish
)

var (
	InvalidID           = errors.New("invalid video id")
	NoDocumentFound     = errors.New("No update of video found")
	FailedGetVideoState = errors.New("Failed to get state form video mate ")
)

func NewMongoRepo(db *mongo.Client) videoUpload.MetaRepository {
	return &mongolRepo{
		db:         db,
		database:   os.Getenv("MONGO_DATABASE"),
		collection: os.Getenv("MONGO_COLLECTION"),
	}
}

type VideoID struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

// insert a new video meta , setting expireAt for ttl
func (m mongolRepo) Insert(ctx context.Context, userId, userName string) (InsertedID string, err error) {

	rep, err := m.db.Database(m.database).Collection(m.collection).InsertOne(
		ctx,
		bson.D{
			{Key: "userId", Value: userId},
			{Key: "userName", Value: userName},
			{Key: "expireAt", Value: time.Now().UnixNano() / 1e6},
		},
	)

	if err != nil {
		return "", err
	}

	if videoId, ok := rep.InsertedID.(primitive.ObjectID); ok {
		return videoId.Hex(), nil
	}

	return "", errors.New("videoId is empty")

}

// this is standalone deployment, if want to add transation, then using mongo in replica set and sharded clusters
// update state of video meta to null => Encoding and remove expireAt (ttl)
func (m mongolRepo) UpdateState(ctx context.Context, videoId, userId string) error {

	objectId, err := primitive.ObjectIDFromHex(videoId)
	if err != nil {
		log.Println("Invalid video id", videoId)
		return InvalidID
	}

	// Check if the document exists and has the expected state.
	res := m.db.Database(m.database).Collection(m.collection).FindOne(ctx, bson.D{
		{Key: "_id", Value: objectId},
		{Key: "userId", Value: userId},
	})

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			log.Println("No document found", videoId)
			return NoDocumentFound
		}
		return res.Err()
	}

	var document meta
	if err := res.Decode(&document); err != nil {
		log.Println("UpdateState decode error : ", err)
		return FailedGetVideoState
	}

	// if State not set, document.State will be 0
	// if State already is Encoding, then return succeed (err = nil)
	if document.State == Encoding {
		return nil
	}

	res = m.db.Database(m.database).Collection(m.collection).FindOneAndUpdate(
		ctx,
		bson.D{
			{Key: "_id", Value: objectId},
			{Key: "userId", Value: userId},
		},
		bson.D{
			{Key: "$unset",
				Value: bson.D{
					{Key: "expireAt", Value: ""},
				},
			},
			{Key: "$set",
				Value: bson.D{
					{Key: "state", Value: Encoding},
					{Key: "uploadTime", Value: primitive.Timestamp{T: uint32(time.Now().Unix())}},
				},
			},
		})
	return res.Err()
}

// this is standalone deployment, if want to add transation, then using mongo in replica set and sharded clusters
// undo  state of video meta to Encoding  => null and add expireAt (ttl) back
func (m mongolRepo) UndoUpdateState(ctx context.Context, videoId, userId string) error {

	objectId, err := primitive.ObjectIDFromHex(videoId)
	if err != nil {
		log.Println("Invalid video id", videoId)
		return InvalidID
	}

	// Check if the document exists and has the expected state.
	res := m.db.Database(m.database).Collection(m.collection).FindOne(ctx, bson.D{
		{Key: "_id", Value: objectId},
		{Key: "userId", Value: userId},
		{Key: "state", Value: Encoding},
	})

	if res.Err() != nil {
		// ErrNoDocuments only happend on state is not Encoding that mean state still is Uploading (not exists) or video been Encoded
		// or Documents alerdy lost that we dont need to update
		// therefore, we return successed ( error = nil )
		if res.Err() == mongo.ErrNoDocuments {
			return nil
		}
		return res.Err()
	}

	res = m.db.Database(m.database).Collection(m.collection).FindOneAndUpdate(
		ctx,
		bson.D{
			{Key: "_id", Value: objectId},
			{Key: "userId", Value: userId},
		},
		bson.D{
			{Key: "$unset",
				Value: bson.D{
					{Key: "state", Value: ""},
				},
			},
			{Key: "$set",
				Value: bson.D{
					{Key: "expireAt", Value: time.Now().UnixNano() / 1e6},
					{Key: "uploadTime", Value: primitive.Timestamp{T: uint32(time.Now().Unix())}},
				},
			},
		})
	return res.Err()
}

//  remove video
func (m mongolRepo) Remove(ctx context.Context, videoId, userId string) error {

	objectId, err := primitive.ObjectIDFromHex(videoId)
	if err != nil {
		log.Println("Invalid video id", videoId)
		return InvalidID
	}

	res := m.db.Database(m.database).Collection(m.collection).FindOneAndDelete(
		ctx,
		bson.D{
			{Key: "_id", Value: objectId},
			{Key: "userId", Value: userId},
		},
	)

	if res.Err() != mongo.ErrNoDocuments || res.Err() == nil {
		return nil
	}

	return res.Err()
}
