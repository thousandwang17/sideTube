/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 16:07:51
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-24 17:18:50
 * @FilePath: /jwtGenerate/internal/jwtGenerate/metaRepository/mongoRepo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package metaRepository

import (
	"context"
	"errors"
	"os"
	"sideTube/jwtGenerate/internal/jwtGenerate"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type mongolRepo struct {
	client  *mongo.Client
	userCol *mongo.Collection
}

var (
	ErrInvalidID        = errors.New("invalid video id")
	ErrInvalidMessageID = errors.New("invalid message id")
	ErrGenerateIDFailed = errors.New(" messageID generate failed ")
)

func NewMongoRepo(db *mongo.Client) jwtGenerate.MetaRepository {
	wcMajority := writeconcern.New(writeconcern.WMajority(), writeconcern.WTimeout(1*time.Second))
	wcMajorityCollectionOpts := options.Collection().SetWriteConcern(wcMajority)

	return &mongolRepo{
		client:  db,
		userCol: db.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_USER_COLLECTION"), wcMajorityCollectionOpts),
	}
}

type MessageID struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

func (m mongolRepo) LogInCheck(c context.Context, account, passWordSH256 string) (info jwtGenerate.UserInfo, err error) {

	// user collection inc Messages count
	res := m.userCol.FindOneAndUpdate(
		c,
		bson.D{
			{Key: "account", Value: account},
			{Key: "passWord", Value: passWordSH256},
		},
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "lastLogInTime", Value: primitive.Timestamp{T: uint32(time.Now().Unix())}},
				},
			},
		})

	if res.Err() != nil {
		return info, res.Err()
	}

	if err = res.Decode(info); err != nil {
		return info, err
	}

	return info, nil
}
