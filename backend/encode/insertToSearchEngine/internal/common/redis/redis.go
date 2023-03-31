/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 20:49:23
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-16 12:23:21
 * @FilePath: /toSearchEngine/lockRepo/redis.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package redis

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func init() {
	var DB int
	dbString := os.Getenv("REDIS_DB")

	if dbString != "" {
		if dB, err := strconv.Atoi(dbString); err != nil {
			log.Fatalln("redis db strconv error : ", err)
		} else {
			DB = dB
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       DB,                          // use default DB
	})

	client = rdb
}

func GetRedisClient() *redis.Client {
	return client
}
