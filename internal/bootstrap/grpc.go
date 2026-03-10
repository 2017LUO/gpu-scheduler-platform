package bootstrap

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

func NewGRPCServer(opts ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(opts...)
}

func NewGRPCListener(addr string) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("listen grpc on %s: %w", addr, err)
	}
	return ln, nil
}
