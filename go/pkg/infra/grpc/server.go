package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/latchai/latch/pkg/core"
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

	aws, err := core.NewAWSProvider()
	if err != nil {
		log.Fatal(err)
	}

	cr, err := core.NewAWSCommandRunner()
	if err != nil {
		log.Fatal(err)
	}

	pyVersionA := strings.Split(req.Job.PythonVersion, ".")
	pyVersion := pyVersionA[0] + "." + pyVersionA[1]

	core.RunJob(aws, cr, &core.Job{PythonDependencies: req.Job.PythonPackages, PythonVersion: pyVersion})

	return &pb.LaunchJobReply{Success: true}, nil
}

func formatYAML(instance, numWorkers, depen string) string {
	template := strings.Replace(template, "$ARG_INSTANCE", instance, 1)
	template = strings.Replace(template, "$ARG_MAX_WORKERS", numWorkers, 1)
	template = strings.Replace(template, "$ARG_DEP_LIST", depen, 1)
	return template

}

func depenMapToYAML(depen map[string]string) string {
	list := ""
	for k, v := range depen {
		list += fmt.Sprintf("\n  - pip install %s==%s", k, v)
	}
	return strings.TrimPrefix(list, "\n")
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
