package grpc

import (
	"context"
	"fmt"
	"net"

	pb "github.com/latchai/latch/pkg/infra/servicepb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedDaemonServer
}

func (s *server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	return &pb.PingReply{Success: true}, nil
}

func (s *server) LaunchJob(ctx context.Context, req *pb.LaunchJobRequest) (*pb.LaunchJobReply, error) {
	return &pb.LaunchJobReply{Success: true}, nil
}

func ListenAndServe(socketPath string) error {

	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		return fmt.Errorf("Failed to open UNIX socket on %s: %w", socketPath, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDaemonServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil

}
