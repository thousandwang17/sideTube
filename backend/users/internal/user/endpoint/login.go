/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-28 15:37:30
 * @FilePath: /user/internal/user/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"sideTube/user/internal/common/simpleKit/endpoint"
	"sideTube/user/internal/user/service"
)

func MakeLoginEndPoint(v service.UserCommend) endpoint.EndPoint {
	return func(c context.Context, req interface{}) (interface{}, error) {
		data := req.(LoginRequest)
		rep, err := v.Login(c, data.Email, data.PassWord)
		return LoginRespond{string(rep.AT), string(rep.RT), err}, err
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"email,lte=50"`
	PassWord string `json:"passWord" validate:"required,lte=30"`
}

type LoginRespond struct {
	AccessToken  string `json:"-"`
	RefreshToken string `json:"-"`
	Err          error  `json:"error,omitempty"`
}
