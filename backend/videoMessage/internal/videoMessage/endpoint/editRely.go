/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-13 21:21:47
 * @FilePath: /videoMessage/internal/videoMessage/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"sideTube/videoMessage/internal/common/simpleKit/endpoint"
	"sideTube/videoMessage/internal/videoMessage/service"
)

func MakeEditReplyEndPoint(v service.VideoMessageCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req := data.(EditReplyRequest)
		err := v.EditReply(c, req.ReplyID, req.NewMessage)
		return EditReplyRespond{err}, err
	}
}

type EditReplyRequest struct {
	ReplyID    string `json:"reply_id" validate:"required"`
	NewMessage string `json:"Message" validate:"required"`
}

type EditReplyRespond struct {
	Err error `json:"error,omitempty"`
}
