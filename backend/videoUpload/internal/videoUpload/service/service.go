/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-28 14:56:43
 * @FilePath: /videoUpload/internal/videoUpload/service/service.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"context"
	"errors"
	"io"
	"log"
	"sideTube/videoUpload/internal/videoUpload"
)

var (
	ErrVaild            = errors.New("args vaild failed")
	FaildToFinishUpload = errors.New("Faild To Finish Upload")
)

type service struct {
	metaRepo  videoUpload.MetaRepository
	fileRepo  videoUpload.VideoRepository
	queueRepo videoUpload.MessagQqueue
}

type VideoCommend interface {
	EditInfo() error
	StartUpload(c context.Context, chundCount int32) (key string, err error)
	UploadPart(c context.Context, videoId string, chunkId int64, data io.ReadSeeker) error
	FinishUpload(c context.Context, videoId string) error
	AbortUpload(c context.Context, videoId string) error
}

func NewVideoCommend(db videoUpload.MetaRepository, v videoUpload.VideoRepository, q videoUpload.MessagQqueue) VideoCommend {
	return &service{
		metaRepo:  db,
		fileRepo:  v,
		queueRepo: q,
	}
}

func (v service) EditInfo() error {
	return nil
}

// Init new video upload
func (v service) StartUpload(c context.Context, chundCount int32) (string, error) {

	userId := c.Value("uid").(string)
	name := c.Value("userName").(string)
	videoId, err := v.metaRepo.Insert(c, userId, name)
	if err != nil {
		return "", err
	}

	if videoId == "" {
		return "", errors.New("videoID create failed")
	}

	_, err = v.fileRepo.CreateMultipartUpload(c, videoId)

	if err != nil {
		v.metaRepo.Remove(c, videoId, userId)
		return "", err
	}

	return videoId, nil
}

func (v service) UploadPart(c context.Context, videoId string, chunkId int64, data io.ReadSeeker) error {

	if videoId == "" {
		return ErrVaild
	}

	if chunkId == 0 {
		return ErrVaild
	}

	if err := v.fileRepo.UploadPart(c,
		videoUpload.VideoRepoMeta{
			Id:     videoId,
			PartID: chunkId,
		}, data,
	); err != nil {
		return err
	}

	return nil
}

func (v service) FinishUpload(c context.Context, videoId string) error {
	userId := c.Value("uid").(string)

	if err := v.metaRepo.UpdateState(c, videoId, userId); err != nil {
		return err
	}

	if err := v.queueRepo.NotifyMergeVideo(c, videoId, userId); err != nil {
		// undo meta of video when message queue return error
		if err := v.metaRepo.UndoUpdateState(c, videoId, userId); err != nil {
			log.Println("failed to undo err: ", err)
		}
		log.Println("failed to Notify MergeVideo err: ", err)
		return FaildToFinishUpload
	}
	return nil
}

func (v service) AbortUpload(c context.Context, videoId string) error {
	userId := c.Value("uid").(string)

	if err := v.metaRepo.Remove(c, videoId, userId); err != nil {
		return err
	}

	if err := v.fileRepo.AbortUpload(c, videoId); err != nil {
		return err
	}

	return nil
}
