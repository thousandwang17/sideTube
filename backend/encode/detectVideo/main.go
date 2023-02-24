/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 16:24:43
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-19 19:48:53
 * @FilePath: /encode/encodeVideo/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"detectVideo/internal/common/rabbitmq"
	locker "detectVideo/internal/common/redis"
	"detectVideo/internal/worker"
	"detectVideo/internal/worker/handler"
	"detectVideo/internal/worker/lockRepo"
	"detectVideo/internal/worker/queueRepo"
	"detectVideo/internal/worker/videoRepo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// when worker upexpectly shutdown or connect lost, if we want to avoid two worker process same message ,
// there use distrube lock system and re-queue message with delay strategy, to ensure every message only handle by a worker
// and implementing a retry strategy to maximize the number of messages that can be successfully handled.
// if message reach max retries, we send to dead-letter queue for urther processing.
func main() {
	rabbit := rabbitmq.GetRabbitClient()
	worker := handler.New(
		videoRepo.NewLoacl(os.Getenv("VIDEO_PATH")),
		lockRepo.NewRedis(locker.GetRedisClient()),
		queueRepo.NewRabbitmq(rabbit),
	)

	go worker.Handle()

	gracefullyShutDown(worker)
	rabbit.Close()
}

// gracefully shutdown
// Wait for interrupt signal to gracefully shutdown the worker
func gracefullyShutDown(srv worker.Handler) {

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	srv.Shutdown()
	log.Println(" worker closed")
}
