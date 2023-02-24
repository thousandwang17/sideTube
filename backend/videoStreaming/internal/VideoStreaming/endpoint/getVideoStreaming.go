/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-23 16:13:41
 * @FilePath: /VideoStreaming/internal/VideoStreaming/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sideTube/VideoStreaming/internal/VideoStreaming"
	"sideTube/VideoStreaming/internal/VideoStreaming/service"
	"sideTube/VideoStreaming/internal/common/simpleKit/endpoint"
)

var ErrRequestFailed = errors.New("Only allow to requset mpd, video")

type videoStreaming struct {
	E error
}

func (v videoStreaming) Error() string {
	return v.E.Error()
}

func (v videoStreaming) StatusCode() int {
	return http.StatusUnprocessableEntity
}

func MakeGetVideoStreamingEndPoint(v service.VideoStreamingCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req, ok := data.(GetVideoStreamingRequest)
		if ok {
			res, err := v.GetVideoStreaming(c, req.Id, req.Start, req.End)
			return GetVideoStreamingRespond{res.Data, VideoStreaming.TypeVideo}, err
		} else if req, ok := data.(GetVideoMpdRequest); ok {
			res, err := v.GetMpdByVideoId(c, req.FileName)
			return GetVideoStreamingRespond{res, VideoStreaming.TypeMpd}, err
		}

		return nil, videoStreaming{ErrRequestFailed}
	}
}

type GetVideoMpdRequest struct {
	FileName string ` validate:"required,fileName"`
}

type GetVideoStreamingRequest struct {
	Id    string `json:"video_id" validate:"required,fileName"`
	Start int64  `validate:"gte=0"`
	End   int64  `validate:"required"`
}

type GetVideoStreamingRespond struct {
	Data io.ReadCloser
	Type string
}
