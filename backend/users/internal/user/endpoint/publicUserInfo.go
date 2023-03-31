/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-06 19:23:42
 * @FilePath: /user/internal/user/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"sideTube/user/internal/common/simpleKit/endpoint"
	"sideTube/user/internal/user"
	"sideTube/user/internal/user/service"
)

func MakePublicUserInfoEndPoint(v service.UserCommend) endpoint.EndPoint {
	return func(c context.Context, req interface{}) (interface{}, error) {
		data := req.(PublicUserInfoRequest)
		rep, err := v.GetPublicUserInfo(c, data.UserId)
		return PublicUserInfoRespond{rep, err}, err
	}
}

type PublicUserInfoRequest struct {
	UserId string `json:"user_id" validate:"required,alphanum"`
}

type PublicUserInfoRespond struct {
	user.PublicUserInfo
	Err error `json:"error,omitempty"`
}
