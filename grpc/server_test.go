package grpc

import (
	"context"
	"testing"

	pb "github.com/go-eyas/toolkit/grpc/example"
	"google.golang.org/grpc"
)

type RouteGuide struct{}

func (RouteGuide) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	panic("no impleament")
}
func (RouteGuide) ListFeatures(rectangle *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
	panic("no impleament")
}
func (RouteGuide) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
	panic("no impleament")
}
func (RouteGuide) RouteChat(stream pb.RouteGuide_RouteChatServer) error {
	panic("no impleament")
}

func TestServer(t *testing.T) {
	_, err := Server(&ServerConfig{
		Addr: ":6060",
		Register: func(ser *grpc.Server) {
			pb.RegisterRouteGuideServer(ser, &RouteGuide{})
		},
	})

	if err != nil {
		panic(err)
	}

}
