package resource

import (
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewAssumeRole(name string, actions []string, resources []string) iam.Role {
	role := iam.NewRole(r.S, jsii.String(name),
		&iam.RoleProps{
			AssumedBy: iam.NewServicePrincipal(jsii.String("codebuild.amazonaws.com"), nil),
			RoleName:  jsii.String(name),
		},
	)

	role.AddToPolicy(iam.NewPolicyStatement(&iam.PolicyStatementProps{
		Actions:   vToP(actions),
		Effect:    iam.Effect_ALLOW,
		Resources: vToP(resources),
	}))

	return role
}
