package resource

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	lb "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	logs "github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
)

type ResourceService struct {
	S awscdk.Stack
}

type IResourceService interface {
	NewVpc(vpcName string, cidr string) ec2.Vpc

	NewEcrRepository(repositoryName string) ecr.Repository

	NewCluster(e NewClusterProps) ecs.Cluster
	NewTaskDefinition(taskName string) ecs.FargateTaskDefinition
	AddContainer(e AddContainerProps) ecs.ContainerDefinition
	NewService(e NewServiceProps) ecs.FargateService
	NewServiceConnection(e NewServiceConnectionProps)

	NewAlb(name string, vpc ec2.Vpc) lb.ApplicationLoadBalancer
	NewTargetGroup(e NewTargetGroupProps) lb.ApplicationTargetGroup
}

type NewClusterProps struct {
	ClusterName string
	NameSpace   string

	LogBucket s3.IBucket
	LogGroup  logs.ILogGroup
	Vpc       ec2.IVpc
}

type AddContainerProps struct {
	ContainerName   string
	Env             map[string]*string
	Port            float64
	PortMappingName string

	Image    ecs.ContainerImage
	LogGroup logs.ILogGroup
	Task     ecs.TaskDefinition
}

type NewServiceProps struct {
	ServiceName string
	Port        float64

	Cluster        ecs.ICluster
	LogGroup       logs.ILogGroup
	Subnets        []awsec2.ISubnet
	TaskDefinition ecs.TaskDefinition
}

type NewServiceConnectionProps struct {
	ToConnection   ec2.Connections
	ToPort         float64
	FromConnection ec2.Connections
}

type NewTargetGroupProps struct {
	Name    string
	Port    float64
	Service ecs.FargateService
	Vpc     ec2.Vpc
}
