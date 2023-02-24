/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-01-29 20:16:34
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

func MakeEditVideoMetaEndPoint(v service.VideoStudioCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req := data.(EditVideoMetaRequest)
		err := v.EditVideoMeta(c, ChannelStudio.VideoEditMeta{
			Id:          req.VideoId,
			Title:       req.Titie,
			Description: req.Desc,
		})
		return EditVideoMetaRespond{}, err
	}
}

type EditVideoMetaRequest struct {
	VideoId string `json:"video_id" validate:"required"`
	Titie   string `json:"title" validate:"required,max=100"`
	Desc    string `json:"desc" validate:"required,max=5000"`
}

type EditVideoMetaRespond struct{}
