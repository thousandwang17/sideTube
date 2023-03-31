/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-30 16:21:15
 * @FilePath: /user/internal/user/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"encoding/json"
	"net/http"
	"sideTube/user/internal/middleware"
	videoEnpoint "sideTube/user/internal/user/endpoint"

	"context"
	"sideTube/user/internal/common/simpleKit/endpoint"
	httptransport "sideTube/user/internal/common/simpleKit/httpTransport"
	"sideTube/user/internal/user/service"

	"github.com/go-playground/validator"
)

func GetHistoryListRegister(svc service.UserCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeGetHistoryListEndPoint(svc)
	ep = middleware.ValidMiddleWare(v)(ep)
	ep = middleware.JwtMiddleWare()(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeGetHistoryListRequest,
		encodeGetHistoryListResponse,
		httptransport.NewServerBefore(middleware.JwtServerBerore()),
	)
}

func decodeGetHistoryListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request videoEnpoint.GetHistoryListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeGetHistoryListResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
