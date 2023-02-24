/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-07 15:14:47
 * @FilePath: /VideoStreaming/internal/VideoStreaming/service/service.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"context"
	"errors"
	"io"
	"sideTube/VideoStreaming/internal/VideoStreaming"
)

var (
	ErrVaild = errors.New("args vaild failed")
)

const (
	private = iota
	publish
)

type service struct {
	metaRepo  VideoStreaming.MetaRepository
	videoRepo VideoStreaming.VideoRepository
}

type VideoStreamingCommend interface {
	GetVideoMeta(c context.Context, videoID string) (data VideoStreaming.VideoMeta, err error)
	GetVideoStreaming(c context.Context, videoID string, pointer int64, length int64) (data VideoStreaming.Video, err error)
	GetMpdByVideoId(c context.Context, videoID string) (data io.ReadCloser, err error)
}

func NewVideoStreamingCommend(db VideoStreaming.MetaRepository, fileReop VideoStreaming.VideoRepository) VideoStreamingCommend {
	return &service{
		metaRepo:  db,
		videoRepo: fileReop,
	}
}

func (v service) GetVideoMeta(c context.Context, videoID string) (data VideoStreaming.VideoMeta, err error) {

	data, err = v.metaRepo.GetVideoMetaById(c, videoID)

	if err != nil {
		return data, err
	}

	return data, nil
}

// get video streaming of video or audio , not including audio data
func (v service) GetVideoStreaming(c context.Context, videoID string, start int64, end int64) (streaming VideoStreaming.Video, err error) {

	data, err := v.videoRepo.GetVideo(c, videoID, start, end)

	if err != nil {
		return VideoStreaming.Video{Data: nil}, err
	}
	return VideoStreaming.Video{Data: data}, nil
}

// get  mpd file of video
func (v service) GetMpdByVideoId(c context.Context, videoID string) (data io.ReadCloser, err error) {
	return v.videoRepo.GetMpdFile(c, videoID)
}
