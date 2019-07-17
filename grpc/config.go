package grpc

import (
	"google.golang.org/grpc"
)

type Server *grpc.Server
type Client *grpc.ClientConn

type ServerConfig struct {
	Addr     string
	Opts     []grpc.ServerOption
	Register func(Server)
}

type ClientConfig struct {
	Addr     string
	Opts     []grpc.DialOption
	Register func(Client)
}
