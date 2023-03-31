/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-30 21:20:41
 * @FilePath: /user/internal/user/service/service.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"sideTube/user/internal/user"
	"sideTube/user/thrid-part/JwtGererater"
	"time"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
)

var (
	ErrVaild = errors.New("args vaild failed")
)

const (
	private = iota
	publish
)

type service struct {
	metaRepo        user.MetaRepository
	refreshTokenSrv JwtGererater.JwtTokenClient
}

type UserCommend interface {
	Login(ctx context.Context, email, password string) (user.LoginToken, error)
	Register(ctx context.Context, email, password, name string) error
	GetPublicUserInfo(c context.Context, user_id string) (user.PublicUserInfo, error)
	History(ctx context.Context, skip, length int64) ([]user.VideoMeta, error)
	LogOut(_ context.Context)
}

func NewUserCommend(db user.MetaRepository, srv JwtGererater.JwtTokenClient) UserCommend {
	return service{
		metaRepo:        db,
		refreshTokenSrv: srv,
	}
}

func (v service) Login(c context.Context, email, password string) (tokens user.LoginToken, err error) {

	hasher := sha256.New()
	hasher.Write([]byte(password))
	passwordHash := hex.EncodeToString(hasher.Sum(nil))

	userInfo, err := v.metaRepo.LogInCheck(c, email, passwordHash)
	if err != nil {
		return user.LoginToken{}, err
	}

	return v.accessToken(c, userInfo)
}

func (v service) LogOut(_ context.Context) {
	return
}

func (v service) Register(c context.Context, email, password, name string) error {

	hasher := sha256.New()
	hasher.Write([]byte(password))
	passwordHash := hex.EncodeToString(hasher.Sum(nil))

	err := v.metaRepo.Register(c, email, passwordHash, name)
	if err != nil {
		return err
	}

	return nil
}

func (v service) accessToken(c context.Context, userInfo user.UserInfo) (user.LoginToken, error) {

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":  userInfo.UserId,
		"name": userInfo.UserName,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET_KEY")))

	if err != nil {
		log.Println("access token err : ", err)
		return user.LoginToken{}, errors.New("Faild to generate access token ")
	}

	// add grpc request header via context
	md := metadata.New(map[string]string{"AccessToken": tokenString})
	ctx := metadata.NewOutgoingContext(c, md)

	resp, err := v.refreshTokenSrv.RefreshToken(ctx,
		&JwtGererater.Requset{AccessToken: tokenString})

	if err != nil {
		log.Println("refresh token err : ", err)
		return user.LoginToken{}, errors.New("Faild to refresh access token ")
	}

	return user.LoginToken{
			AT: user.AccessToken(tokenString),
			RT: user.RefreshToken(resp.Data.GetRefreshToken())},
		nil
}

func (v service) GetPublicUserInfo(c context.Context, user_id string) (user.PublicUserInfo, error) {

	userInfo, err := v.metaRepo.GetPublicUserInfo(c, user_id)
	if err != nil {
		return user.PublicUserInfo{}, err
	}

	return userInfo, nil
}

func (v service) History(ctx context.Context, skip, length int64) ([]user.VideoMeta, error) {
	userId := ctx.Value("uid").(string)

	data, err := v.metaRepo.GetHistoryList(ctx, userId, skip, length)

	if err != nil {
		return nil, err
	}

	if data == nil {
		return []user.VideoMeta{}, nil
	}

	return data, nil
}
