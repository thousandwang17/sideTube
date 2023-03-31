/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-28 21:07:36
 * @FilePath: /jwtGenerate/internal/jwtGenerate/service/service.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"context"
	"errors"
	"log"
	"os"
	"sideTube/jwtGenerate/internal/jwtGenerate"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrVaild = errors.New("args vaild failed")
)

const (
	private = iota
	publish
)

type service struct {
}

type JwtGenerateCommend interface {
	RefreshToken(c context.Context) (token jwtGenerate.RefreshToken, err error)
}

func NewjwtGenerateCommend() JwtGenerateCommend {
	return service{}
}

func (v service) RefreshToken(c context.Context) (jwtGenerate.RefreshToken, error) {
	userId := c.Value("uid").(string)
	name := c.Value("userName").(string)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	ApiToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":  userId,
		"name": name,
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	toString, err := ApiToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	if err != nil {
		log.Println("SignedString err: ", err)
		return "", nil
	}

	return jwtGenerate.RefreshToken(toString), nil
}
