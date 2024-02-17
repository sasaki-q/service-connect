package resource

import (
	kms "github.com/aws/aws-cdk-go/awscdk/v2/awskms"
	logs "github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewLogGroup(name string, key kms.IKey) logs.LogGroup {
	return logs.NewLogGroup(r.S, jsii.String(name), &logs.LogGroupProps{
		EncryptionKey: key,
		LogGroupName:  jsii.String(name),
	})
}

func (r *ResourceService) GetLogGroupFromName(name string) logs.ILogGroup {
	return logs.LogGroup_FromLogGroupName(r.S, jsii.String(name), jsii.String(name))
}
