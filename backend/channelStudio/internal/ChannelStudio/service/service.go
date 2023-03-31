/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-31 21:24:44
 * @FilePath: /ChannelStudio/internal/ChannelStudio/service/service.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"context"
	"sideTube/ChannelStudio/internal/ChannelStudio"
)

const (
	minPicturWidth = 720

	maxPicturRatio = 0.75
	minPicturRatio = 0.5

	private = iota
	publish
)

type service struct {
	metaRepo  ChannelStudio.MetaRepository
	queueRepo ChannelStudio.MessagQqueue
	fileRepo  ChannelStudio.FileReop
}

type VideoStudioCommend interface {
	EditVideoMeta(c context.Context, data ChannelStudio.VideoEditMeta) error
	EditVideoPublicState(c context.Context, videoId string, state int8) error
	GetVideoList(c context.Context, skip int64, length int64) (videos []ChannelStudio.VideoMeta, count int64, err error)
	GetPublicVideoList(c context.Context, userId string, skip int64, length int64) (data []ChannelStudio.VideoMeta, count int64, err error)
}

func NewVideoStudioCommend(db ChannelStudio.MetaRepository, queue ChannelStudio.MessagQqueue, fileRepo ChannelStudio.FileReop) VideoStudioCommend {
	return &service{
		metaRepo:  db,
		queueRepo: queue,
		fileRepo:  fileRepo,
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

	data, err = v.metaRepo.GetVideoList(c, userId, skip, length, false)

	if err != nil {
		return nil, 0, err
	}

	return data, count, nil
}

// edit the video meta by owner, including title and desc
func (v service) EditVideoMeta(c context.Context, data ChannelStudio.VideoEditMeta) error {
	userId := c.Value("uid").(string)

	meta, err := v.metaRepo.GetVideoById(c, data.Id, userId)
	if err != nil {
		return err
	}

	pictureName := ""
	if data.Picture != nil {
		if data.Config.Width < minPicturWidth {
			return serviceErr{ErrNotReachMinWidth}
		}

		hw := float32(data.Config.Height) / float32(data.Config.Width)
		if hw < minPicturRatio || hw > maxPicturRatio {
			return serviceErr{ErrHigthWidthRatio}
		}

		// TODO : convert to WebP format
		pictureName, err = v.fileRepo.SetVideoPicture(c, data.Id, data.Extension, data.Picture)
		if err != nil {
			return err
		}
	}

	if err := v.metaRepo.UpdateVideoMeta(c, userId, data, pictureName); err != nil {
		return err
	}

	// if title has been change , we should regnerate serach index of the video
	if meta.Title != data.Title {
		if err := v.queueRepo.NotifySearchEngine(c, data.Id, userId); err != nil {
			return err
		}
	}

	return nil
}

// edit the video publish state  by owner
func (v service) EditVideoPublicState(c context.Context, videoId string, state int8) error {
	userId := c.Value("uid").(string)

	rep, err := v.metaRepo.GetVideoById(c, videoId, userId)
	if err != nil {
		return err
	}

	if rep.Title == "" || rep.Description == "" {
		return serviceErr{ErrSetPermission}
	}

	// already updated
	if rep.Permission == state {
		return nil
	}

	if err := v.metaRepo.UpdateVideoState(c, videoId, userId, state); err != nil {
		return err
	}

	if err := v.queueRepo.NotifySearchEngine(c, videoId, userId); err != nil {
		return err
	}

	return nil
}

// Get Public videos of the owner for other viewing .
// If in a production project, it should count the number of users' videos asynchronously, and the time complexity of getting the count should be O(1)
// another simple is use catch repo like redis or memcached
func (v service) GetPublicVideoList(c context.Context, userId string, skip int64, length int64) (data []ChannelStudio.VideoMeta, count int64, err error) {

	if skip == 0 {
		count, err = v.metaRepo.GetPublicVideoCount(c, userId)
		if err != nil {
			return nil, 0, err
		}
	}

	data, err = v.metaRepo.GetVideoList(c, userId, skip, length, true)

	if err != nil {
		return nil, 0, err
	}

	return data, count, nil
}
