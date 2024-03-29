/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-31 17:44:01
 * @FilePath: /ChannelStudio/internal/ChannelStudio/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"image"
	"io"
	"sideTube/ChannelStudio/internal/ChannelStudio"
	"sideTube/ChannelStudio/internal/ChannelStudio/service"
	"sideTube/ChannelStudio/internal/common/simpleKit/endpoint"
)

func MakeEditVideoMetaEndPoint(v service.VideoStudioCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req := data.(EditVideoMetaRequest)
		err := v.EditVideoMeta(c, ChannelStudio.VideoEditMeta{
			Id:          req.VideoId,
			Title:       req.Title,
			Description: req.Desc,
			Picture:     req.PictureRS,
			Extension:   req.Extension,
			Config:      req.Config,
		})
		return EditVideoMetaRespond{}, err
	}
}

type EditVideoMetaRequest struct {
	VideoId string `json:"video_id" validate:"required"`

	Title     string `json:"title" validate:"required,max=100"`
	Desc      string `json:"desc" validate:"required,max=5000"`
	PictureRS io.ReadSeekCloser
	Extension string `validate:"required_with=PictureRS,omitempty,oneof=png jpeg"`

	Config image.Config
}

type EditVideoMetaRespond struct{}
