/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-02 19:32:41
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-26 18:27:16
 * @FilePath: /videoUpload/internal/videoUpload/endpoint.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package endpoint

import (
	"context"
	"io"
	"sideTube/videoUpload/internal/common/simpleKit/endpoint"
	"sideTube/videoUpload/internal/videoUpload/service"
)

func MakeUploadPartEndPoint(v service.VideoCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req := data.(UploadPartRequest)
		err := v.UploadPart(c, req.VideoId, req.PartId, req.Streaming)
		return UploadPartRespond{}, err
	}
}

type UploadPartRequest struct {
	Streaming io.ReadSeeker `validate:"required"`
	VideoId   string        `validate:"required,alphanum"`
	PartId    int64         `validate:"required,gte=1,lte=100000"`
}

type UploadPartRespond struct{}
