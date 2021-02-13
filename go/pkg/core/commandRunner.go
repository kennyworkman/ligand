package core

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"golang.org/x/crypto/ssh"
)

type CommandRunner interface {
	WaitConnectionPossible(node *EC2Node, cons *Console) error
	WaitNodeRunning(node *EC2Node, cons *Console) (bool, error)
	Run(cmd string, node *EC2Node, cons *Console) error
	RsyncUp(source string, node *EC2Node, cons *Console) error
	RsyncDown() error
}

type AWSCommandRunner struct {
	sess *session.Session
}

// A CommandRunner can run shell script on a Node.
func NewAWSCommandRunner() (*AWSCommandRunner, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")})
	if err != nil {
		return nil, err
	}

	return &AWSCommandRunner{sess}, nil

}

func (awscr *AWSCommandRunner) Run(cmd string, node *EC2Node, cons *Console) error {

	key := os.Getenv("HOME") + "/.ssh/latch.pem"
	runner := exec.Command("ssh", "-oStrictHostKeyChecking=no", fmt.Sprintf("%s@%s", node.instanceOsUser, node.publicIpAddress), fmt.Sprintf("-i%s", key), "-y", cmd)

	runner.Stdin = os.Stdin
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	runner.Run()
	return nil

}

// Asynchronously attempts to form TCP handshake with given node. Returns when successful.
func (awscr *AWSCommandRunner) WaitConnectionPossible(node *EC2Node, cons *Console) error {

	success := make(chan bool, 1)
	listener, err := ioutil.TempFile("", "ligand_listener")
	if err != nil {
		return err
	}
	fmt.Println(listener.Name())

	// Spawns process to attempt connection
	go attemptConnect(listener, node)

	scanner := bufio.NewScanner(listener)
	for true {
		for scanner.Scan() {
			fmt.Println("SUCCESS: ", scanner.Text())
		}
	}

	select {
	case _ = <-success:
		return nil
	}

}

func attemptConnect(listener *os.File, node *EC2Node) {

	key := os.Getenv("HOME") + "/.ssh/latch.pem"
	runner := exec.Command("ssh", "-oStrictHostKeyChecking=no", fmt.Sprintf("%s@%s", node.instanceOsUser, node.publicIpAddress), fmt.Sprintf("-i%s", key), "-y", "echo 'ROSA PARKS'")

	runner.Stdin = os.Stdin
	runner.Stdout = listener
	runner.Stderr = os.Stderr
	runner.Run()

	time.Sleep(10 * time.Second)
	fmt.Printf("starting again...")
	go attemptConnect(listener, node)

}

func (awscr *AWSCommandRunner) WaitNodeRunning(node *EC2Node, cons *Console) (bool, error) {

	// Wait for instance status = 'running'.
	svc := ec2.New(awscr.sess)

	success := make(chan bool, 1)
	go pingAvailability(svc, node, success)

	select {
	case suc := <-success:
		return suc, nil
	}

}

// Concurrent routine to ping for instance availability.
func pingAvailability(svc *ec2.EC2, node *EC2Node, success chan bool) {

	fmt.Print(".")
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(node.instanceID)},
	}
	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(err.Error())
				time.Sleep(1 * time.Second)
				go pingAvailability(svc, node, success)
				return
			}
		} else {
			fmt.Println(err.Error())
			return
		}
	}

	state := result.Reservations[0].Instances[0].State
	ip := result.Reservations[0].Instances[0].PublicIpAddress
	if *state.Name == "running" {
		node.publicIpAddress = *ip
		success <- true
		return
	}
	time.Sleep(1 * time.Second)
	go pingAvailability(svc, node, success)
}

func (awscr *AWSCommandRunner) RsyncUp(source string, node *EC2Node, console *Console) error {

	for i := 0; i < 2; i++ {
		fmt.Println("connecting now...")
		key := os.Getenv("HOME") + "/.ssh/latch.pem"
		_, sourceFile := filepath.Split(source)
		target := "/home/ubuntu/" + sourceFile

		runner := exec.Command("rsync", "-avz", "-progress", "-e", fmt.Sprintf("ssh -i%s", key), source, fmt.Sprintf("%s@%s:%s", node.instanceOsUser, node.publicIpAddress, target), "-y")
		fmt.Printf("cmd: %+v", runner)

		runner.Stdin = os.Stdin
		runner.Stdout = os.Stdout
		runner.Stderr = os.Stderr
		runner.Run()
		time.Sleep(10 * time.Second)

	}
	return nil
}

func (awscr *AWSCommandRunner) RsyncDown() error {
	return nil
}

type DockerRunner struct {
}

// TODO: refactor below
// raw bytes
func getPublicKey() (string, error) {
	byteKey, err := ioutil.ReadFile("/Users/kenny/.ssh/id_rsa.pub")
	if err != nil {
		return "", err
	}
	return string(byteKey), nil
}

// implementation for ssh client
func publicKey(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

func runCommand(cmd string, conn *ssh.Client) {
	sess, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close()
	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stdout, sessStdOut)
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stderr, sessStderr)
	err = sess.Run(cmd) // eg., /usr/bin/whoami
	if err != nil {
		panic(err)
	}
}
