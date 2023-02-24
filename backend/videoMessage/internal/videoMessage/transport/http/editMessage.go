/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-11 11:47:07
 * @FilePath: /videoEditMessage/internal/videoEditMessage/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"context"
	"encoding/json"
	"net/http"
	"sideTube/videoMessage/internal/common/simpleKit/endpoint"
	httptransport "sideTube/videoMessage/internal/common/simpleKit/httpTransport"
	"sideTube/videoMessage/internal/middleware"
	"sideTube/videoMessage/internal/videoMessage/service"

	videoEnpoint "sideTube/videoMessage/internal/videoMessage/endpoint"

	"github.com/go-playground/validator"
)

func EditMessageRegister(svc service.VideoMessageCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeEditMessageEndPoint(svc)
	ep = middleware.JwtMiddleWare()(ep)
	ep = middleware.ValidMiddleWare(v)(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeEditMessageRequest,
		encodeEditMessageResponse,
		httptransport.NewServerBefore(middleware.JwtServerBerore()),
	)
}

func decodeEditMessageRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request videoEnpoint.EditMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeEditMessageResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
