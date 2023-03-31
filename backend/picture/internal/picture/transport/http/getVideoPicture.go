/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-31 21:32:32
 * @FilePath: /picture/internal/picture/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"errors"
	"io"
	"net/http"
	"sideTube/picture/internal/middleware"
	pictureEndpoint "sideTube/picture/internal/picture/endpoint"

	"context"
	"sideTube/picture/internal/common/simpleKit/endpoint"
	httptransport "sideTube/picture/internal/common/simpleKit/httpTransport"

	"sideTube/picture/internal/picture/service"

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

func GetpictureRegister(svc service.PictureCommend, v *validator.Validate) *httptransport.HttpTransport {
	var ep endpoint.EndPoint
	ep = pictureEndpoint.MakeGetpictureEndPoint(svc)
	ep = middleware.ValidMiddleWare(v)(ep)

	return httptransport.NewHttpTransport(
		ep,
		decodeGetpictureRequest,
		encodeGetpictureResponse,
	)
}

func decodeGetpictureRequest(_ context.Context, r *http.Request) (interface{}, error) {
	videoName := mux.Vars(r)["fileName"]
	return pictureEndpoint.GetPictureRequest{FileName: videoName}, nil
}

func encodeGetpictureResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(pictureEndpoint.GetPictureRespond)
	defer resp.Data.Close()
	w.Header().Add("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=600")
	_, err := io.Copy(w, resp.Data)
	return err
}
