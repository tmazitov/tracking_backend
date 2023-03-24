package grpc

import (
	"context"
	"log"
	"net"

	"github.com/tmazitov/tracking_backend.git/internal/aaa/proto/api"
	"google.golang.org/grpc"
)

type Server struct{}

// CheckAuth is check an authorization of the user by the token
func (s *Server) CheckAuth(ctx context.Context, in *api.CheckRequest) (*api.CheckResponse, error) {
	return &api.CheckResponse{Result: true}, nil
}

func SetupServer(grpcPort string) {
	grpc := grpc.NewServer()
	srv := &Server{}
	api.RegisterAAAServer(grpc, srv)

	l, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalln(err)
		return
	}

	if err = grpc.Serve(l); err != nil {
		log.Fatalln(err)
		return
	}

	log.Printf("Success run grpc server on port %s\n", grpcPort)
}
