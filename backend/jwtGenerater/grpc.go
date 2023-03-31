/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-24 21:04:32
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-28 13:55:34
 * @FilePath: /jwtGenerater/grpc.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	pb "sideTube/jwtGenerate/internal/jwtGenerate/pb"
	"sideTube/jwtGenerate/internal/jwtGenerate/service"

	"google.golang.org/grpc"

	transportgrpc "sideTube/jwtGenerate/internal/jwtGenerate/transport/grpc"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	refreshToken kitgrpc.Handler
	pb.UnimplementedJwtTokenServer
}

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer.
func startGrpcServer(svc service.JwtGenerateCommend) *grpc.Server {

	grpcListener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", os.Getenv("GRPC_ADDRESS"), os.Getenv("GRPC_PORT")))
	if err != nil {
		log.Println("transport", "gRPC", "during", "Listen", "err", err)
		os.Exit(1)
	}


	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	pb.RegisterJwtTokenServer(baseServer, newGRPCServer(svc))

	go func() {
		log.Println("grpc server is staring")
		if err := 	baseServer.Serve(grpcListener); err != nil  {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return baseServer
}

func newGRPCServer(svc service.JwtGenerateCommend) pb.JwtTokenServer {
	return &grpcServer{
		refreshToken: transportgrpc.RefreshTokenRegister(svc),
	}
}

func (s *grpcServer) RefreshToken(ctx context.Context, req *pb.Requset) (*pb.Respond, error) {
	_, rep, err := s.refreshToken.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.Respond), nil
}
