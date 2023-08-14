package aws

import (
	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/configservice/configserviceiface"
	"github.com/gruntwork-io/cloud-nuke/config"
	"github.com/gruntwork-io/go-commons/errors"
)

type ConfigServiceRecorders struct {
	Client        configserviceiface.ConfigServiceAPI
	Region        string
	RecorderNames []string
}

func (csr ConfigServiceRecorders) ResourceName() string {
	return "config-recorders"
}

func (csr ConfigServiceRecorders) ResourceIdentifiers() []string {
	return csr.RecorderNames
}

func (csr ConfigServiceRecorders) MaxBatchSize() int {
	return 50
}

func (csr ConfigServiceRecorders) GetAndSetIdentifiers(configObj config.Config) ([]string, error) {
	identifiers, err := csr.getAll(configObj)
	if err != nil {
		return nil, err
	}

	csr.RecorderNames = awsgo.StringValueSlice(identifiers)
	return csr.RecorderNames, nil
}

func (csr ConfigServiceRecorders) Nuke(configServiceRecorderNames []string) error {
	if err := csr.nukeAll(configServiceRecorderNames); err != nil {
		return errors.WithStackTrace(err)
	}
	return nil
}
