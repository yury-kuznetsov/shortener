package main

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/yury-kuznetsov/shortener/api/pb"
	"github.com/yury-kuznetsov/shortener/internal/grpcsrv"
	"github.com/yury-kuznetsov/shortener/internal/uricoder"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func startGrpcServer(coder *uricoder.Coder, wg *sync.WaitGroup) (*grpc.Server, net.Listener, error) {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		return nil, nil, err
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))
	pb.RegisterServiceServer(server, grpcsrv.NewCoderServer(coder))

	go func() {
		defer wg.Done()
		if err := server.Serve(listen); err != nil {
			log.Fatalf("gRPC server Serve: %v", err)
		}
	}()

	return server, listen, nil
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}
	// псевдоавторизация (0 - гость, 1 - авторизованный)
	userID := 0
	tokens := md.Get("token")
	if len(tokens) > 0 && tokens[0] == "SECRET_KEY" {
		userID = 1
	}
	newCtx := context.WithValue(ctx, grpcsrv.KeyUserID, userID)

	return handler(newCtx, req)
}
