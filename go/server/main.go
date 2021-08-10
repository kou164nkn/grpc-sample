package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/kou164nkn/grpc-sample/go/deepthought"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
)

const portNumber = 13333

func main() {
	kep := keepalive.EnforcementPolicy{
		MinTime: 10 * time.Second,
	}
	serv := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kep))

	hs := health.NewServer()
	healthgrpc.RegisterHealthServer(serv, hs)
	hs.Resume()

	deepthought.RegisterComputeServer(serv, &Server{})

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		fmt.Println("failed to listen", err)
		os.Exit(1)
	}

	// return after being closed `l`, so there is no need close(l) in the main func.
	serv.Serve(l)
}
