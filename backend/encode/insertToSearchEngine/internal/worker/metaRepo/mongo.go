/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-22 20:59:24
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-16 16:35:26
 * @FilePath: /generateMPD/internal/worker/metaRepo/mongo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package metaRepo

import (
	"context"
	"errors"
	"os"
	"toSearchEngine/internal/worker"

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
		collection: os.Getenv("MONGO_VIDEO_COLLECTION"),
	}
}

type VideoID struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

// update video status to successed ,
func (m mongolRepo) GetVideoMeta(c context.Context, missions map[string]*worker.Mission) ([]worker.VideoMeta, error) {

	objectIDs := make([]primitive.ObjectID, len(missions))
	for i := range missions {
		objectID, err := primitive.ObjectIDFromHex(missions[i].VideoId)
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, objectID)
	}
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	cursor, err := m.db.Database(m.database).Collection(m.collection).Find(c, filter)
	if err != nil {
		return nil, err
	}

	var VideoMeta []worker.VideoMeta
	if err := cursor.All(c, &VideoMeta); err != nil {
		return nil, err
	}

	for i := range VideoMeta {
		VideoMeta[i].Id = VideoMeta[i].M_id.Hex()
		missions[VideoMeta[i].Id].DBExist = true
	}

	return VideoMeta, nil
}
