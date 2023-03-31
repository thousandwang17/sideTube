/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-28 14:32:21
 * @FilePath: /user/internal/user/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"encoding/json"
	"net/http"
	"sideTube/user/internal/middleware"
	userEnpoint "sideTube/user/internal/user/endpoint"
	"sideTube/user/internal/user/service"

	"context"
	"sideTube/user/internal/common/simpleKit/endpoint"
	httptransport "sideTube/user/internal/common/simpleKit/httpTransport"

	"github.com/go-playground/validator"
)

func RegisterRegister(svc service.UserCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = userEnpoint.MakeRegisterEndPoint(svc)
	ep = middleware.ValidMiddleWare(v)(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeRegisterRequest,
		encodeRegisterResponse,
	)
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request userEnpoint.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeRegisterResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
