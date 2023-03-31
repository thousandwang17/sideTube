/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-03 15:46:26
 * @FilePath: /recommend/internal/recommend/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"encoding/json"
	"net/http"
	"sideTube/recommend/internal/middleware"
	videoEnpoint "sideTube/recommend/internal/recommend/endpoint"

	"context"
	"sideTube/recommend/internal/common/simpleKit/endpoint"
	httptransport "sideTube/recommend/internal/common/simpleKit/httpTransport"
	"sideTube/recommend/internal/recommend/service"

	"github.com/go-playground/validator"
)

func GetRelationVideoListRegister(svc service.RecommendCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeGetRelationVideoListEndPoint(svc)
	ep = middleware.ValidMiddleWare(v)(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeGetRelationVideoListRequest,
		encodeGetRelationVideoListResponse,
		httptransport.NewServerBefore(middleware.JwtServerBerore()),
	)
}

func decodeGetRelationVideoListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request videoEnpoint.GetRelationVideoListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeGetRelationVideoListResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
