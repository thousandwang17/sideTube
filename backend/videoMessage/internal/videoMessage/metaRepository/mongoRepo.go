/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 16:07:51
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-24 14:17:55
 * @FilePath: /videoMessage/internal/videoMessage/metaRepository/mongoRepo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package metaRepository

import (
	"context"
	"errors"
	"log"
	"os"
	"sideTube/videoMessage/internal/videoMessage"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type mongolRepo struct {
	client     *mongo.Client
	videoCol   *mongo.Collection
	messageCol *mongo.Collection
	ReplyCol   *mongo.Collection
}

var (
	ErrInvalidID        = errors.New("invalid video id")
	ErrInvalidMessageID = errors.New("invalid message id")
	ErrGenerateIDFailed = errors.New(" messageID generate failed ")
)

func NewMongoRepo(db *mongo.Client) videoMessage.MetaRepository {
	wcMajority := writeconcern.New(writeconcern.WMajority(), writeconcern.WTimeout(1*time.Second))
	wcMajorityCollectionOpts := options.Collection().SetWriteConcern(wcMajority)

	return &mongolRepo{
		client:     db,
		videoCol:   db.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_VIDEO_COLLECTION"), wcMajorityCollectionOpts),
		messageCol: db.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_MESSAGE_COLLECTION"), wcMajorityCollectionOpts),
		ReplyCol:   db.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_Reply_COLLECTION"), wcMajorityCollectionOpts),
	}
}

type MessageID struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

func (m mongolRepo) MessageList(c context.Context, videoID string, skip, length int64) (data []videoMessage.VideoMessageMeta, err error) {
	sort := bson.D{{Key: "createTime", Value: -1}}
	filter := bson.D{{Key: "videoId", Value: videoID}}

	opts := options.Find().SetSort(sort).SetLimit(length).SetSkip(skip)

	cursor, err := m.messageCol.Find(c, filter, opts)
	if err != nil {
		return nil, err
	}
	var results []videoMessage.VideoMessageMeta
	if err = cursor.All(c, &results); err != nil {
		return nil, err
	}

	// if  results length  is 0  return empty slice replace nil
	if results == nil {
		return []videoMessage.VideoMessageMeta{}, nil
	}
	return results, nil
}

func (m mongolRepo) Message(c context.Context, userId, userName, videoId, message string) (messageId string, err error) {

	objectId, err := primitive.ObjectIDFromHex(videoId)
	if err != nil {
		log.Println("Invalid video id", videoId)
		return "", ErrInvalidID
	}

	// // Step 1: Define the callback that specifies the sequence of operations to perform inside the transaction.
	// callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
	// Important: You must pass sessCtx as the Context parameter to the operations for them to be executed in the
	// transaction.

	// message collection insert new message
	rep, err := m.messageCol.InsertOne(c,
		bson.D{
			{Key: "userId", Value: userId},
			{Key: "userName", Value: userName},
			{Key: "videoId", Value: videoId},
			{Key: "message", Value: message},
			{Key: "createTime", Value: primitive.Timestamp{T: uint32(time.Now().Unix())}},
		},
	)

	if err != nil {
		return "", err
	}

	// user collection inc Messages count
	res := m.videoCol.FindOneAndUpdate(
		c,
		bson.D{
			{Key: "_id", Value: objectId},
		},
		bson.D{
			{Key: "$inc",
				Value: bson.D{
					{Key: "Messages", Value: 1},
				},
			},
		})

	if res.Err() != nil {
		return "", res.Err()
	}

	if messageID, ok := rep.InsertedID.(primitive.ObjectID); ok {
		return messageID.Hex(), nil
	}

	return "", ErrGenerateIDFailed
	// }
	// Step 2: Start a session and run the callback using WithTransaction.
	// session, err := m.client.StartSession()
	// if err != nil {
	// 	return "", err
	// }

	// defer session.EndSession(c)
	// result, err := session.WithTransaction(c, callback)
	// log.Printf("result: %v\n", result)

	// if err != nil {
	// 	return "", err
	// }

	// id, ok := result.(string)
	// if !ok {
	// 	return "", err
	// }

	// return id, nil
}

func (m mongolRepo) EditMessage(c context.Context, userId, messageId, message string) (err error) {

	objectId, err := primitive.ObjectIDFromHex(messageId)
	if err != nil {
		log.Println("Invalid message id", messageId)
		return ErrInvalidMessageID
	}

	res := m.messageCol.FindOneAndUpdate(
		c,
		bson.D{
			{Key: "_id", Value: objectId},
		},
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "message", Value: message},
					{Key: "lastUpdatetime", Value: primitive.Timestamp{T: uint32(time.Now().Unix())}},
				},
			},
		})

	return res.Err()
}

func (m mongolRepo) DeleteMessage(c context.Context, userId, messageId string) (err error) {

	objectId, err := primitive.ObjectIDFromHex(messageId)
	if err != nil {
		log.Println("Invalid message id", messageId)
		return ErrInvalidMessageID
	}

	// Step 1: Define the callback that specifies the sequence of operations to perform inside the transaction.
	// callback := func(sessCtx mongo.SessionContext) (interface{}, error) {

	// Important: You must pass sessCtx as the Context parameter to the operations for them to be executed in the
	// transaction.

	// message collection delete the message
	res := m.messageCol.FindOneAndDelete(
		c,
		bson.D{
			{Key: "_id", Value: objectId},
			{Key: "userId", Value: userId},
		},
	)

	if res.Err() == mongo.ErrNoDocuments {
		return nil
	} else if res.Err() != nil {
		return res.Err()
	}

	// user collection dec Messages count
	resVideo := m.videoCol.FindOneAndUpdate(
		c,
		bson.D{
			{Key: "_id", Value: objectId},
		},
		bson.D{
			{Key: "$inc",
				Value: bson.D{
					{Key: "Messages", Value: -1},
				},
			},
		})

	if resVideo.Err() != nil {
		return res.Err()
	}

	return nil
	// }

	// Step 2: Start a session and run the callback using WithTransaction.
	// session, err := m.client.StartSession()
	// if err != nil {
	// 	return err
	// }

	// defer session.EndSession(c)
	// _, err = session.WithTransaction(c, callback)

	// if err != nil {
	// 	return err
	// }
	// return nil
}

