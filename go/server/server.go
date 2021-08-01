package main

import (
	"context"
	"time"

	"github.com/kou164nkn/grpc-sample/go/deepthought"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Type to implement ComputeServer interface
type Server struct {
	// A mechanism to prevent build errors
	// when RPCs are added to the proto file in the future
	// and the interface is extended.
	deepthought.UnimplementedComputeServer
}

// Verify at compile time that the interface has been implemented
var _ deepthought.ComputeServer = &Server{}

/*
func (s *Server) Boot(req *deepthought.BootRequest, stream deepthought.Compute_BootServer) error {
	panic("not implemented")
}

func (s *Server) Infer(ctx context.Context, req *deepthought.InferRequest) (*deepthought.InferResponse, error) {
	panic("not implemented")
}
*/

func (s *Server) Boot(req *deepthought.BootRequest, stream deepthought.Compute_BootServer) error {
	for {
		select {
		// Exit when the client cancels request.
		case <-stream.Context().Done():
			return nil
		// otherwise, wait 1 second and send data.
		case <-time.After(1 * time.Second):
		}

		if err := stream.Send(&deepthought.BootResponse{
			Message: "I THINK THEREFORE I AM.",
		}); err != nil {
			return err
		}
	}
}

func (s *Server) Infer(ctx context.Context, req *deepthought.InferRequest) (*deepthought.InferResponse, error) {
	switch req.Query {
	case "Life", "Universe", "Everything":
	default:
		// Use the predefined codes as a basis since gRPC defines commonly used error codes.
		// https://grpc.github.io/grpc/core/md_doc_statuscodes.html
		return nil, status.Error(codes.InvalidArgument, "Contemplate your query")
	}

	// check the client specifies timeout.
	deadline, ok := ctx.Deadline()

	// answer if not specified or if there is enough time.
	if !ok || time.Until(deadline) > 750*time.Millisecond {
		time.Sleep(750 * time.Millisecond)
		return &deepthought.InferResponse{
			Answer: 42,
			// Description: []string{"I checked it"},
		}, nil
	}

	// return DEADLINE_EXCEED (code 4) error if there is not enough time.
	return nil, status.Error(codes.DeadlineExceeded, "It would take longer")
}
