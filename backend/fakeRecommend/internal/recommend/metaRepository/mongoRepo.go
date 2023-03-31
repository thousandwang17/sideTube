/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 16:07:51
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-03 16:16:04
 * @FilePath: /recommend/internal/recommend/metaRepository/mongoRepo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package metaRepository

import (
	"context"
	"errors"
	"os"
	"sideTube/recommend/internal/recommend"

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

// NOTE : this recommend list just get videos that public and  been encoded order by time
// in production should use Collaborative filtering algorithms to build recommend list ,
// Apache Mahout, Apache Spark and  TensorFlow are good choose to impelemnet recommendation system
func NewMongoRepo(db *mongo.Client) recommend.MetaRepository {
	return &mongolRepo{
		db.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION")),
	}
}

type VideoID struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

func (m mongolRepo) GetHomeVideoList(ctx context.Context, userId string, skip int64, length int64) ([]recommend.VideoMeta, error) {

	filter := bson.D{{Key: "permission", Value: 1}, {Key: "mpd", Value: bson.M{"$exists": true}}}
	sort := bson.D{{Key: "mpd", Value: -1}}

	opts := options.Find().SetSort(sort).SetLimit(length).SetSkip(skip)

	cursor, err := m.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var results []recommend.VideoMeta
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	for i := range results {
		results[i].Id = results[i].M_id.Hex()
	}

	return results, nil
}

func (m mongolRepo) GetRelationVideoList(ctx context.Context, videoID string, skip int64, length int64) ([]recommend.VideoMeta, error) {

	filter := bson.D{{Key: "permission", Value: 1}, {Key: "mpd", Value: bson.M{"$exists": true}}}
	sort := bson.D{{Key: "mpd", Value: -1}}

	opts := options.Find().SetSort(sort).SetLimit(length).SetSkip(skip)

	cursor, err := m.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var results []recommend.VideoMeta
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	for i := range results {
		results[i].Id = results[i].M_id.Hex()
	}

	return results, nil
}
