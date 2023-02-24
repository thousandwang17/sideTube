/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-05 16:19:49
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-01-22 18:05:35
 * @FilePath: /videoUpload/internal/transport/http/startUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package http

import (
	"context"
	"encoding/json"
	"net/http"
	"sideTube/videoUpload/internal/common/simpleKit/endpoint"
	httptransport "sideTube/videoUpload/internal/common/simpleKit/httpTransport"
	"sideTube/videoUpload/internal/middleware"
	videoEnpoint "sideTube/videoUpload/internal/videoUpload/endpoint"
	"sideTube/videoUpload/internal/videoUpload/service"
)

func StartUploadRegister(svc service.VideoCommend) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeStartUploadEndPoint(svc)
	ep = middleware.JwtMiddleWare()(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeStartUploadRequest,
		encodeResponse,
		httptransport.NewServerBefore(middleware.JwtServerBerore()),
	)
}

func decodeStartUploadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request videoEnpoint.StartUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
