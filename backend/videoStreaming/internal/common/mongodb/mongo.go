/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-04 19:59:38
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-08-20 17:08:13
 * @FilePath: /VideoStreaming/internal/common/mongodb/mongo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	credential := options.Credential{
		Username: os.Getenv("MONGODB_USERNAME"),
		Password: os.Getenv("MONGODB_PASSWORD"),
	}

	uri := fmt.Sprintf("mongodb://%s", os.Getenv("MONGODB_HOST"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(credential))

	if err != nil {
		log.Fatal("mongodb connect failed : ", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatal("mongodb Ping failed : ", err)
	}

	log.Print("mongo server is connect")

	mongoClient = client
}

func GetMongoClient() *mongo.Client {
	return mongoClient
}
