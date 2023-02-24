/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-01 13:53:53
 * @FilePath: /ChannelStudio/internal/ChannelStudio/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"context"
	"encoding/json"
	"net/http"
	videoEnpoint "sideTube/ChannelStudio/internal/ChannelStudio/endpoint"
	"sideTube/ChannelStudio/internal/ChannelStudio/service"
	"sideTube/ChannelStudio/internal/common/simpleKit/endpoint"
	httptransport "sideTube/ChannelStudio/internal/common/simpleKit/httpTransport"
	"sideTube/ChannelStudio/internal/middleware"

	"github.com/go-playground/validator"
)

func EditVideoPublicStateRegister(svc service.VideoStudioCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeEditVideoPublicStateEndPoint(svc)
	ep = middleware.ValidMiddleWare(v)(ep)
	ep = middleware.JwtMiddleWare()(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeEditVideoPublicStateRequest,
		encodeEditVideoPublicStateResponse,
		httptransport.NewServerBefore(middleware.JwtServerBerore()),
	)
}

func decodeEditVideoPublicStateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request videoEnpoint.EditVideoPublicStateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeEditVideoPublicStateResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
