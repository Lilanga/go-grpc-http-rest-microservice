package rest

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	v1 "github.com/Lilanga/go-grpc-http-rest-microservice/pkg/api/v1"
	"github.com/Lilanga/go-grpc-http-rest-microservice/pkg/logger"
	"github.com/Lilanga/go-grpc-http-rest-microservice/pkg/protocol/rest/middleware"
)

const (
	// relativePath is path to dist folder from server executable
	relativePath = "."
)

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, relativePath+"/todo-service.swagger.json")
}

// RunServer runs HTTP/REST gateway
func RunServer(ctx context.Context, grpcHost, grpcPort, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := v1.RegisterTodoServiceHandlerFromEndpoint(ctx, mux, grpcHost+":"+grpcPort, opts); err != nil {
		logger.Log.Fatal("failed to start http gateway", zap.String("reason", err.Error()))
	}

	// Serve the swagger-ui and swagger file
	// need to add swagger middleware to serve the files like logger
	smux := http.NewServeMux()
	smux.Handle("/", mux)
	smux.HandleFunc("/swagger.json", serveSwagger)
	fs := http.FileServer(http.Dir(relativePath + "/"))
	smux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", fs))

	srv := &http.Server{
		Addr: ":" + httpPort,
		Handler: middleware.AddRequestID(
			middleware.AddLogger(logger.Log, smux)),
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Log.Warn("shutting down http server...")

			srv.Shutdown(ctx)

			<-ctx.Done()
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	logger.Log.Info("Starting HTTP gateway...")
	return srv.ListenAndServe()
}
