package grpc

import (
	"net"

	"google.golang.org/grpc"
)

// NewServer 初始化 grpc 服务器
func NewServer(conf *ServerConfig) (*grpc.Server, error) {
	if conf.Opts == nil {
		conf.Opts = make([]grpc.ServerOption, 0)
	}

	lis, err := net.Listen("tcp", conf.Addr)
	if err != nil {
		return nil, err
	}
	ser := grpc.NewServer(conf.Opts...)

	conf.Register(ser)

	ser.Serve(lis)

	return ser, nil
}
