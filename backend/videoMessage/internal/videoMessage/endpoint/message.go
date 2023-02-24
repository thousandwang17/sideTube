/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-08 17:51:11
 * @FilePath: /videoMessage/internal/videoMessage/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"sideTube/videoMessage/internal/common/simpleKit/endpoint"
	"sideTube/videoMessage/internal/videoMessage/service"
)

func MakeMessageEndPoint(v service.VideoMessageCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req := data.(MessageRequest)
		id, err := v.Message(c, req.VideoId, req.Message)
		return MessageRespond{id, err}, err
	}
}

type MessageRequest struct {
	VideoId string `json:"video_id" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type MessageRespond struct {
	MessageId string `json:"MessageId,omitempty" validate:"required"`
	Err       error  `json:"error,omitempty"`
}
