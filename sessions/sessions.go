// Package sessions provides functions that return AWS sessions to use in the AWS SDK.
package sessions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	maxRetriesOnRecoverableFailures = 8 // Default provided by SDK is 3 which means requests are retried up to only 2 seconds.
	credsTimeout                    = 10 * time.Second
	clientTimeout                   = 30 * time.Second
)

// NewSession returns a session configured against the "default" AWS profile.
func NewSession() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            *newConfig(),
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}
	if aws.StringValue(sess.Config.Region) == "" {
		return nil, fmt.Errorf("region is missing")
	}
	return sess, nil
}

func newConfig() *aws.Config {
	c := &http.Client{
		Timeout: clientTimeout,
	}
	return aws.NewConfig().
		WithHTTPClient(c).
		WithCredentialsChainVerboseErrors(true).
		WithMaxRetries(maxRetriesOnRecoverableFailures)
}
