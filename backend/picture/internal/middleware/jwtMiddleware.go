/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-04 16:46:31
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-01-22 17:30:32
 * @FilePath: /picture/internal/middleware/jwtMiddleware.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sideTube/picture/internal/common/simpleKit/endpoint"
	httptransport "sideTube/picture/internal/common/simpleKit/httpTransport"

	"github.com/golang-jwt/jwt"
)

type jWTerror struct{}

func (_ jWTerror) StatusCode() int {
	return http.StatusForbidden
}

func (_ jWTerror) Error() string {
	return "auth failed"
}

func JwtServerBerore() httptransport.ServerBefore {
	return func(ctx context.Context, r *http.Request) context.Context {
		ctx = context.WithValue(ctx, "jwtToken", r.Header.Get("Authorization"))
		return ctx
	}
}

func JwtMiddleWare() endpoint.MiddleWare {
	return func(next endpoint.EndPoint) endpoint.EndPoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {

			alg := os.Getenv("JWT_HEADER_ALG")
			jwt_secret_key := os.Getenv("JWT_SECRET_KEY")

			if alg == "" || jwt_secret_key == "" {
				fmt.Println("losing environment variables")
				return nil, errors.New("losing environment variables")
			}

			// Parse takes the token string and a function for looking up the key. The latter is especially
			// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
			// head of the token to identify which key to use, but the parsed token (head and claims) is provided
			// to the callback, providing flexibility.
			token, err := jwt.Parse(ctx.Value("jwtToken").(string), func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				if token.Header["alg"] != os.Getenv("JWT_HEADER_ALG") {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})

			if err != nil {
				return nil, jWTerror{}
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if claims["uid"] != nil {
					if uid, ok := claims["uid"].(string); ok {
						ctx = context.WithValue(ctx, "uid", uid)
					} else {
						return nil, jWTerror{}
					}
				} else {
					log.Println("losing require alg: uid ")
					return nil, jWTerror{}
				}
			} else {
				return nil, jWTerror{}
			}
			return next(ctx, request)
		}
	}

}
