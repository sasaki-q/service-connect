package resource

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	kms "github.com/aws/aws-cdk-go/awscdk/v2/awskms"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewKey(name string) kms.Key {
	return kms.NewKey(r.S, jsii.String(name), &kms.KeyProps{
		Alias: jsii.String(name),
		Policy: awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
			Statements: &[]awsiam.PolicyStatement{
				awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
					Actions:   &[]*string{jsii.String("*")},
					Effect:    awsiam.Effect_ALLOW,
					Resources: &[]*string{jsii.String("*")},
				}),
			},
		}),
	})
}

func (r *ResourceService) GetKeyFromName(name string) kms.IKey {
	return kms.Alias_FromAliasName(r.S, jsii.String(name), jsii.String(name))
}
