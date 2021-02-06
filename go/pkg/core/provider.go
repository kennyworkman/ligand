package core

import (
	"log"
)

type NodeProvider interface {
	GetNode() (Node, error)
	DestroyNode(Node) error
}

type Node interface {
	/* TODO: add methods */
	// probably expose some methods for access
	RunJob(job *Job) error
	ID() (string, error)
}

// this is like the communication between computer <> script to run.
// there needs to be sense of the local dev environment in the model
// how do we both portforward and move things around?

// while the job is executing, there needs to be
//	* communication with a console stdout
//	* moving around of files

func RunJob(np NodeProvider, job *Job) {
	node, err := np.GetNode()
	if err != nil {
		log.Fatal(err)
	}

	// err = node.RunJob(job)
	// // err = node.RunJob(job)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err = np.DestroyNode(node)
	if err != nil {
		log.Fatal(err)
	}
}

// A Job is a unit of computation that is launched locally on some cluster
type Job struct {
	Script             string // Absolute path from local machine
	PythonDependencies map[string]string
	PythonVersion      string
	// TODO: other dependencies
}
