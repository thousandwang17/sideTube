/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-02 19:32:41
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-26 18:26:11
 * @FilePath: /videoUpload/internal/videoUpload/endpoint.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"context"
	"encoding/json"
	"net/http"
	"sideTube/videoUpload/internal/common/simpleKit/endpoint"
	httptransport "sideTube/videoUpload/internal/common/simpleKit/httpTransport"
	"sideTube/videoUpload/internal/middleware"
	videoEnpoint "sideTube/videoUpload/internal/videoUpload/endpoint"
	"sideTube/videoUpload/internal/videoUpload/service"
	"strconv"

	"github.com/go-playground/validator"
)

func UploadPartRegister(svc service.VideoCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeUploadPartEndPoint(svc)
	ep = middleware.ValidMiddleWare(v)(ep)
	ep = middleware.JwtMiddleWare()(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeUploadPartRequest,
		encodeUploadPartResponse,
		httptransport.NewServerBefore(middleware.JwtServerBerore()),
	)
}

func decodeUploadPartRequest(_ context.Context, r *http.Request) (interface{}, error) {
	video_id := r.FormValue("video_id")
	part_id := r.FormValue("part_id")
	partId, _ := strconv.Atoi(part_id)

	videoChunkData, _, err := r.FormFile("streaming_data")
	if err != nil && err != http.ErrMissingFile {
		return nil, err
	}

	return videoEnpoint.UploadPartRequest{
		VideoId:   video_id,
		PartId:    int64(partId),
		Streaming: videoChunkData,
	}, nil
}

func encodeUploadPartResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
