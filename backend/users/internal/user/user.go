package user

import (
	"context"
	"time"

	"sideTube/users/internal/user/repository"
)

type userInfo struct {
	id          string // uuid
	name        string
	icon        string
	email       string
	instruction string
	createTime  time.Time
}

type user struct {
	userInfo

	subList         []userInfo
	repository      repository.Repository
	cacheRepository repository.Repository
}

type User interface {
	sub(c context.Context, userId string) (bool, error)
	showSubList(c context.Context, count, skip int) ([]userInfo, error)

	// unset()
}

func New() User {
	return user{}
}

func (u user) sub(c context.Context, userId string) (bool, error) {
	return false, nil
}

func (u user) showSubList(c context.Context, count, skip int) ([]userInfo, error) {
	return []userInfo{}, nil
}
