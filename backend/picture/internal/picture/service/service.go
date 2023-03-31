/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-02 16:16:20
 * @FilePath: /picture/internal/picture/service/service.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"context"
	"errors"
	"io"
	"sideTube/picture/internal/picture"
)

var (
	ErrVaild = errors.New("args vaild failed")
)

const (
	private = iota
	publish
)

type service struct {
	videoRepo picture.PictureRepository
}

type PictureCommend interface {
	GetVideoPicture(c context.Context, videoID string) (data io.ReadCloser, err error)
}

func NewpictureCommend(fileReop picture.PictureRepository) PictureCommend {
	return &service{
		videoRepo: fileReop,
	}
}

// get  mpd file of video
func (v service) GetVideoPicture(c context.Context, videoID string) (data io.ReadCloser, err error) {
	return v.videoRepo.GetPicture(c, videoID)
}
