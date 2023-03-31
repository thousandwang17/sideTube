/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-08 20:26:35
 * @FilePath: /VideoStreaming/internal/VideoStreaming/service/service.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sideTube/VideoStreaming/internal/VideoStreaming"
)

var (
	ErrVaild = errors.New("args vaild failed")
)

const (
	private = iota
	publish
)

type serviceError struct{}

func (_ serviceError) StatusCode() int {
	return http.StatusForbidden
}

func (_ serviceError) Error() string {
	return "video is private"
}

type service struct {
	metaRepo  VideoStreaming.MetaRepository
	videoRepo VideoStreaming.VideoRepository
}

type VideoStreamingCommend interface {
	GetVideoMeta(c context.Context, videoID string) (data VideoStreaming.VideoMeta, err error)
	GetVideoStreaming(c context.Context, videoID string, pointer int64, length int64) (data VideoStreaming.Video, err error)
	GetMpdByVideoId(c context.Context, videoID string) (data io.ReadCloser, err error)
	IncVideoViews(c context.Context, videoID string) error
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

	if data.Permission == 0 {
		return data, serviceError{}
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

// if in produciotn for count views, using Apache Kafka with Apache Spark or Apache Flink to control db pressure
func (v service) IncVideoViews(c context.Context, videoID string) (err error) {
	userId := c.Value("uid").(string)
	return v.metaRepo.IncVideoViews(c, videoID, userId)
}
