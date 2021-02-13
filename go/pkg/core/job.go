package core

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	pb "github.com/latchai/latch/pkg/infra/servicepb"
)

func NewJob(req *pb.LaunchJobRequest) (*Job, error) {

	pyVersionA := strings.Split(req.Job.PythonVersion, ".")
	pyVersion := pyVersionA[0] + "." + pyVersionA[1]

	tempScript, err := scriptToTemp(req.Job.Script)
	if err != nil {
		return nil, err
	}

	return &Job{Script: tempScript, PythonDependencies: req.Job.PythonPackages, PythonVersion: pyVersion}, nil

}

// Convert a script called from interpretor to a "ligand-less" script in a `tmp` path.
func scriptToTemp(scriptPath string) (string, error) {

	// stream edit
	unParsed, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		return "", err
	}

	_, fileName := filepath.Split(scriptPath)

	// Replace:
	//	* `import ligand`
	//	* `ligand.init()`
	importParse := strings.Replace(string(unParsed), "import ligand", "", 1)
	methodParse := strings.Replace(importParse, "ligand.init()", "", 1)

	// Construct new temp. file deferring to native os. Write parsed script.
	f, err := ioutil.TempFile("", fileName)
	if err != nil {
		return "", err
	}

	// Write parsed
	err = ioutil.WriteFile(f.Name(), []byte(methodParse), 0666)
	if err != nil {
		return "", err
	}

	return f.Name(), nil

}

// A Job is a unit of computation that is launched locally on some cluster
type Job struct {
	Script             string // Absolute path from local machine
	PythonDependencies map[string]string
	PythonVersion      string
	// TODO: other dependencies
}
