package core

import (
	"fmt"

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

	return &AWSProvider{sess}, nil
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
		// ImageId:      aws.String("ami-0828a1066dc750737"),
		ImageId:      aws.String("ami-07dd19a7900a1f049"),
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})
	if err != nil {
		return nil, err
	}
	instanceID := *runResult.Instances[0].InstanceId
	instanceType := *runResult.Instances[0].InstanceType
	availabilityZone := *runResult.Instances[0].Placement.AvailabilityZone
	fmt.Printf("instance: %+v", *runResult.Instances[0])
	fmt.Println("Created instance", instanceID)

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
	publicDnsName    string
	instanceOsUser   string
}

func (ec2n *EC2Node) GetHostName() string {
	return ec2n.publicDnsName + ":22"
}
