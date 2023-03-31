/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-30 16:21:17
 * @FilePath: /User/internal/User/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package endpoint

import (
	"context"
	"sideTube/user/internal/common/simpleKit/endpoint"
	"sideTube/user/internal/user"
	"sideTube/user/internal/user/service"
)

func MakeGetHistoryListEndPoint(v service.UserCommend) endpoint.EndPoint {
	return func(c context.Context, data interface{}) (interface{}, error) {
		req := data.(GetHistoryListRequest)
		res, err := v.History(c, req.Skip, req.Limit)
		return GetHistoryListRespond{res}, err
	}
}

type GetHistoryListRequest struct {
	Skip  int64 `json:"skip" validate:"gte=0,lte=1000"`
	Limit int64 `json:"limit" validate:"required,gte=5,lte=100"`
}

type GetHistoryListRespond struct {
	List []user.VideoMeta `json:"list" `
}
