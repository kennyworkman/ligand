package core

import (
	"fmt"
	"log"
	"path/filepath"
)

type NodeProvider interface {
	GetNode() (*EC2Node, error)
	DestroyNode(*EC2Node) error
}

// this is like the communication between computer <> script to run.
// there needs to be sense of the local dev environment in the model
// how do we both portforward and move things around?

// while the job is executing, there needs to be
//	* communication with a console stdout
//	* moving around of files

func RunJob(np NodeProvider, cr CommandRunner, job *Job) {

	fmt.Println("Provisioning your node from AWS...")
	node, err := np.GetNode()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Waiting until AWS starts your node...")
	success, err := cr.WaitNodeRunning(node, &Console{})
	if err != nil {
		fmt.Printf("out chan: %+v", success)
		log.Fatal(err)
	}

	// Setup machine command
	fmt.Println("Setting up your node's OS and installing dependencies...")
	err = cr.WaitConnectionPossible(node, &Console{})
	if err != nil {
		log.Fatal(err)
	}

	err = cr.Run(setupCommand(job), node, &Console{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Copying your entry point to your node...")
	err = cr.RsyncUp(job.Script, node, &Console{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running your entrypoint:\n")
	_, sourceFile := filepath.Split(job.Script)
	err = cr.Run(fmt.Sprintf("python%s /home/ubuntu/%s", job.PythonVersion, sourceFile), node, &Console{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Destroying your node.")
	err = np.DestroyNode(node)
	if err != nil {
		log.Fatal(err)
	}
}

// Helper function for to activate correct env in AWS AMI. Also install
// additional python depen.
func setupCommand(job *Job) string {
	cmd := "sudo apt-get install software-properties-common"
	cmd += "\nsudo killall apt apt-get"
	cmd += "\nsudo add-apt-repository ppa:deadsnakes/ppa -y"
	cmd += "\nsudo apt-get update"
	cmd += fmt.Sprintf("\nsudo apt-get install python%s-dev -y", job.PythonVersion)
	cmd += fmt.Sprintf("\nsudo apt-get install python3-pip -y")
	cmd += fmt.Sprintf("\npython%s --version", job.PythonVersion)
	cmd += fmt.Sprintf("\nexport LATCH_REMOTE=true")
	for k, v := range job.PythonDependencies {
		if k == "ligand" {
			continue
		}
		cmd += fmt.Sprintf("\npython%s -m pip install %s==%s", job.PythonVersion, k, v)
	}
	cmd += fmt.Sprintf("\npython%s -m pip show six", job.PythonVersion)
	return cmd
}
