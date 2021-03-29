package lib

import (
	"fmt"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	awsIAM "github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/s3"
	as "github.com/clok/awssession"
	"github.com/clok/kemba"
)

var (
	ksess = kemba.New("gw-aws-audit:lib:session")
)

type Session struct {
	session *awsSession.Session
}

func (s *Session) init() error {
	ksess.Println("init")
	sess, err := as.New()
	if err != nil {
		return err
	}
	s.session = sess
	return nil
}

func (s *Session) getOrInit() *awsSession.Session {
	if s.session == nil {
		ksess.Extend("getOrInit").Println("creating new session")
		if err := s.init(); err != nil {
			_ = fmt.Errorf("unable to build session: %e", err)
			return nil
		}
	}
	return s.Session()
}

func (s *Session) Session() *awsSession.Session {
	return s.session
}

func (s *Session) GetIAMClient() *awsIAM.IAM {
	return awsIAM.New(s.getOrInit())
}

func (s *Session) GetEC2Client() *ec2.EC2 {
	return ec2.New(s.getOrInit())
}

func (s *Session) GetRDSClient() *rds.RDS {
	return rds.New(s.getOrInit())
}

func (s *Session) GetS3Client() *s3.S3 {
	return s3.New(s.getOrInit())
}

func (s *Session) GetCloudWatchClient() *cloudwatch.CloudWatch {
	return cloudwatch.New(s.getOrInit())
}
