package resource

import (
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	kms "github.com/aws/aws-cdk-go/awscdk/v2/awskms"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewKey(name string, principal string) kms.Key {
	return kms.NewKey(r.S, jsii.String(name), &kms.KeyProps{
		Alias: jsii.String(name),
		Policy: iam.NewPolicyDocument(&iam.PolicyDocumentProps{
			Statements: &[]iam.PolicyStatement{
				iam.NewPolicyStatement(
					&iam.PolicyStatementProps{
						Actions: &[]*string{jsii.String("kms:*")},
						Effect:  iam.Effect_ALLOW,
						Principals: &[]iam.IPrincipal{
							iam.NewAccountRootPrincipal(),
							iam.NewServicePrincipal(jsii.String(principal), nil),
						},
						Resources: &[]*string{jsii.String("*")},
					},
				),
			},
		}),
	})
}

func (r *ResourceService) GetKeyFromName(name string) kms.IKey {
	return kms.Alias_FromAliasName(r.S, jsii.String(name), jsii.String(name))
}
