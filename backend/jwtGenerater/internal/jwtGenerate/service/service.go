/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-24 17:38:59
 * @FilePath: /jwtGenerate/internal/jwtGenerate/service/service.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"context"
	"errors"
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
	metaRepo jwtGenerate.MetaRepository
}

type JwtGenerateCommend interface {
	AccessToken(c context.Context, account, passWordSH256 string) (tokens jwtGenerate.LoginToken, err error)
	RefreshToken(c context.Context) (token jwtGenerate.RefreshToken, err error)
}

func NewjwtGenerateCommend(db jwtGenerate.MetaRepository) JwtGenerateCommend {
	return service{
		metaRepo: db,
	}
}

func (v service) AccessToken(c context.Context, account, passWordSH256 string) (tokens jwtGenerate.LoginToken, err error) {

	userInfo, err := v.metaRepo.LogInCheck(c, account, passWordSH256)
	if err != nil {
		return jwtGenerate.LoginToken{}, err
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":  userInfo.UserId,
		"name": userInfo.UserName,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(os.Getenv("JWT_ACCESS_SECRET_KEY"))

	refreshToken, err := v.RefreshToken(contextAddValue(c, userInfo))

	return jwtGenerate.LoginToken{
			AT: jwtGenerate.AccessToken(tokenString),
			RT: refreshToken},
		nil
}

func (v service) RefreshToken(c context.Context) (jwtGenerate.RefreshToken, error) {
	userId := c.Value("uid").(string)
	name := c.Value("name").(string)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	ApiToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":  userId,
		"name": name,
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	toString, err := ApiToken.SignedString(os.Getenv("JWT_SECRET_KEY"))

	if err != nil {
		return "", nil
	}

	return jwtGenerate.RefreshToken(toString), nil
}

func contextAddValue(ctx context.Context, u jwtGenerate.UserInfo) context.Context {
	ctx = context.WithValue(ctx, "uid", u.UserId)
	ctx = context.WithValue(ctx, "name", u.UserName)
	return ctx
}
