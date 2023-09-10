/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-06-27 08:57:12
 * @FilePath: /search/internal/search/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"encoding/json"
	"net/http"
	"sideTube/search/internal/middleware"
	videoEnpoint "sideTube/search/internal/search/endpoint"

	"context"
	"sideTube/search/internal/common/simpleKit/endpoint"
	httptransport "sideTube/search/internal/common/simpleKit/httpTransport"
	"sideTube/search/internal/search/service"

	"github.com/go-playground/validator"
)

func SearchVideoRegister(svc service.SearchCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeSearchVideoEndPoint(svc)
	ep = middleware.ValidMiddleWare(v)(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeSearchVideoRequest,
		encodeSearchVideoResponse,
		httptransport.NewServerBefore(middleware.JwtServerBerore()),
	)
}

func decodeSearchVideoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request videoEnpoint.SearchVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeSearchVideoResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
