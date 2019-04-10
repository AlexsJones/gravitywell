package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
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
	input := &cloudformation.CreateStackInput{TemplateURL:
		aws.String("https://amazon-eks.s3-us-west-2.amazonaws.com/cloudformation/2019-02-11/amazon-eks-nodegroup.yaml"),
		StackName: aws.String("eks-deployment-stack")}
	input.RoleARN = aws.String("arn:aws:iam::976358786557:role/EKSRole")

	input.Parameters = []*cloudformation.Parameter{{
		ParameterKey:   aws.String("NodeGroupName"),
		ParameterValue: aws.String(""),
	},{
		ParameterKey:   aws.String("ClusterControlPlaneSecurityGroup"),
		ParameterValue: aws.String("sg-c76452b6"),
	},{
		ParameterKey:   aws.String("KeyName"),
		ParameterValue: aws.String(""),
	},{
		ParameterKey:   aws.String("NodeImageId"),
		ParameterValue: aws.String(""),
	},{
		ParameterKey:   aws.String("Subnets"),
		ParameterValue: aws.String("subnet-7d52a956, subnet-980d99d3, subnet-dec28384, subnet-644e241d"),
	},{
		ParameterKey:   aws.String("VpcId"),
		ParameterValue: aws.String(""),
	},{
		ParameterKey:   aws.String("ClusterName"),
		ParameterValue: aws.String(""),
	},
	}

	cso, err := svc.CreateStack(input)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	fmt.Println(cso.GoString())
}
