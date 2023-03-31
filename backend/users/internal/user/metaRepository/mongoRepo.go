/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 16:07:51
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-30 18:32:36
 * @FilePath: /user/internal/user/metaRepository/mongoRepo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package metaRepository

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"sideTube/user/internal/user"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrDatabase = errors.New("dataBase error!!")
	InvalidID   = errors.New("invalid video id")
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

type mongolRepo struct {
	client     *mongo.Client
	userCol    *mongo.Collection
	videoCol   *mongo.Collection
	viewLogCol *mongo.Collection
}

func NewMongoRepo(db *mongo.Client) user.MetaRepository {
	// wcMajority := writeconcern.New(writeconcern.WMajority(), writeconcern.WTimeout(1*time.Second))
	// wcMajorityCollectionOpts := options.Collection().SetWriteConcern(wcMajority)

	return &mongolRepo{
		client:     db,
		userCol:    db.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_USER_COLLECTION")),
		videoCol:   db.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_VIDEO_COLLECTION")),
		viewLogCol: db.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_HISTORY_COLLECTION")),
	}
}

type MessageID struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

func (m mongolRepo) LogInCheck(c context.Context, email, password string) (user.UserInfo, error) {

	// user collection inc Messages count
	res := m.userCol.FindOneAndUpdate(
		c,
		bson.D{
			{Key: "email", Value: email},
			{Key: "passWord", Value: password},
		},
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "lastLogInTime", Value: primitive.Timestamp{T: uint32(time.Now().Unix())}},
				},
			},
		})

	info := user.UserInfo{}
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return info, mongoErr{errors.New("Account is not exist or typing a wrong password")}
		}

		log.Println("LogInCheck err: ", res.Err())
		return info, ErrDatabase
	}

	if err := res.Decode(&info); err != nil {
		log.Println("LogInCheck Decode err: ", err)
		return info, ErrDatabase
	}

	info.UserId = info.ID.Hex()
	return info, nil
}

func (m mongolRepo) Register(c context.Context, email, password, name string) error {

	// user collection inc Messages count
	_, err := m.userCol.InsertOne(
		c,
		bson.D{
			{Key: "email", Value: email},
			{Key: "name", Value: name},
			{Key: "passWord", Value: password},
			{Key: "createTime", Value: primitive.Timestamp{T: uint32(time.Now().Unix())}}})

	if err != nil {
		// Check if the error is a duplicate key violation
		if mongo.IsDuplicateKeyError(err) {
			return mongoErr{errors.New("The Email alerdy been used")}
		}
		log.Println("Register err: ", err)
		return ErrDatabase
	}

	return nil
}

func (m mongolRepo) GetPublicUserInfo(c context.Context, user_id string) (user.PublicUserInfo, error) {
	objectId, err := primitive.ObjectIDFromHex(user_id)

	info := user.PublicUserInfo{}

	if err != nil {
		log.Println("Invalid user id", user_id)
		return info, InvalidID
	}

	// user collection inc Messages count
	res := m.userCol.FindOne(
		c,
		bson.D{
			{Key: "_id", Value: objectId},
		})

	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return info, mongoErr{errors.New("Account is not exist or typing a wrong password")}
		}

		log.Println("LogInCheck err: ", res.Err())
		return info, ErrDatabase
	}

	if err := res.Decode(&info); err != nil {
		log.Println("LogInCheck Decode err: ", err)
		return info, ErrDatabase
	}

	info.UserId = info.ID.Hex()
	return info, nil
}

func (m mongolRepo) GetHistoryList(ctx context.Context, userId string, skip int64, length int64) ([]user.VideoMeta, error) {

	filter := bson.D{{Key: "user_id", Value: userId}}
	sort := bson.D{{Key: "date", Value: -1}}
	opts := options.Find().SetSort(sort).SetLimit(length).SetSkip(skip)

	cursor, err := m.viewLogCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var historyResults []user.HistoryMeta
	if err = cursor.All(ctx, &historyResults); err != nil {
		return nil, err
	}

	idMap := map[primitive.ObjectID]int{}
	ids := make([]primitive.ObjectID, 0, len(historyResults))
	for i := range historyResults {
		objectID, err := primitive.ObjectIDFromHex(historyResults[i].Video_id)
		if err != nil {
			continue
		}
		ids = append(ids, objectID)
		idMap[objectID] = i
	}

	/// video meta ///
	videofilter := bson.M{"_id": bson.M{"$in": ids}}

	videoCursor, err := m.videoCol.Find(ctx, videofilter)
	if err != nil {
		return nil, err
	}

	var results []user.VideoMeta
	if err = videoCursor.All(ctx, &results); err != nil {
		return nil, err
	}

	for i := range results {
		results[i].Id = results[i].M_id.Hex()
		results[i].HistoryTime = historyResults[idMap[results[i].M_id]].HistoryTime
	}

	return results, nil
}
