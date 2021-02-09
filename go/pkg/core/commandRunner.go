package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"golang.org/x/crypto/ssh"
)

type CommandRunner interface {
	Run(cmd string, node *EC2Node, cons *Console) error
	RsyncUp() error
	RsyncDown() error
}

type AWSCommandRunner struct {
	sess *session.Session
}

func NewAWSCommandRunner() (*AWSCommandRunner, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")})
	if err != nil {
		return nil, err
	}

	return &AWSCommandRunner{sess}, nil

}

func (awscr *AWSCommandRunner) Run(cmd string, node *EC2Node, cons *Console) error {

	// Wait for instance status = 'running'.
	svc := ec2.New(awscr.sess)
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(node.instanceID)},
	}

	// TODO: Hack, leverage channels/concurrency
	for i := 0; i < 100; i++ {
		result, err := svc.DescribeInstances(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(err.Error())
					fmt.Println("Trying again...")
					time.Sleep(1 * time.Second)
					continue
				}
			} else {
				fmt.Println(err.Error())
				return err
			}
		}
		state := result.Reservations[0].Instances[0].State
		fmt.Println(state)
		ip := result.Reservations[0].Instances[0].PublicIpAddress
		if *state.Name == "running" {
			fmt.Printf("\nip name: %+v", *ip)
			node.publicIpAddress = *ip
			break
		}
		time.Sleep(1 * time.Second)
	}

	for i := 0; i < 2; i++ {

		fmt.Printf("connecting")
		key := os.Getenv("HOME") + "/.ssh/latch.pem"
		runner := exec.Command("ssh", "-oStrictHostKeyChecking=no", fmt.Sprintf("%s@%s", node.instanceOsUser, node.publicIpAddress), fmt.Sprintf("-i%s", key), cmd)

		runner.Stdin = os.Stdin
		runner.Stdout = os.Stdout
		runner.Stderr = os.Stderr
		runner.Run()

		fmt.Printf("sleeping and trying again")
		time.Sleep(5 * time.Second)
	}

	return nil
}

func (awscr *AWSCommandRunner) RsyncUp() error {
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
