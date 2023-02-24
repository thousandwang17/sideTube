/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-12 23:13:16
 * @FilePath: /videoMessage/internal/videoMessage/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"sideTube/videoMessage/internal/common/simpleKit/endpoint"
	"sideTube/videoMessage/internal/videoMessage/service"
)

func MakeReplyEndPoint(v service.VideoMessageCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req := data.(ReplyRequest)
		id, err := v.Reply(c, req.MessageId, req.Message)
		return ReplyRespond{id, err}, err
	}
}

type ReplyRequest struct {
	MessageId string `json:"message_id" validate:"required"`
	Message   string `json:"message" validate:"required"`
}

type ReplyRespond struct {
	ReplyId string `json:"ReplyId,omitempty" `
	Err     error  `json:"error,omitempty"`
}
