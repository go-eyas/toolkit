package grpc

import "google.golang.org/grpc"

func NewClient(conf *ClientConfig) (*grpc.ClientConn, error) {
	if conf.Opts == nil {
		conf.Opts = make([]grpc.DialOption, 0)
	}
	conn, err := grpc.Dial(conf.Addr, conf.Opts...)

	if err != nil {
		return nil, err
	}
	conf.Register(conn)

	return conn, nil
}