// reply part

func (m mongolRepo) ReplyList(c context.Context, MessageID string, skip, length int64) (data []videoMessage.VideoMessageReplyMeta, err error) {
	sort := bson.D{{Key: "createTime", Value: -1}}
	filter := bson.D{{Key: "messageId", Value: MessageID}}

	opts := options.Find().SetSort(sort).SetLimit(length).SetSkip(skip)

	cursor, err := m.ReplyCol.Find(c, filter, opts)
	if err != nil {
		return nil, err
	}
	var results []videoMessage.VideoMessageReplyMeta
	if err = cursor.All(c, &results); err != nil {
		return nil, err
	}

	// if  results length  is 0  return empty slice replace nil
	if results == nil {
		return []videoMessage.VideoMessageReplyMeta{}, nil
	}

	return results, nil
}

func (m mongolRepo) Reply(c context.Context, userId, userName, messageID, message string) (messageId string, err error) {

	objectId, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		log.Println("Invalid message id", messageID, err.Error())
		return "", ErrInvalidMessageID
	}

	// ** NOTE : Transaction only can work on replica mongos cluster , if for local development should mark "Transaction codes"

	// Step 1: Define the callback that specifies the sequence of operations to perform inside the transaction.
	// callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
	// Important: You must pass sessCtx as the Context parameter to the operations for them to be executed in the
	// transaction.

	// Reply collection insert new message
	rep, err := m.ReplyCol.InsertOne(c,
		bson.D{
			{Key: "userId", Value: userId},
			{Key: "userName", Value: userName},
			{Key: "messageId", Value: messageID},
			{Key: "message", Value: message},
			{Key: "createTime", Value: primitive.Timestamp{T: uint32(time.Now().Unix())}},
		},
	)

	if err != nil {
		return "", err
	}

	// user collection inc Messages count
	res := m.messageCol.FindOneAndUpdate(
		c,
		bson.D{
			{Key: "_id", Value: objectId},
		},
		bson.D{
			{Key: "$inc",
				Value: bson.D{
					{Key: "Replies", Value: 1},
				},
			},
		})

	if res.Err() != nil {
		return "", res.Err()
	}

	if messageID, ok := rep.InsertedID.(primitive.ObjectID); ok {
		return messageID.Hex(), nil
	}

	return "", ErrGenerateIDFailed
	// }
	// Step 2: Start a session and run the callback using WithTransaction.
	// session, err := m.client.StartSession()
	// if err != nil {
	// 	return "", err
	// }

	// defer session.EndSession(c)
	// result, err := session.WithTransaction(c, callback)
	// log.Printf("result: %v\n", result)

	// if err != nil {
	// 	return "", err
	// }

	// id, ok := result.(string)
	// if !ok {
	// 	return "", err
	// }

	// return id, nil
}

func (m mongolRepo) EditReply(c context.Context, userId, messageID, message string) (err error) {

	objectId, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		log.Println("Invalid message id", messageID)
		return ErrInvalidMessageID
	}

	res := m.ReplyCol.FindOneAndUpdate(
		c,
		bson.D{
			{Key: "_id", Value: objectId},
		},
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "message", Value: message},
					{Key: "lastUpdatetime", Value: primitive.Timestamp{T: uint32(time.Now().Unix())}},
				},
			},
		})

	return res.Err()
}

func (m mongolRepo) DeleteReply(c context.Context, userId, messageID string) (err error) {

	objectId, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		log.Println("Invalid message id", messageID)
		return ErrInvalidMessageID
	}

	// Step 1: Define the callback that specifies the sequence of operations to perform inside the transaction.
	// callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
	// Important: You must pass sessCtx as the Context parameter to the operations for them to be executed in the
	// transaction.

	// message collection delete the message
	res := m.ReplyCol.FindOneAndDelete(
		c,
		bson.D{
			{Key: "_id", Value: objectId},
			{Key: "userId", Value: userId},
		},
	)

	if res.Err() == mongo.ErrNoDocuments {
		return nil
	} else if res.Err() != nil {
		return res.Err()
	}

	// user collection dec Reply count
	resMessage := m.messageCol.FindOneAndUpdate(
		c,
		bson.D{
			{Key: "_id", Value: objectId},
		},
		bson.D{
			{Key: "$inc",
				Value: bson.D{
					{Key: "Replies", Value: -1},
				},
			},
		})

	if resMessage.Err() != nil {
		return res.Err()
	}

	return nil
	// }

	// // Step 2: Start a session and run the callback using WithTransaction.
	// session, err := m.client.StartSession()
	// if err != nil {
	// 	return err
	// }

	// defer session.EndSession(c)
	// _, err = session.WithTransaction(c, callback)

	// if err != nil {
	// 	return err
	// }
	// return nil

}
