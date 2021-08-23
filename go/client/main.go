package main

import (
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: client HOST:PORT")
	}
	addr := os.Args[1]

	kp := keepalive.ClientParameters{
		Time: 1 * time.Minute,
	}

	// Connect in plaintext by specifying grpc.WithInsecure().
	// Don't actually do it because it's a security issue.
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kp))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer conn.Close()

	err = CallBoot(conn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = CallInfer(conn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
