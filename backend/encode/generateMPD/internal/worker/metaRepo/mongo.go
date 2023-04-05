/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-22 20:59:24
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-05 16:48:16
 * @FilePath: /generateMPD/internal/worker/metaRepo/mongo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package metaRepo

import (
	"context"
	"errors"
	"generateMPD/internal/worker"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongolRepo struct {
	db         *mongo.Client
	database   string
	collection string
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

func NewMongoRepo(db *mongo.Client) worker.MetaRepo {
	return &mongolRepo{
		db:         db,
		database:   os.Getenv("MONGO_DATABASE"),
		collection: os.Getenv("MONGO_COLLECTION"),
	}
}

type VideoID struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

// update video status to successed ,
func (m mongolRepo) UpdateStateAndFileNames(ctx context.Context, mission worker.Mission, mpdFileName, pngFileName, duration string) error {

	objectId, err := primitive.ObjectIDFromHex(mission.VideoId)
	if err != nil {
		log.Println("Invalid video id", mission.VideoId)
		return InvalidID
	}

	filter := bson.D{
		{Key: "_id", Value: objectId},
		{Key: "userId", Value: mission.UserId},
	}

	update := bson.A{
		bson.M{"$set": bson.M{
			"png": bson.M{
				"$ifNull": bson.A{"$png", pngFileName},
			},
			"mpd":      mpdFileName,
			"duration": duration,
		},
		},
	}

	res := m.db.Database(m.database).Collection(m.collection).FindOneAndUpdate(
		ctx,
		filter,
		update)

	return res.Err()
}
