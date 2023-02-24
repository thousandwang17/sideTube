/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-24 21:38:58
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-24 22:28:19
 * @FilePath: /jwtGenerater/internal/jwtGenerate/transport/grpc/login.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package grpc

import (
	videoEnpoint "sideTube/jwtGenerate/internal/jwtGenerate/endpoint"
	JwtGererater "sideTube/jwtGenerate/internal/jwtGenerate/pb"
	"sideTube/jwtGenerate/internal/jwtGenerate/service"
	"sideTube/jwtGenerate/internal/middleware"

	"context"
	"sideTube/jwtGenerate/internal/common/simpleKit/endpoint"

	kitEndPoint "github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

func RefreshTokenRegister(svc service.JwtGenerateCommend) *grpctransport.Server {
	var ep endpoint.EndPoint
	ep = videoEnpoint.MakeRefreshTokenEndPoint(svc)
	ep = middleware.JwtMiddleWare()(ep)
	kitep := kitEndPoint.Endpoint(ep)

	return grpctransport.NewServer(
		kitep,
		decodeGRPCSumRequest,
		encodeGRPCSumResponse,
	)

}

func decodeGRPCSumRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	// req := grpcReq.(*JwtGererater.RequsetLogin)
	return videoEnpoint.RefreshTokenRequest{}, nil
}

func encodeGRPCSumResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(videoEnpoint.RefreshTokenRespond)
	return &JwtGererater.RespondLogin{}, nil
}
