package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"

	v1 "github.com/Lilanga/go-grpc-http-rest-microservice/pkg/api/v1"
	"github.com/Lilanga/go-grpc-http-rest-microservice/pkg/logger"
	"github.com/Lilanga/go-grpc-http-rest-microservice/pkg/protocol/grpc/middleware"
	"google.golang.org/grpc"
)

// RunServer runs gRPC service to publish Todo service
func RunServer(ctx context.Context, v1API v1.TodoServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}
	opts = middleware.AddLogging(logger.Log, opts)

	// register service
	server := grpc.NewServer(opts...)
	v1.RegisterTodoServiceServer(server, v1API)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Log.Warn("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	logger.Log.Info("starting gRPC server...")
	return server.Serve(listen)
}
