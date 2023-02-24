/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-07 13:02:33
 * @FilePath: /VideoStreaming/internal/VideoStreaming/repository.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package VideoStreaming

import (
	"context"
	"io"
)

type MetaRepository interface {
	GetVideoMetaById(c context.Context, userId string) (VideoMeta, error)
}

type VideoRepository interface {
	GetVideo(c context.Context, videoId string, pointer int64, length int64) (data io.ReadCloser, err error)
	GetMpdFile(c context.Context, videoId string) (data io.ReadCloser, err error)
}
