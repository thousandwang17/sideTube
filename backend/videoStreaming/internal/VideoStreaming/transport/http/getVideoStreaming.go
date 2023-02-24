/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-23 16:16:04
 * @FilePath: /VideoStreaming/internal/VideoStreaming/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"sideTube/VideoStreaming/internal/VideoStreaming"
	videoEnpoint "sideTube/VideoStreaming/internal/VideoStreaming/endpoint"
	"sideTube/VideoStreaming/internal/middleware"
	"strconv"
	"strings"

	"context"
	"sideTube/VideoStreaming/internal/VideoStreaming/service"
	"sideTube/VideoStreaming/internal/common/simpleKit/endpoint"
	httptransport "sideTube/VideoStreaming/internal/common/simpleKit/httpTransport"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

var ErrorFileRange = errors.New("invalid Range header")

type fileRangeError struct{}

func (f fileRangeError) Error() string {
	return ErrorFileRange.Error()
}

func (f fileRangeError) StatusCode() int {
	return http.StatusRequestedRangeNotSatisfiable
}

func GetVideoStreamingRegister(svc service.VideoStreamingCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeGetVideoStreamingEndPoint(svc)
	ep = middleware.ValidMiddleWare(v)(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeGetVideoStreamingRequest,
		encodeGetVideoStreamingResponse,
	)
}

func decodeGetVideoStreamingRequest(_ context.Context, r *http.Request) (interface{}, error) {

	videoName := mux.Vars(r)["fileName"]

	ext := filepath.Ext(videoName)

	if ext == ".mpd" {
		return videoEnpoint.GetVideoMpdRequest{FileName: videoName}, nil
	} else {
		rangeHeader := r.Header.Get("Range")
		if rangeHeader == "" {
			return nil, fileRangeError{}
		}

		parts := strings.Split(rangeHeader, "=")
		if len(parts) != 2 || parts[0] != "bytes" {
			return nil, fileRangeError{}
		}

		ranges := strings.Split(parts[1], "-")
		if len(ranges) != 2 {
			return nil, fileRangeError{}
		}

		start, err := strconv.ParseInt(ranges[0], 10, 64)
		if err != nil {
			return nil, fileRangeError{}
		}

		end, err := strconv.ParseInt(ranges[1], 10, 64)
		if err != nil {
			end = -1
		}

		return videoEnpoint.GetVideoStreamingRequest{
			Id:    videoName,
			Start: start,
			End:   end,
		}, nil
	}

}

func encodeGetVideoStreamingResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(videoEnpoint.GetVideoStreamingRespond)
	defer resp.Data.Close()
	if resp.Type == VideoStreaming.TypeMpd {
		w.Header().Add("Content-Type", "application/dash+xml")
	} else {
		w.Header().Add("Content-Type", "arraybuffer")
	}
	_, err := io.Copy(w, resp.Data)
	return err
}
