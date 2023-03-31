/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-02 21:03:30
 * @FilePath: /videoUpload/internal/videoUpload/repository.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package videoUpload

import (
	"context"
	"io"
)

type VideoRepository interface {
	CreateMultipartUpload(c context.Context, id string) (data VideoRepoMeta, err error)
	UploadPart(c context.Context, v VideoRepoMeta, file io.ReadSeeker) error
	CompleteMultipartUpload(c context.Context, id string) error
	AbortUpload(c context.Context, id string) error
}

type MetaRepository interface {
	Insert(c context.Context, userId, userName string) (id string, err error)
	UpdateState(c context.Context, videoId, userId string) error
	UndoUpdateState(c context.Context, videoId, userId string) error
	Remove(c context.Context, videoId, userId string) error
}

type MessagQqueue interface {
	NotifyMergeVideo(c context.Context, videoId, userId string) (err error)
}
