package grpc

import (
	"google.golang.org/grpc"
)

type Config struct {
	Addr string
}

type ServerConfig struct {
	Config
	Opts     []grpc.ServerOption
	Register func(*grpc.Server)
}

type ClientConfig struct {
	Config
	Opts     []grpc.DialOption
	Register func(*grpc.ClientConn)
}
