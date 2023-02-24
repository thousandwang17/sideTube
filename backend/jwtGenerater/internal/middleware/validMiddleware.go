/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-18 17:47:14
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-07 13:57:36
 * @FilePath: /jwtGenerate/internal/middleware/validMiddleware.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package middleware

import (
	"context"
	"net/http"
	"sideTube/jwtGenerate/internal/common/simpleKit/endpoint"

	"github.com/go-playground/validator"
)

const httpErrCode = http.StatusUnprocessableEntity

type VaildError struct {
	E error
}

func (v VaildError) Error() string {
	return v.E.Error()
}

func (v VaildError) StatusCode() int {
	return httpErrCode
}

func ValidMiddleWare(valid *validator.Validate) endpoint.MiddleWare {
	return func(next endpoint.EndPoint) endpoint.EndPoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			err := valid.Struct(request)
			if err != nil {
				e := VaildError{err}
				return nil, e
			}
			return next(ctx, request)
		}
	}
}
