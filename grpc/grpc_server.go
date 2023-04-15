package grpc

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"sentinel-filter/service"
)

type server struct {
	// pb.UnimplementedBloomFilterServiceServer
	Filter *service.BloomFilter
}

func (s *server) Add(ctx context.Context, in *AddRequest) (*emptypb.Empty, error) {
	s.Filter.Add(in.Key)
	return &emptypb.Empty{}, nil
}

func (s *server) Check(ctx context.Context, in *CheckRequest) (*CheckResponse, error) {
	return &CheckResponse{Exists: s.Filter.Check(in.Key)}, nil
}
