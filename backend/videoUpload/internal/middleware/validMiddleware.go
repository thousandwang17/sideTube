/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-18 17:47:14
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-19 19:39:09
 * @FilePath: /videoUpload/internal/middleware/validMiddleware.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package middleware

import (
	"context"
	"net/http"
	"sideTube/videoUpload/internal/common/simpleKit/endpoint"

	"github.com/go-playground/validator"
)

type validerror struct {
	err error
}

func (_ validerror) StatusCode() int {
	return http.StatusBadRequest
}

func (v validerror) Error() string {
	return v.err.Error()
}

func validMiddleWare() endpoint.MiddleWare {
	return func(next endpoint.EndPoint) endpoint.EndPoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			validate := validator.New()
			err := validate.Struct(request)
			if err != nil {
				return nil, validerror{err}
			}
			return next(ctx, request)
		}
	}

}
