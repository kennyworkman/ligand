package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func NewAWSProvider() (*AWSProvider, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")})
	if err != nil {
		return nil, err
	}

	// Check if aws CLI installed.
	_, err = exec.LookPath("aws")
	if err != nil {
		log.Println(error.Error(err))
		// TODO: better error handling
		// We don't necessarily need the cli to use go sdk...
		downloadAWS()
	}

	// Check that aws CLI is configured
	_, err = ioutil.ReadFile(os.Getenv("HOME") + "/.aws/config")
	if err != nil {
		log.Println(error.Error(err))
		// TODO: better error handling
		configureAWS()
	}
	_, err = ioutil.ReadFile(os.Getenv("HOME") + "/.aws/credentials")
	if err != nil {
		log.Println(error.Error(err))
		// TODO: better error handling
		configureAWS()
	}

	// if not

	// setup provider
	// prompt for cli download if no .aws
	// creat latch .pem file

	return &AWSProvider{sess}, nil
}

// Works on osx.
func downloadAWS() {
	cmd := exec.Command("curl", "https://awscli.amazonaws.com/AWSCLIV2.pkg", "-oAWSCLIV2.pkg")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = exec.Command("sudo", "installer", "-pkg", "AWSCLIV2.pkg", "-target", "/")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Contingent on aws cli binary present.
func configureAWS() {
	cmd := exec.Command("aws", "config")
	cmd.Run()
}

// concrete type of provider
type AWSProvider struct {
	sess *session.Session
}

func (awsp *AWSProvider) GetNode() (*EC2Node, error) {
	// Create EC2 service client
	svc := ec2.New(awsp.sess)

	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		// ImageId: aws.String("ami-0b1a80ce62c464a55"),
		// 20.04 w/o anything
		// ImageId:      aws.String("ami-07dd19a7900a1f049"),
		// ImageId:      aws.String("ami-0b1a80ce62c464a55"),
		ImageId:      aws.String("ami-0828a1066dc750737"),
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		KeyName:      aws.String("latch"),
	})
	if err != nil {
		return nil, err
	}
	instanceID := *runResult.Instances[0].InstanceId
	instanceType := *runResult.Instances[0].InstanceType
	availabilityZone := *runResult.Instances[0].Placement.AvailabilityZone
	fmt.Printf("instance: %+v", *runResult.Instances[0])
	fmt.Println("\nCreated instance", instanceID)

	return &EC2Node{instanceID: instanceID, instanceType: instanceType, availabilityZone: availabilityZone, instanceOsUser: "ubuntu"}, nil
}

func (awsp *AWSProvider) DestroyNode(node *EC2Node) error {

	// Create EC2 service client
	svc := ec2.New(awsp.sess)
	runResult, err := svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(node.instanceID)},
	})
	if err != nil {
		fmt.Printf("error: ", err)
		return err
	}
	fmt.Printf("\ndestroy results: ", runResult)

	return nil

}

type EC2Node struct {
	instanceID       string
	instanceType     string
	availabilityZone string
	publicIpAddress  string
	instanceOsUser   string
}

func (ec2n *EC2Node) GetHostName() string {
	return ec2n.publicIpAddress + ":22"
}
