package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/kou164nkn/grpc-sample/go/deepthought"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CallBoot() error {
	dialer := &net.Dialer{}
	dialFunc := func(ctx context.Context, a string) (net.Conn, error) {
		return dialer.DialContext(ctx, "unix", a)
	}
	conn, err := grpc.Dial("/tmp/sample.sock", grpc.WithInsecure(), grpc.WithContextDialer(dialFunc))
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

func CallInfer() error {
	dialer := &net.Dialer{}
	dialFunc := func(ctx context.Context, a string) (net.Conn, error) {
		return dialer.DialContext(ctx, "unix", a)
	}
	conn, err := grpc.Dial("/tmp/sample.sock", grpc.WithInsecure(), grpc.WithContextDialer(dialFunc))
	if err != nil {
		return err
	}
	defer conn.Close()

	cc := deepthought.NewComputeClient(conn)

	// Set the deadline at 2 seconds from now
	// The Client would call Infer request twice
	shortDuration := 2000 * time.Millisecond
	deadline := time.Now().Add(shortDuration)

	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// gRPC server expects to receive the following messages
	queryMessages := [...]string{"Life", "Universe", "Everything"}

	for _, msg := range queryMessages {
		// resp, err := cc.Infer(ctx, &deepthought.InferRequest{Query: msg})
		resp, err := Infer(cc, ctx, msg)
		if err != nil {
			return nil
		}
		fmt.Printf("Infer: %s\n", resp.String())
		fmt.Printf("Infer Answer: %d\n", resp.GetAnswer())
		// fmt.Printf("Infer Description: %s\n", resp.GetDescription())
	}

	return nil
}

func Infer(cc deepthought.ComputeClient, ctx context.Context, msg string) (*deepthought.InferResponse, error) {
	return cc.Infer(ctx, &deepthought.InferRequest{Query: msg})
}
