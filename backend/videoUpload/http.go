/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-04 17:36:26
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-09-02 11:37:26
 * @FilePath: /videoUpload/internal/transport/http/http.go
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
	"sideTube/videoUpload/internal/common/mongodb"
	"sideTube/videoUpload/internal/common/rabbitmq"
	"sideTube/videoUpload/internal/videoUpload/messageQueue"
	"sideTube/videoUpload/internal/videoUpload/metaRepository"
	"sideTube/videoUpload/internal/videoUpload/service"
	vts "sideTube/videoUpload/internal/videoUpload/transport/http"
	"sideTube/videoUpload/internal/videoUpload/videoRepository"
	"syscall"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func startHttpServer() {
	// Get aws session
	// AWSsession := awsS3.GetAWSSession()
	mongoDB := mongodb.GetMongoClient()
	rabbitmq := rabbitmq.GetRabbitClient()
	// invoke video-upload service
	svc := service.NewVideoCommend(
		metaRepository.NewMongoRepo(mongoDB),
		videoRepository.NewLoacl(os.Getenv("VIDEO_PATH")),
		messageQueue.NewMessageRepo(rabbitmq),
		// videoRepository.NewS3(AWSsession),
	)

	vaild := validator.New()

	// register apis
	r := mux.NewRouter()
	s := r.PathPrefix("/api/videoUpload").Subrouter()

	s.Handle("/start", vts.StartUploadRegister(svc, vaild)).Methods("POST")

	s.Handle("/updatePart", vts.UploadPartRegister(svc, vaild)).Methods("POST")

	s.Handle("/abort", vts.AbortUploadRegister(svc, vaild)).Methods("POST")

	s.Handle("/finish", vts.FinishUploadRegister(svc, vaild)).Methods("POST")

	// start Http server
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", os.Getenv("HTTP_ADDRESS"), os.Getenv("HTTP_PORT")),
		Handler: accessControl(r),
	}

	go func() {

		log.Println("HTTP server is staring")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	gracefullyShutDown(srv)
}

// gracefully shutdown
// Wait for interrupt signal to gracefully shutdown the server with
// a timeout of 5 seconds.
func gracefullyShutDown(srv *http.Server) {

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("HTTP server close")
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CROS_ALLOW_ORIGIN"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
