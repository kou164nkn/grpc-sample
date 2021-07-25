package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kou164nkn/grpc-sample/go/deepthought"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	err := subMain()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func subMain() error {
	if len(os.Args) != 2 {
		return errors.New("usage: client HOST:PORT")
	}
	addr := os.Args[1]

	// Connect in plaintext by specifying grpc.WithInsecure().
	// Don't actually do it because it's a security issue.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create the RPC client fron conn.
	// Since gRPC uses HTTP/2 streams, multi clients can use the same `conn`.
	// Also, multiple RPC client methods can be invoked simultaneously.
	cc := deepthought.NewComputeClient(conn)

	// Cancel `Boot` from client after 2.5 seconds
	ctx, cancel := context.WithCancel(context.Background())
	go func(cancel func()) {
		time.Sleep(2500 * time.Millisecond)
		cancel()
	}(cancel)

	stream, err := cc.Boot(ctx, &deepthought.BootRequest{})
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			// `io.EOF` indicates successful completeion of the stream.
			if err == io.EOF {
				break
			}
			// the package `status` translates between error and gRPC status.
			// `status.Code` retrieves the status code of gRPC.
			if status.Code(err) == codes.Canceled {
				break
			}
			return fmt.Errorf("receiving boot response: %w", err)
		}
		fmt.Printf("Boot: %s\n", resp.Message)
	}
	return nil
}
