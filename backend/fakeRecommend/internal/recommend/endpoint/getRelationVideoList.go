/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-03 15:47:59
 * @FilePath: /recommend/internal/recommend/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"sideTube/recommend/internal/common/simpleKit/endpoint"
	"sideTube/recommend/internal/recommend"
	"sideTube/recommend/internal/recommend/service"
)

func MakeGetRelationVideoListEndPoint(v service.RecommendCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req := data.(GetRelationVideoListRequest)
		res, err := v.GetRelationVideoList(c, req.VideoId, req.Skip, req.Limit)
		return GetRelationVideoListRespond{res}, err
	}
}

type GetRelationVideoListRequest struct {
	VideoId string `json:"video_id" validate:"required,alphanum"`
	Skip    int64  `json:"skip" validate:"gte=0,lte=100000"`
	Limit   int64  `json:"limit" validate:"required,gte=5,lte=100"`
}

type GetRelationVideoListRespond struct {
	List []recommend.VideoMeta `json:"list" `
}
