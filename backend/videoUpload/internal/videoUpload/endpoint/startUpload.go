/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-02 19:32:41
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-01-22 17:42:22
 * @FilePath: /videoUpload/internal/videoUpload/endpoint.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package endpoint

import (
	"context"
	"sideTube/videoUpload/internal/common/simpleKit/endpoint"
	"sideTube/videoUpload/internal/videoUpload/service"
)

func MakeStartUploadEndPoint(v service.VideoCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		res := data.(StartUploadRequest)
		videoId, err := v.StartUpload(c, res.TotalChunks)
		return StartUploadRespond{videoId}, err
	}
}

type StartUploadRequest struct {
	TotalChunks int32 `json:"totalChunks" validate:"required,alphanum"`
}

type StartUploadRespond struct {
	VideoId string `json:"videoId" `
}
