package resource

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	lb "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewAlb(name string, vpc ec2.Vpc) lb.ApplicationLoadBalancer {
	return lb.NewApplicationLoadBalancer(r.S, jsii.String(name), &lb.ApplicationLoadBalancerProps{
		Vpc:              vpc,
		InternetFacing:   jsii.Bool(true),
		LoadBalancerName: jsii.String(name),
		VpcSubnets:       &ec2.SubnetSelection{Subnets: vpc.PublicSubnets()},
	})
}

func (r *ResourceService) NewTargetGroup(e NewTargetGroupProps) lb.ApplicationTargetGroup {
	return lb.NewApplicationTargetGroup(r.S, jsii.String(e.Name), &lb.ApplicationTargetGroupProps{
		TargetGroupName: jsii.String(e.Name),
		TargetType:      lb.TargetType_IP,
		Port:            jsii.Number(80),
		Vpc:             e.Vpc,
		Targets:         &[]lb.IApplicationLoadBalancerTarget{e.Service},
		HealthCheck: &lb.HealthCheck{
			Path:     jsii.String("/hc"),
			Port:     jsii.String(fmt.Sprintf("%g", e.Port)),
			Interval: awscdk.Duration_Seconds(jsii.Number(300)),
		},
	})
}

func (r *ResourceService) AddListener(e AddListenerProps) lb.ApplicationListener {
	return e.ALB.AddListener(jsii.String(e.Id), &lb.BaseApplicationListenerProps{
		Protocol:            lb.ApplicationProtocol_HTTP,
		Port:                jsii.Number(e.Port),
		DefaultTargetGroups: &[]lb.IApplicationTargetGroup{e.TargetGroup},
	})
}
