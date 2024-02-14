package resource

import (
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewEcrRepository(repositoryName string) ecr.Repository {
	return ecr.NewRepository(r.S, jsii.String(repositoryName), &ecr.RepositoryProps{
		ImageScanOnPush:    jsii.Bool(true),
		ImageTagMutability: ecr.TagMutability_IMMUTABLE,
		RepositoryName:     jsii.String(repositoryName),
	})
}
