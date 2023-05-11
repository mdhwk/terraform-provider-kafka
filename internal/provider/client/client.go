package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl"
	kaws "github.com/twmb/franz-go/pkg/sasl/aws"
)

type (
	Config struct {
		BootstrapServers []string
		IAM              *IAM
	}
	IAM struct {
		RoleArn     string
		SessionName string
	}
)

func NewClient(cfg *Config) (*kadm.Client, error) {
	var opts = []kgo.Opt{
		kgo.SeedBrokers(cfg.BootstrapServers...),
		//kgo.WithLogger(kgo.BasicLogger(os.Stdout, kgo.LogLevelWarn, nil)),
	}

	switch {
	case cfg.IAM != nil:
		m, err := doAwsAuth(cfg.IAM)
		if err != nil {
			return nil, err
		}
		opts = append(opts, kgo.SASL(m))
	}

	opts = append(opts, kgo.Dialer((&tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}}).DialContext))

	c, err := kgo.NewClient(opts...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return kadm.NewClient(c), nil
}

func doAwsAuth(iam *IAM) (sasl.Mechanism, error) {
	s, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	})
	if err != nil {
		return nil, err
	}

	var c kaws.Auth

	if iam.RoleArn != "" {
		c, err = doAwsRoleAssumeAuth(s, iam.RoleArn, iam.SessionName)
	} else {
		c, err = doAwsDefaultAuth(s)
	}
	if err != nil {
		return nil, err
	}

	return c.AsManagedStreamingIAMMechanism(), nil
}

func doAwsDefaultAuth(s *session.Session) (kaws.Auth, error) {
	val, err := s.Config.Credentials.GetWithContext(context.TODO())
	if err != nil {
		return kaws.Auth{}, err
	}

	a := kaws.Auth{
		AccessKey:    val.AccessKeyID,
		SecretKey:    val.SecretAccessKey,
		SessionToken: val.SessionToken,
	}

	return a, nil
}

func doAwsRoleAssumeAuth(s *session.Session, roleArn, sessionName string) (kaws.Auth, error) {
	sc := sts.New(s)
	if sessionName == "" {
		sessionName = "terraform-provider-kafka"
	}

	res, err := sc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(sessionName),
	})
	if err != nil {
		return kaws.Auth{}, err
	}

	a := kaws.Auth{
		AccessKey:    *res.Credentials.AccessKeyId,
		SecretKey:    *res.Credentials.SecretAccessKey,
		SessionToken: *res.Credentials.SessionToken,
	}

	return a, nil
}
