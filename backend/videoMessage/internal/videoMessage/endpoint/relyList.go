/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-10 18:14:10
 * @FilePath: /videoMessage/internal/videoMessage/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"sideTube/videoMessage/internal/common/simpleKit/endpoint"
	"sideTube/videoMessage/internal/videoMessage"
	"sideTube/videoMessage/internal/videoMessage/service"
)

func MakeReplyListEndPoint(v service.VideoMessageCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req := data.(ReplyListRequest)
		messages, err := v.ReplyList(c, req.MessageId, req.Skip, req.Limit)
		return ReplyListRespond{List: messages, Err: err}, err
	}
}

type ReplyListRequest struct {
	MessageId string `json:"message_id" validate:"required"`
	Skip      int64  `json:"skip" validate:"gte=0"`
	Limit     int64  `json:"limit" validate:"required,gte=5,lte=100"`
}

type ReplyListRespond struct {
	List []videoMessage.VideoMessageReplyMeta `json:"list,omitempty" `
	Err  error                                `json:"error,omitempty" `
}
