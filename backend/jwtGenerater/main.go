/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-30 13:22:07
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-28 13:48:41
 * @FilePath: /jwtGenerate/cmd/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sideTube/jwtGenerate/internal/jwtGenerate/service"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

var (
	grpcShutDown chan struct{}
	httpShutDown chan struct{}
)

func main() {

	svc := service.NewjwtGenerateCommend()

	grpcServer := startGrpcServer(svc)
	httpServer := startHttpServer(svc)
	gracefullyShutDown(httpServer, grpcServer)

}

// gracefully shutdown
// Wait for interrupt signal to gracefully shutdown the server with
// a timeout of 5 seconds.
func gracefullyShutDown(hsrv *http.Server, gsrv *grpc.Server) {

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	go func() {
		gsrv.GracefulStop()
		close(grpcShutDown)
	}()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := hsrv.Shutdown(ctx); err != nil {
			log.Println("Server Shutdown: ", err)
		}
		close(httpShutDown)
	}()

	select {
	case <-grpcShutDown:
		log.Println("gRPC server closed")
	case <-time.After(5 * time.Second):
		log.Println("gRPC server shutdown timed out")
	}

	select {
	case <-httpShutDown:
		log.Println("HTTP server closed")
	case <-time.After(5 * time.Second):
		log.Println("HTTP server shutdown timed out")
	}

	log.Println("server close")
}
