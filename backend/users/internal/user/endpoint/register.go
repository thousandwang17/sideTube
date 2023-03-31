/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-28 15:38:07
 * @FilePath: /user/internal/user/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"sideTube/user/internal/common/simpleKit/endpoint"
	"sideTube/user/internal/user/service"
)

func MakeRegisterEndPoint(v service.UserCommend) endpoint.EndPoint {
	return func(c context.Context, req interface{}) (interface{}, error) {
		data := req.(RegisterRequest)
		err := v.Register(c, data.Email, data.PassWord, data.Name)
		return RegisterRespond{err}, err
	}
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"email,lte=50"`
	PassWord string `json:"passWord" validate:"required,lte=30"`
	Name     string `json:"name" validate:"required,lte=50"`
}

type RegisterRespond struct {
	Err error `json:"error,omitempty"`
}
