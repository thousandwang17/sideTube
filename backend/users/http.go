/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-04 17:36:26
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-29 17:59:38
 * @FilePath: /user/internal/transport/http/http.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sideTube/user/internal/user/service"
	vts "sideTube/user/internal/user/transport/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func startHttpServer(svc service.UserCommend) *http.Server {

	vaild := validator.New()

	// register apis
	r := mux.NewRouter()
	u := r.PathPrefix("/user").Subrouter()

	u.Handle("/login", vts.LoginRegister(svc, vaild)).Methods("POST")
	u.Handle("/logout", vts.LogoutRegister(svc, vaild)).Methods("POST")
	u.Handle("/register", vts.RegisterRegister(svc, vaild)).Methods("POST")
	u.Handle("/history", vts.GetHistoryListRegister(svc, vaild)).Methods("POST")
	u.Handle("/info", vts.PublicUserInfoRegister(svc, vaild)).Methods("POST")

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

	return srv

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CROS_ALLOW_ORIGIN"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
