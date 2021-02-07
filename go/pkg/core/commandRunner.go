package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
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
					time.Sleep(1 * time.Second)
					continue
				}
			} else {
				fmt.Println(err.Error())
				return err
			}
		}
		state := result.Reservations[0].Instances[0].State
		dns := result.Reservations[0].Instances[0].PublicDnsName
		fmt.Printf("\ndescribe result: %+v", state)
		if *state.Name == "running" {
			fmt.Printf("\ndns name: %+v", *dns)
			node.publicDnsName = *dns
			fmt.Printf("\nstate: %+v", *state.Name)
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Send ssh key to instance.
	svc2 := ec2instanceconnect.New(awscr.sess)
	sshBytes, err := ioutil.ReadFile(os.Getenv("HOME") + "/.ssh/id_rsa.pub")
	if err != nil {
		log.Fatal(err)
	}
	input2 := &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: aws.String(node.availabilityZone),
		InstanceId:       aws.String(node.instanceID),
		InstanceOSUser:   aws.String(node.instanceOsUser),
		SSHPublicKey:     aws.String(strings.TrimSpace(string(sshBytes))),
	}

	fmt.Println("\nsending public key")
	sendSSHRes, err := svc2.SendSSHPublicKey(input2)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ec2instanceconnect.ErrCodeAuthException:
				fmt.Println(ec2instanceconnect.ErrCodeAuthException, aerr.Error())
			case ec2instanceconnect.ErrCodeInvalidArgsException:
				fmt.Println(ec2instanceconnect.ErrCodeInvalidArgsException, aerr.Error())
			case ec2instanceconnect.ErrCodeServiceException:
				fmt.Println(ec2instanceconnect.ErrCodeServiceException, aerr.Error())
			case ec2instanceconnect.ErrCodeThrottlingException:
				fmt.Println(ec2instanceconnect.ErrCodeThrottlingException, aerr.Error())
			case ec2instanceconnect.ErrCodeEC2InstanceNotFoundException:
				fmt.Println(ec2instanceconnect.ErrCodeEC2InstanceNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
			return aerr
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return err
		}
	}
	fmt.Printf("\nsshResults %+v", sendSSHRes)

	key := os.Getenv("HOME") + "/.ssh/id_rsa"
	runner := exec.Command("ssh", "-oStrictHostKeyChecking=no", fmt.Sprintf("%s@%s", node.instanceOsUser, node.publicDnsName), fmt.Sprintf("-i%s", key))
	fmt.Printf("\ncommand: %+v", runner)
	runner.Stdin = os.Stdin
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	runner.Run()

	// pk, _ := ioutil.ReadFile(os.Getenv("HOME") + "/.ssh/id_rsa")
	// signer, err := ssh.ParsePrivateKey(pk)
	// if err != nil {
	// 	panic(err)
	// }

	// // Communicate with instance via ssh.
	// cfg := &ssh.ClientConfig{
	// 	User: node.instanceOsUser,
	// 	Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	// }
	// cfg.SetDefaults()

	// fmt.Printf("\n:user %+v", node.instanceOsUser)
	// fmt.Printf("\nhostname: %+v", node.GetHostName())

	// fmt.Println("\nsleeping")
	// time.Sleep(5 * time.Second)
	// fmt.Println("\nattempting connection")
	// client, err := ssh.Dial("tcp", node.GetHostName(), cfg)
	// if err != nil {
	// 	fmt.Printf("\nssh error: %+v", err)
	// 	return err
	// }
	// defer client.Close()

	// session, err := client.NewSession()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("we have a session! %+v", session)

	// fmt.Println("launched")
	// runCommand("cat /dev/urandom", conn)
	// time.Sleep(2 * time.Second)

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
