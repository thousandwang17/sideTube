/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-31 20:58:37
 * @FilePath: /ChannelStudio/internal/ChannelStudio/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"context"
	"encoding/json"
	"errors"
	"image"
	"log"
	"net/http"
	videoEnpoint "sideTube/ChannelStudio/internal/ChannelStudio/endpoint"
	"sideTube/ChannelStudio/internal/ChannelStudio/service"
	"sideTube/ChannelStudio/internal/common/simpleKit/endpoint"
	httptransport "sideTube/ChannelStudio/internal/common/simpleKit/httpTransport"
	"sideTube/ChannelStudio/internal/middleware"

	"github.com/go-playground/validator"

	_ "image/jpeg"
	_ "image/png"
)

func EditVideoMetaRegister(svc service.VideoStudioCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeEditVideoMetaEndPoint(svc)
	ep = middleware.ValidMiddleWare(v)(ep)
	ep = middleware.JwtMiddleWare()(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeEditVideoMetaRequest,
		encodeEditVideoMetaResponse,
		httptransport.NewServerBefore(middleware.JwtServerBerore()),
	)
}

func decodeEditVideoMetaRequest(_ context.Context, r *http.Request) (interface{}, error) {

	video_id := r.FormValue("video_id")
	title := r.FormValue("title")
	desc := r.FormValue("desc")
	format := ""
	pngFile, header, err := r.FormFile("png")

	config := image.Config{}
	if err != nil && err != http.ErrMissingFile {
		pngFile = nil
	} else if header != nil {
		contentType := header.Header.Get("Content-Type")
		if contentType != "image/png" && contentType != "image/jpeg" {
			return nil, errors.New("invalid file type")
		}

		// Decode the uploaded file to check its format
		config, format, err = image.DecodeConfig(pngFile)

		if err != nil {
			log.Println("error decoding file", err)
			return nil, errors.New("error decoding file")
		}
		pngFile.Seek(0, 0)

		if format != "png" && format != "jpeg" {
			return nil, errors.New("invalid file format")
		}
	}
	return videoEnpoint.EditVideoMetaRequest{
		VideoId:   video_id,
		Title:     title,
		Desc:      desc,
		PictureRS: pngFile,
		Extension: format,
		Config:    config,
	}, nil
}

func encodeEditVideoMetaResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
