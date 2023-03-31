/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-09 18:09:55
 * @FilePath: /ChannelStudio/internal/ChannelStudio/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"sideTube/ChannelStudio/internal/ChannelStudio"
	"sideTube/ChannelStudio/internal/ChannelStudio/service"
	"sideTube/ChannelStudio/internal/common/simpleKit/endpoint"
)

func MakeGetPublicVideoListEndPoint(v service.VideoStudioCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req := data.(GetPublicVideoListRequest)
		res, count, err := v.GetPublicVideoList(c, req.UserId, req.Skip, req.Limit)
		return GetPublicVideoListRespond{res, count}, err
	}
}

type GetPublicVideoListRequest struct {
	UserId string `json:"user_id" validate:"required,alphanum"`
	Skip   int64  `json:"skip" validate:"gte=0,lte=100000"`
	Limit  int64  `json:"limit" validate:"required,gte=5,lte=100"`
}

type GetPublicVideoListRespond struct {
	List  []ChannelStudio.VideoMeta `json:"list" `
	Count int64                     `json:"count,optimize"`
}
