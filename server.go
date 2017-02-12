package shrty

import (
	"fmt"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
)

type shrtygRPCServer struct {
	s ShortenedURLService
}

// RunGRPCServer initializes and runs the an instance of ShrtyServer for gRPC handling.
func RunGRPCServer(s ShortenedURLService, port int) error {
	shrtygRPC := &shrtygRPCServer{s}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	RegisterShrtyServer(grpcServer, shrtygRPC)
	grpcServer.Serve(lis)
	return nil
}

func (ss *shrtygRPCServer) Shorten(ctx context.Context, sr *ShortenRequest) (*ShortenResponse, error) {
	grpclog.Printf("attempting to shorten %v", sr)
	su, err := ss.s.Shorten(sr.URL)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Error while building the shrt url")
	}

	resp := &ShortenResponse{su.ShrtURL, su.Token}
	return resp, nil
}

func (ss *shrtygRPCServer) Expand(ctx context.Context, er *ExpandRequest) (*ExpandResponse, error) {
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
