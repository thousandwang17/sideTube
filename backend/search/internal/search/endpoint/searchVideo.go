/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-19 17:01:47
 * @FilePath: /search/internal/search/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"sideTube/search/internal/common/simpleKit/endpoint"
	"sideTube/search/internal/search"
	"sideTube/search/internal/search/service"
)

func MakeSearchVideoEndPoint(v service.SearchCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req := data.(SearchVideoRequest)
		res, err := v.Serach(c, req.Query, req.Skip, req.Limit)
		return SearchVideoRespond{res}, err
	}
}

type SearchVideoRequest struct {
	Query string `json:"query" validate:"required,max=100"`
	Skip  int64  `json:"skip" validate:"gte=0,lte=100000"`
	Limit int64  `json:"limit" validate:"required,gte=5,lte=100"`
}

type SearchVideoRespond struct {
	List []search.VideoMeta `json:"list" `
}
