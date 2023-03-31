/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 16:07:51
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-19 17:23:54
 * @FilePath: /ChannelStudio/internal/ChannelStudio/metaRepository/mongoRepo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package metaRepository

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sideTube/search/internal/search"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongolRepo struct {
	client   *mongo.Client
	videoCol *mongo.Collection
	// count of public video
	channelVideoCountVol *mongo.Collection
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

func NewMongoRepo(db *mongo.Client) search.MetaRepository {
	return &mongolRepo{
		client:   db,
		videoCol: db.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_VIDEO_COLLECTION")),
	}
}

type VideoID struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

func (m mongolRepo) GetPublicVideosByIds(ctx context.Context, ids []string) ([]search.VideoMeta, error) {

	objectIDs := make([]primitive.ObjectID, len(ids))
	for i := range ids {
		objectID, err := primitive.ObjectIDFromHex(ids[i])
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, objectID)
	}

	filter := bson.D{{Key: "_id", Value: bson.M{"$in": objectIDs}}}
	filter = append(filter, bson.E{Key: "permission", Value: 1})

	sort := bson.D{{Key: "uploadTime", Value: -1}}

	opts := options.Find().SetSort(sort)

	cursor, err := m.videoCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var results []search.VideoMeta
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	for i := range results {
		results[i].Id = results[i].M_id.Hex()
	}

	return results, nil
}
