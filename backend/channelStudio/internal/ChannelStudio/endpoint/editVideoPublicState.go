/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-01-30 17:04:30
 * @FilePath: /ChannelStudio/internal/ChannelStudio/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"sideTube/ChannelStudio/internal/ChannelStudio/service"
	"sideTube/ChannelStudio/internal/common/simpleKit/endpoint"
)

func MakeEditVideoPublicStateEndPoint(v service.VideoStudioCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req := data.(EditVideoPublicStateRequest)
		err := v.EditVideoPublicState(c, req.VideoId, req.State)
		return EditVideoPublicStateRespond{}, err
	}
}

type EditVideoPublicStateRequest struct {
	VideoId string `json:"video_id" validate:"required"`
	State   uint8  `json:"state" validate:"gte=0,lte=1"`
}

type EditVideoPublicStateRespond struct{}
