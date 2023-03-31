/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-28 20:58:49
 * @FilePath: /jwtGenerate/internal/jwtGenerate/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"sideTube/jwtGenerate/internal/common/simpleKit/endpoint"
	"sideTube/jwtGenerate/internal/jwtGenerate/service"
)

func MakeRefreshTokenEndPoint(v service.JwtGenerateCommend) endpoint.EndPoint {
	return func(c context.Context, _ interface{}) (interface{}, error) {
		token, err := v.RefreshToken(c)
		return RefreshTokenRespond{string(token), err}, err
	}
}

type RefreshTokenRequest struct {
}

type RefreshTokenRespond struct {
	Token string `json:"-"`
	Err   error  `json:"error,omitempty"`
}
