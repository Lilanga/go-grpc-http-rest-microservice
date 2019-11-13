package main

import (
	"fmt"
	"os"

	"github.com/Lilanga/go-grpc-http-rest-microservice/pkg/cmd"
)

func main() {
	if err := cmd.RunGRPCServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
