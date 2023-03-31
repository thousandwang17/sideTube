/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-06 21:50:03
 * @FilePath: /user/internal/user/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"encoding/json"
	"net/http"
	userEnpoint "sideTube/user/internal/user/endpoint"
	"sideTube/user/internal/user/service"
	"time"

	"context"
	"sideTube/user/internal/common/simpleKit/endpoint"
	httptransport "sideTube/user/internal/common/simpleKit/httpTransport"

	"github.com/go-playground/validator"
)

func LogoutRegister(svc service.UserCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = userEnpoint.MakeLogoutEndPoint(svc)

	return httptransport.NewHttpTransport(
		ep,
		decodeLogoutRequest,
		encodeLogoutResponse,
	)
}

func decodeLogoutRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request userEnpoint.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeLogoutResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	// res := response.(userEnpoint.LogoutRespond)
	w.Header().Add("Content-Type", "application/json")

	AccessCookie := http.Cookie{
		Name:     "AccessToken",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, &AccessCookie)
	RefreshCookie := http.Cookie{
		Name:     "RefreshToken",
		Value:    "res.RefreshToken",
		Expires:  time.Now().Add(-1 * time.Hour),
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, &RefreshCookie)

	return json.NewEncoder(w).Encode(response)
}
