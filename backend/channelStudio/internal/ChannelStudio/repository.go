/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-31 21:11:08
 * @FilePath: /ChannelStudio/internal/ChannelStudio/repository.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package ChannelStudio

import (
	"context"
	"io"
)

type MetaRepository interface {
	UpdateVideoMeta(c context.Context, userId string, data VideoEditMeta, pictureName string) error
	UpdateVideoState(c context.Context, userId string, videoId string, state int8) error

	GetVideoList(ctx context.Context, userId string, skip int64, length int64, public bool) ([]VideoMeta, error)
	GetVideoCount(c context.Context, userId string) (int64, error)
	GetPublicVideoCount(c context.Context, userId string) (int64, error)
	GetVideoById(ctx context.Context, videoId, userId string) (VideoMeta, error)
}

type MessagQqueue interface {
	NotifySearchEngine(c context.Context, videoId, userId string) (err error)
}

type FileReop interface {
	SetVideoPicture(c context.Context, videoId, extension string, data io.ReadSeekCloser) (imageName string, err error)
}
