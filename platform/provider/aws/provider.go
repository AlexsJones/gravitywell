package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

type AWSProvider struct {
	AWSClient *session.Session
	Region    string
}
