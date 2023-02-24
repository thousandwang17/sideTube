/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-01-25 22:36:24
 * @FilePath: /ChannelStudio/internal/ChannelStudio/service/service.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"context"
	"errors"
	"sideTube/ChannelStudio/internal/ChannelStudio"
)

var (
	ErrVaild = errors.New("args vaild failed")
)

const (
	private = iota
	publish
)

type service struct {
	metaRepo ChannelStudio.MetaRepository
}

type VideoStudioCommend interface {
	EditVideoMeta(c context.Context, data ChannelStudio.VideoEditMeta) error
	EditVideoPublicState(c context.Context, videoId string, state uint8) error
	GetVideoList(c context.Context, skip int64, length int64) (videos []ChannelStudio.VideoMeta, count int64, err error)
}

func NewVideoStudioCommend(db ChannelStudio.MetaRepository) VideoStudioCommend {
	return &service{
		metaRepo: db,
	}
}

// get all videos of the owner
func (v service) GetVideoList(c context.Context, skip int64, length int64) (data []ChannelStudio.VideoMeta, count int64, err error) {
	userId := c.Value("uid").(string)

	if skip == 0 {
		count, err = v.metaRepo.GetVideoCount(c, userId)
		if err != nil {
			return nil, 0, err
		}
	}

	data, err = v.metaRepo.GetVideoList(c, userId, skip, length)

	if err != nil {
		return nil, 0, err
	}

	return data, count, nil
}

// edit the video meta by owner, including title and desc
func (v service) EditVideoMeta(c context.Context, data ChannelStudio.VideoEditMeta) error {
	userId := c.Value("uid").(string)

	if err := v.metaRepo.UpdateVideoMeta(c, userId, data); err != nil {
		return err
	}

	return nil
}

// edit the video publish state  by owner
func (v service) EditVideoPublicState(c context.Context, videoId string, state uint8) error {
	userId := c.Value("uid").(string)

	if err := v.metaRepo.UpdateVideoState(c, videoId, userId, state); err != nil {
		return err
	}
	return nil
}
