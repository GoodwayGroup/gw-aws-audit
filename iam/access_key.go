package iam

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awsIAM "github.com/aws/aws-sdk-go/service/iam"
)

type AccessKey struct {
	id          *string
	createdDate *time.Time
	status      *string
	lastUsed    *awsIAM.AccessKeyLastUsed
	userName    *string
}

func (ak *AccessKey) Deactivate() error {
	kl := k.Extend("AccessKey:Deactivate")
	client := session.GetIAMClient()

	var err error
	var result *awsIAM.UpdateAccessKeyOutput
	result, err = client.UpdateAccessKey(&awsIAM.UpdateAccessKeyInput{
		AccessKeyId: ak.id,
		Status:      aws.String("Inactive"),
		UserName:    ak.userName,
	})
	if err != nil {
		return err
	}

	kl.Log(result)
	return nil
}

func (ak *AccessKey) Activate() error {
	kl := k.Extend("AccessKey:Activate")
	client := session.GetIAMClient()

	var err error
	var result *awsIAM.UpdateAccessKeyOutput
	result, err = client.UpdateAccessKey(&awsIAM.UpdateAccessKeyInput{
		AccessKeyId: ak.id,
		Status:      aws.String("Active"),
		UserName:    ak.userName,
	})
	if err != nil {
		return err
	}

	kl.Log(result)
	return nil
}

func (ak *AccessKey) Delete() error {
	kl := k.Extend("AccessKey:Delete")
	client := session.GetIAMClient()

	var err error
	var result *awsIAM.DeleteAccessKeyOutput
	result, err = client.DeleteAccessKey(&awsIAM.DeleteAccessKeyInput{
		AccessKeyId: ak.id,
		UserName:    ak.userName,
	})
	if err != nil {
		return err
	}

	kl.Log(result)
	return nil
}
