/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-23 20:38:06
 * @FilePath: /recommend/internal/recommend/service/service.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"context"
	"errors"
	"math/rand"
	"sideTube/recommend/internal/recommend"
	"time"
)

var (
	ErrVaild = errors.New("args vaild failed")
)

const (
	private = iota
	publish
)

type service struct {
	metaRepo recommend.MetaRepository
}

type RecommendCommend interface {
	GetHomeVideoList(c context.Context, skip int64, length int64) (data []recommend.VideoMeta, err error)
	GetRelationVideoList(c context.Context, video_id string, skip int64, length int64) (data []recommend.VideoMeta, err error)
}

func NewrecommendCommend(db recommend.MetaRepository) RecommendCommend {
	return &service{
		metaRepo: db,
	}
}

func (v service) GetHomeVideoList(c context.Context, skip int64, length int64) (data []recommend.VideoMeta, err error) {
	userId, ok := c.Value("uid").(string)

	// if uid is empty , that mean requester is a guest
	if !ok {
		userId = ""
	}

	data, err = v.metaRepo.GetHomeVideoList(c, userId, skip, length)

	if err != nil {
		return nil, err
	}

	if data == nil {
		data = []recommend.VideoMeta{}
	}

	return data, nil
}

func (v service) GetRelationVideoList(c context.Context, video_id string, skip int64, length int64) (data []recommend.VideoMeta, err error) {
	data, err = v.metaRepo.GetRelationVideoList(c, video_id, skip, length)

	if err != nil {
		return nil, err
	}

	if data == nil {
		data = []recommend.VideoMeta{}
	}

	// set the random seed based on the current time
	rand.Seed(time.Now().UnixNano())

	// shuffle the array using the Fisher-Yates algorithm
	for i := len(data) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}

	return data, nil
}
