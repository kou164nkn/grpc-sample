package main

import (
	"fmt"
	"net"
	"os"

	"github.com/kou164nkn/grpc-sample/go/deepthought"
	"google.golang.org/grpc"
)

const portNumber = 13333

func main() {
	serv := grpc.NewServer()

	deepthought.RegisterComputeServer(serv, &Server{})

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		fmt.Println("failed to listen", err)
		os.Exit(1)
	}

	// return after being closed `l`, so there is no need close(l) in the main func.
	serv.Serve(l)
}
