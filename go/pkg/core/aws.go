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

type EC2Node struct {
	instanceID string
}

func (awsp *AWSProvider) GetNode() (Node, error) {
	// Create EC2 service client
	svc := ec2.New(awsp.sess)

	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		ImageId:      aws.String("ami-0828a1066dc750737"),
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})
	if err != nil {
		return nil, err
	}
	instanceID := *runResult.Instances[0].InstanceId
	fmt.Println("Created instance", instanceID)

	return &EC2Node{instanceID}, nil
}

func (ec2n *EC2Node) RunJob(job *Job) error {
	return nil
}

func (ec2n *EC2Node) ID() (string, error) {
	return ec2n.instanceID, nil
}

func (awsp *AWSProvider) DestroyNode(node Node) error {
	// Create EC2 service client
	svc := ec2.New(awsp.sess)
	instanceID, err := node.ID()
	if err != nil {
		return err
	}

	runResult, err := svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})
	if err != nil {
		fmt.Printf("error: ", err)
		return err
	}
	fmt.Printf("results: ", runResult)

	return nil

}
