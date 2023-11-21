package server

import (
	"context"
	"fmt"
	pb "github.com/bulutcan99/grpc_weather/proto"
)

type GrpcServer struct {
	pb.UnimplementedWeatherServiceServer
}

func (s *GrpcServer) Register(ctx context.Context, req *pb.RequestRegister) (*pb.ResponseRegister, error) {
	if req.Username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	return &pb.ResponseRegister{
		Username: req.Username,
		Name:     req.Name,
		Success:  true,
	}, nil
}
