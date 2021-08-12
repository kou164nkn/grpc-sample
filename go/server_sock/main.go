package main

import (
	"fmt"
	"net"
	"os"

	"github.com/kou164nkn/grpc-sample/go/deepthought"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	socketName := "/tmp/sample.sock"
	err := os.Remove(socketName)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, err)
	}

	lis, err := net.Listen("unix", socketName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	serv := grpc.NewServer()

	hs := health.NewServer()
	healthgrpc.RegisterHealthServer(serv, hs)
	hs.Resume()

	deepthought.RegisterComputeServer(serv, &Server{})

	// return after being closed `l`, so there is no need close(l) in the main func.
	serv.Serve(lis)
}
