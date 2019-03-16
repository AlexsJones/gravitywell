package aws

import (
	"errors"
	"fmt"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/fatih/color"
	"github.com/jpillora/backoff"
	"time"
)
func (awsp *AWSProvider)Delete(clusterp kinds.ProviderCluster) error {

	ec := eks.New(awsp.AWSClient)

	_, err := ec.DeleteCluster(&eks.DeleteClusterInput{
		Name:aws.String(clusterp.ShortName),
	})
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
		color.Blue(fmt.Sprintf("Started cluster deletion"))

		do, err := ec.DescribeCluster(&eks.DescribeClusterInput{Name: aws.String(clusterp.ShortName)})
		if err != nil {
			return nil
		}
		fmt.Printf("%v",do)

		time.Sleep(b.Duration())
		if b.Attempt() >= 20 {
			return errors.New("max retry attempts hit")
		}
	}

}
