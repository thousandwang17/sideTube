/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-26 18:24:51
 * @FilePath: /videoUpload/internal/videoUpload/endpoint/abortUpload.go
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

	"github.com/go-playground/validator"
)

func AbortUploadRegister(svc service.VideoCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeAbortUploadEndPoint(svc)
	ep = middleware.ValidMiddleWare(v)(ep)
	ep = middleware.JwtMiddleWare()(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeAbortUploadRequest,
		encodeAbortUploadResponse,
		httptransport.NewServerBefore(middleware.JwtServerBerore()),
	)
}

func decodeAbortUploadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request videoEnpoint.AbortUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeAbortUploadResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
