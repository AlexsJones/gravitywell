package aws

import (
	"errors"
	"fmt"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/fatih/color"
	"github.com/jpillora/backoff"
	"github.com/satori/go.uuid"
	"time"
)

func (awsp *AWSProvider) Create(clusterp kinds.ProviderCluster) error {

	ec := eks.New(awsp.AWSClient)

	//SecurityGroupID
	var securityGroupId []*string
	for _, c := range clusterp.SecurityGroupID {
		securityGroupId = append(securityGroupId, aws.String(c))
	}
	//SubnetID
	var subnetId []*string
	for _, c := range clusterp.SubnetID {
		subnetId = append(subnetId, aws.String(c))
	}
	u, err := uuid.NewV4()
	if err != nil {
		return err
	}

	input := &eks.CreateClusterInput{
		ClientRequestToken: aws.String(u.String()),
		Name:               aws.String(clusterp.Name),
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SecurityGroupIds: securityGroupId,
			SubnetIds:        subnetId,
		},
		RoleArn: aws.String(clusterp.RoleARN),
		Version: aws.String(clusterp.KubernetesVersion),
	}

	_, err = ec.CreateCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceInUseException:
				fmt.Println(eks.ErrCodeResourceInUseException, aerr.Error())
			case eks.ErrCodeResourceLimitExceededException:
				fmt.Println(eks.ErrCodeResourceLimitExceededException, aerr.Error())
			case eks.ErrCodeInvalidParameterException:
				fmt.Println(eks.ErrCodeInvalidParameterException, aerr.Error())
			case eks.ErrCodeClientException:
				fmt.Println(eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				fmt.Println(eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				fmt.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
			case eks.ErrCodeUnsupportedAvailabilityZoneException:
				fmt.Println(eks.ErrCodeUnsupportedAvailabilityZoneException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err
	}

	b := &backoff.Backoff{
		Min:    30 * time.Second,
		Max:    time.Second * 60,
		Jitter: true,
	}
	for {
		color.Blue(fmt.Sprintf("Started cluster build"))

		do, err := ec.DescribeCluster(&eks.DescribeClusterInput{Name: aws.String(clusterp.Name)})
		if err != nil {
			return err
		}
		fmt.Printf("%v", do)

		if *do.Cluster.Status == "ACTIVE" {
			color.Green("Cluster running")

			//Build node groups

			return nil
		}

		time.Sleep(b.Duration())
		if b.Attempt() >= 20 {
			return errors.New("max retry attempts hit")
		}
	}
}
