/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-06 15:40:58
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

func GetPublicVideoListRegister(svc service.VideoStudioCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeGetPublicVideoListEndPoint(svc)
	ep = middleware.ValidMiddleWare(v)(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeGetPublicVideoListRequest,
		encodeGetPublicVideoListResponse,
		httptransport.NewServerBefore(middleware.JwtServerBerore()),
	)
}

func decodeGetPublicVideoListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request videoEnpoint.GetPublicVideoListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeGetPublicVideoListResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
