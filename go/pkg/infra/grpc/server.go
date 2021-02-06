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

	core.RunJob(aws, &core.Job{})
	// // cluster config file name
	// filename := "generated.yaml"
	// // pythonVersion := req.Job.PythonVersion
	// depen := depenMapToYAML(req.Job.PythonPackages)
	// template := formatYAML("g4dn.xlarge", strconv.Itoa(1), depen)
	// scriptPath := req.Job.Script
	// dir, file := filepath.Split(scriptPath)
	// rayConfig := filepath.Join(dir, filename)
	// remotePath := filepath.Join("/home/ray/", file)

	// err := ioutil.WriteFile(rayConfig, []byte(template), 0644)
	// if err != nil {
	// 	return &pb.LaunchJobReply{Success: false}, err
	// }

	// // 1. Construct cluster
	// console.Info("\n📡 Constructing your ephemeral cluster...")
	// cmd := exec.Command("ray", "up", rayConfig, "-y")
	// _, err = cmd.Output()
	// if err != nil {
	// 	return &pb.LaunchJobReply{Success: false}, err
	// }

	// input, err := ioutil.ReadFile(scriptPath)
	// if err != nil {
	// 	return &pb.LaunchJobReply{Success: false}, err
	// }

	// lines := strings.Split(string(input), "\n")

	// for i, line := range lines {
	// 	if strings.Contains(line, "import latch") || strings.Contains(line, "latch.init") {
	// 		lines[i] = ""
	// 	}
	// }

	// output := strings.Join(lines, "\n")
	// scriptPath = strings.Replace(scriptPath, ".py", "_remote.py", 1)
	// err = ioutil.WriteFile(scriptPath, []byte(output), 0644)
	// if err != nil {
	// 	return &pb.LaunchJobReply{Success: false}, err
	// }

	// // // 2. Send necessary files to cluster
	// console.Info("📂 Syncing necessary files to the cloud...")
	// cmd = exec.Command("ray", "rsync_up", filename, scriptPath, remotePath)
	// err = cmd.Run()
	// if err != nil {
	// 	return &pb.LaunchJobReply{Success: false}, err
	// }

	// // 3. Execute remote script
	// console.Info("🛠️ Executing your script:")
	// remoteCommand := fmt.Sprintf("sudo env \"PATH=$PATH\" python %s", remotePath)
	// cmd = exec.Command("ray", "exec", filename, remoteCommand)
	// cmd.Stdout = os.Stdout
	// // cmd.Stderr = os.Stderr
	// err = cmd.Run()
	// if err != nil {
	// 	return &pb.LaunchJobReply{Success: false}, err
	// }

	// // // 4. Tear down cluster
	// console.Info("\n🚜 Tearing down your ephermal cluster...")
	// cmd = exec.Command("ray", "down", filename, "-y")
	// err = cmd.Run()
	// if err != nil {
	// 	return &pb.LaunchJobReply{Success: false}, err
	// }

	// return nil
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
