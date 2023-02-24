/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-03 17:20:02
 * @FilePath: /ChannelStudio/internal/ChannelStudio/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"encoding/json"
	"net/http"
	videoEnpoint "sideTube/ChannelStudio/internal/ChannelStudio/endpoint"
	"sideTube/ChannelStudio/internal/middleware"

	"context"
	"sideTube/ChannelStudio/internal/ChannelStudio/service"
	"sideTube/ChannelStudio/internal/common/simpleKit/endpoint"
	httptransport "sideTube/ChannelStudio/internal/common/simpleKit/httpTransport"

	"github.com/go-playground/validator"
)

func GetVideoListRegister(svc service.VideoStudioCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeGetVideoListEndPoint(svc)
	ep = middleware.ValidMiddleWare(v)(ep)
	ep = middleware.JwtMiddleWare()(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeGetVideoListRequest,
		encodeGetVideoListResponse,
		httptransport.NewServerBefore(middleware.JwtServerBerore()),
	)
}

func decodeGetVideoListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request videoEnpoint.GetVideoListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeGetVideoListResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
