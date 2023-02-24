/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-24 17:44:32
 * @FilePath: /jwtGenerate/internal/jwtGenerate/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"encoding/json"
	"net/http"
	videoEnpoint "sideTube/jwtGenerate/internal/jwtGenerate/endpoint"
	"sideTube/jwtGenerate/internal/jwtGenerate/service"
	"sideTube/jwtGenerate/internal/middleware"
	"time"

	"context"
	"sideTube/jwtGenerate/internal/common/simpleKit/endpoint"
	httptransport "sideTube/jwtGenerate/internal/common/simpleKit/httpTransport"
)

func RefreshTokenRegister(svc service.JwtGenerateCommend) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeRefreshTokenEndPoint(svc)
	ep = middleware.JwtMiddleWare()(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeRefreshTokenRequest,
		encodeRefreshTokenResponse,
		httptransport.NewServerBefore(middleware.JwtServerBerore()),
	)
}

func decodeRefreshTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeRefreshTokenResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(videoEnpoint.RefreshTokenRespond)
	w.Header().Add("Content-Type", "application/json")
	cookie := http.Cookie{
		Name:     "RefreshToken",
		Value:    res.Token,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return json.NewEncoder(w).Encode(response)
}
