package grpc

import (
	"context"
	"fmt"
	config_builder "github.com/bulutcan99/grpc_weather/pkg/config"
	pb "github.com/bulutcan99/grpc_weather/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
)

func StartGRPCServer() {
	grpcConn, err := config_builder.ConnectionURLBuilder("grpc")
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", grpcConn)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	pb.RegisterWeatherServiceServer(grpcServer, &GrpcServer{})

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			zap.S().Errorf("Failed to serve: %v\n", err)
		}
	}()
	zap.S().Info("gRPC server is running on port 50051")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	zap.S().Info("Shutting down the gRPC server...")
	grpcServer.GracefulStop()
	zap.S().Info("gRPC server stopped.")
}

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
