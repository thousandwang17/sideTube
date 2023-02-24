/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-03 15:40:10
 * @FilePath: /VideoStreaming/internal/VideoStreaming/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"sideTube/VideoStreaming/internal/VideoStreaming"
	"sideTube/VideoStreaming/internal/VideoStreaming/service"
	"sideTube/VideoStreaming/internal/common/simpleKit/endpoint"
)

func MakeGetVideoMetaEndPoint(v service.VideoStreamingCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req := data.(GetVideoMetaRequest)
		res, err := v.GetVideoMeta(c, req.Id)
		return GetVideoMetaRespond{res}, err
	}
}

type GetVideoMetaRequest struct {
	Id string `json:"video_id" validate:"required"`
}

type GetVideoMetaRespond struct {
	VideoStreaming.VideoMeta
}
