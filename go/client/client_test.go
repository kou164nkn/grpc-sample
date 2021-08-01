package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/kou164nkn/grpc-sample/go/deepthought"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockComputeClient struct{}

// mock func. implemented later
func (mc *MockComputeClient) Boot(ctx context.Context, in *deepthought.BootRequest, opts ...grpc.CallOption) (deepthought.Compute_BootClient, error) {
	return nil, nil
}

func (mc *MockComputeClient) Infer(ctx context.Context, in *deepthought.InferRequest, opts ...grpc.CallOption) (*deepthought.InferResponse, error) {
	switch in.Query {
	case "Life", "Universe", "Everything":
		return &deepthought.InferResponse{Answer: 42}, nil
	default:
		return nil, status.Error(codes.InvalidArgument, "Contemplate your query")
	}
}

func TestInfer(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		query        string
		expectResult *deepthought.InferResponse
		expectError  error
	}{
		"query_is_Life":       {query: "Life", expectResult: &deepthought.InferResponse{Answer: 42}, expectError: nil},
		"query_is_Universe":   {query: "Universe", expectResult: &deepthought.InferResponse{Answer: 42}, expectError: nil},
		"query_is_Everything": {query: "Everything", expectResult: &deepthought.InferResponse{Answer: 42}, expectError: nil},
		"query_is_foo":        {query: "foo", expectResult: nil, expectError: status.Error(codes.InvalidArgument, "Contemplate your query")},
	}

	for name, tt := range cases {
		mc := new(MockComputeClient)
		ctx := new(context.Context)

		t.Run(name, func(t *testing.T) {
			actualResult, actualError := Infer(mc, *ctx, tt.query)

			if !reflect.DeepEqual(tt.expectResult, actualResult) {
				t.Errorf("return InferResponse: want %v but got %v", tt.expectResult, actualResult)
			}

			if !reflect.DeepEqual(tt.expectError, actualError) {
				t.Errorf("return error: want %v but got %v", tt.expectError, actualError)
			}
		})
	}
}
