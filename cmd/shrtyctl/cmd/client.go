package cmd

import (
	"git.kono.sh/bkono/shrty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// NewClient creates a new ShrtyClient for the given server address.
func NewClient(serverAddr string) (shrty.ShrtyClient, *grpc.ClientConn) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		grpclog.Fatalf("failed to dial: %v", err)
	}

	client := shrty.NewShrtyClient(conn)
	return client, conn
}
