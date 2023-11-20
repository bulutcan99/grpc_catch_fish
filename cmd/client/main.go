package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	pb "github.com/bulutcan99/grpc_weather/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

func main() {
	serverAddr := flag.String("server_addr", "localhost:8080", "The server address in the format of host:port")
	flag.Parse()
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})
	fmt.Printf("Client is running on port %v\n", *serverAddr)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, *serverAddr, opts...)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	client := pb.NewWeatherServiceClient(conn)
	req := &pb.RequestRegister{
		Username: "bcgocer",
		Name:     "Bulut Can",
		Email:    "bcgocer@gmail.com",
		Password: "1453",
		City:     "Istanbul",
	}

	res, err := client.Register(context.Background(), req)
	if err != nil {
		fmt.Errorf("error while calling Register RPC: %v", err)
	}

	fmt.Println(&res)
}
