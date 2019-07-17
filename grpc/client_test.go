package grpc

import (
	"testing"
	"context"

	pb "github.com/go-eyas/toolkit/grpc/example"
)

func TestClient(t *testing.T) {
	conn, err := NewClient(&ClientConfig{
		Addr: ":6060",
	})
	if err != nil {
		panic(err)
	}

	client := pb.NewRouteGuideClient(conn)

	res, err := client.GetFeature(context.TODO(), &pb.Point{})
	if err != nil {
		panic(err)
	}
	t.Logf("grpc response: %+v", res)
}
