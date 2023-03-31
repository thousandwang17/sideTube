/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-02 16:14:54
 * @FilePath: /picture/internal/picture/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sideTube/picture/internal/common/simpleKit/endpoint"
	"sideTube/picture/internal/picture/service"
)

var ErrRequestFailed = errors.New("Only allow to requset mpd, video")

type picture struct {
	E error
}

func (v picture) Error() string {
	return v.E.Error()
}

func (v picture) StatusCode() int {
	return http.StatusUnprocessableEntity
}

func MakeGetpictureEndPoint(v service.PictureCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req, _ := data.(GetPictureRequest)
		res, err := v.GetVideoPicture(c, req.FileName)
		return GetPictureRespond{res}, err
	}
}

type GetPictureRequest struct {
	FileName string `validate:"required,fileName,png"`
}

type GetPictureRespond struct {
	Data io.ReadCloser
}
