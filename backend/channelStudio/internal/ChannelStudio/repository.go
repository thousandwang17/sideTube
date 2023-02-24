/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-02 17:36:49
 * @FilePath: /ChannelStudio/internal/ChannelStudio/repository.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package ChannelStudio

import (
	"context"
)

type MetaRepository interface {
	UpdateVideoMeta(c context.Context, userId string, data VideoEditMeta) error
	UpdateVideoState(c context.Context, userId string, videoId string, state uint8) error

	GetVideoList(c context.Context, userId string, skip int64, length int64) ([]VideoMeta, error)
	GetVideoCount(c context.Context, userId string) (int64, error)
}
