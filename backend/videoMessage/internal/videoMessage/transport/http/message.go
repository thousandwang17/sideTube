/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-11 11:45:15
 * @FilePath: /videoMessage/internal/videoMessage/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"encoding/json"
	"net/http"
	"sideTube/videoMessage/internal/middleware"
	videoEnpoint "sideTube/videoMessage/internal/videoMessage/endpoint"
	"sideTube/videoMessage/internal/videoMessage/service"

	"context"
	"sideTube/videoMessage/internal/common/simpleKit/endpoint"
	httptransport "sideTube/videoMessage/internal/common/simpleKit/httpTransport"

	"github.com/go-playground/validator"
)

func MessageRegister(svc service.VideoMessageCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeMessageEndPoint(svc)
	ep = middleware.JwtMiddleWare()(ep)
	ep = middleware.ValidMiddleWare(v)(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeMessageRequest,
		encodeMessageResponse,
		httptransport.NewServerBefore(middleware.JwtServerBerore()),
	)
}

func decodeMessageRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request videoEnpoint.MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeMessageResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
