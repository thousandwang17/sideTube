/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-08 16:53:41
 * @FilePath: /VideoStreaming/internal/VideoStreaming/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"encoding/json"
	"net/http"
	videoEnpoint "sideTube/VideoStreaming/internal/VideoStreaming/endpoint"
	"sideTube/VideoStreaming/internal/middleware"

	"context"
	"sideTube/VideoStreaming/internal/VideoStreaming/service"
	"sideTube/VideoStreaming/internal/common/simpleKit/endpoint"
	httptransport "sideTube/VideoStreaming/internal/common/simpleKit/httpTransport"

	"github.com/go-playground/validator"
)

func IncVideoViewsRegister(svc service.VideoStreamingCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeIncVideoViewsEndPoint(svc)
	ep = middleware.ValidMiddleWare(v)(ep)
	ep = middleware.JwtMiddleWare()(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeIncVideoViewsRequest,
		encodeIncVideoViewsResponse,
		httptransport.NewServerBefore(middleware.JwtServerBerore()),
	)
}

func decodeIncVideoViewsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request videoEnpoint.IncVideoViewsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeIncVideoViewsResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
