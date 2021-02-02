package core

// A Job is a unit of computation that is launched locally on some cluster
type Job struct {
	Script             string // Absolute path from local machine
	PythonDependencies map[string]string
	PythonVersion      string
	// TODO: other dependencies
}
