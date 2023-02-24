/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-10 20:28:19
 * @FilePath: /videoMessage/internal/videoMessage/repository.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package videoMessage

import (
	"context"
)

type MetaRepository interface {
	MessageList(c context.Context, videoID string, skip, length int64) (data []VideoMessageMeta, err error)
	Message(c context.Context, userID, userName, videoID, message string) (messageID string, err error)
	EditMessage(c context.Context, userID, messageID, message string) error
	DeleteMessage(c context.Context, userID, messageID string) error

	// respond
	ReplyList(c context.Context, messageID string, skip, length int64) (data []VideoMessageReplyMeta, err error)
	Reply(c context.Context, userID, userName, ReplyID, message string) (messageID string, err error)
	EditReply(c context.Context, userID, messageID, message string) error
	DeleteReply(c context.Context, userID, messageID string) error
}
