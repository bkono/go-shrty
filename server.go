package shrty

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
)

// ShrtygRPCServer implements the gRPC ShrtyService.
type ShrtygRPCServer struct {
	s ShortenedURLService
}

// NewgRPCServer factories a new gRPC server instance.
func NewgRPCServer(s ShortenedURLService) *ShrtygRPCServer {
	return &ShrtygRPCServer{s}
}

// Shorten creates a ShortUrl for a given ShortenRequest.
func (ss *ShrtygRPCServer) Shorten(ctx context.Context, sr *ShortenRequest) (*ShortenResponse, error) {
	grpclog.Printf("attempting to shorten %v", sr)
	su, err := ss.s.Shorten(sr.URL)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Error while building the shrt url")
	}

	resp := &ShortenResponse{su.ShrtURL, su.Token}
	return resp, nil
}

// Expand expands a token into its original url.
func (ss *ShrtygRPCServer) Expand(ctx context.Context, er *ExpandRequest) (*ExpandResponse, error) {
	grpclog.Printf("attempting to expand %v", er)
	u, err := ss.s.Expand(er.Token)
	grpclog.Printf("result %v, %v", u, err)
	if err == ErrTokenNotFound {
		return nil, grpc.Errorf(codes.NotFound, "Token not found")
	} else if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Error while expanding token")
	}

	resp := &ExpandResponse{u}
	return resp, nil
}
