package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
)

func main() {

	awsP := os.Getenv("AWS_DEFAULT_PROFILE")
	if awsP == "" {
		fmt.Printf("No default profile\n")
		os.Exit(1)
	}
	awsR := os.Getenv("AWS_DEFAULT_REGION")
	if awsR == "" {
		fmt.Printf("No default region\n")
		os.Exit(1)
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsR),
		Credentials: credentials.NewSharedCredentials("", awsP),
	})
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	svc := cloudformation.New(sess)

	//https://github.com/weaveworks/eksctl/blob/2bf40d7646fb2faa0297c21a3efff4a3a8a1d282/pkg/cfn/manager/api.go
	input := &cloudformation.CreateStackInput{TemplateURL: aws.String("https://amazon-eks.s3-us-west-2.amazonaws.com/cloudformation/2019-02-11/amazon-eks-nodegroup.yaml"),
		StackName: aws.String("eks-deployment-stack")}



}
