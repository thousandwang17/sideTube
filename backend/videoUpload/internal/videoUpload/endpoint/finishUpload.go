/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:27:39
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-01-18 19:40:43
 * @FilePath: /videoUpload/internal/videoUpload/endpoint/finishUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package endpoint

import (
	"context"
	"sideTube/videoUpload/internal/common/simpleKit/endpoint"
	"sideTube/videoUpload/internal/videoUpload/service"
)

func MakeFinsihUploadEndPoint(v service.VideoCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req := data.(FinsihUploadRequest)
		err := v.FinishUpload(c, req.VideoId)
		return FinsihUploadRespond{}, err
	}
}

type FinsihUploadRequest struct {
	VideoId string `json:"video_id" validate:"required"`
}

type FinsihUploadRespond struct {
}
