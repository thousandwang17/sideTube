/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-30 13:22:07
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-28 20:21:45
 * @FilePath: /user/cmd/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sideTube/user/internal/common/mongodb"
	_ "sideTube/user/internal/common/mongodb"
	"sideTube/user/internal/user/metaRepository"
	"sideTube/user/internal/user/service"
	"sideTube/user/thrid-part/JwtGererater"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

func main() {

	mongoDB := mongodb.GetMongoClient()

	// JwtGererater grcp - client
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", os.Getenv("GRPC_ADDRESS"), os.Getenv("GRPC_PORT")), grpc.WithInsecure())
	if err != nil {
		log.Panic(err)
	}

	// srv
	svc := service.NewUserCommend(
		metaRepository.NewMongoRepo(mongoDB),
		JwtGererater.NewJwtTokenClient(conn),
	)

	httpServer := startHttpServer(svc)
	gracefullyShutDown(httpServer)

}

// gracefully shutdown
// Wait for interrupt signal to gracefully shutdown the server with
// a timeout of 5 seconds.
func gracefullyShutDown(hsrv *http.Server) {

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := hsrv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown: ", err)
	}

	log.Println("server close")
}
