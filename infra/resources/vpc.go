package resource

import (
	"fmt"

	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewVpc(vpcName string, cidr string) ec2.Vpc {
	return ec2.NewVpc(r.S, jsii.String(vpcName), &ec2.VpcProps{
		VpcName:     jsii.String(vpcName),
		IpAddresses: ec2.IpAddresses_Cidr(jsii.String(cidr)),
		SubnetConfiguration: &[]*ec2.SubnetConfiguration{
			{Name: jsii.String(fmt.Sprintf("%s-public-", vpcName)), CidrMask: jsii.Number(24), SubnetType: ec2.SubnetType_PUBLIC},
			{Name: jsii.String(fmt.Sprintf("%s-private-", vpcName)), CidrMask: jsii.Number(24), SubnetType: ec2.SubnetType_PRIVATE_WITH_EGRESS},
		},
	})
}
