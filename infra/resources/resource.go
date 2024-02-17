package resource

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	pipeline "github.com/aws/aws-cdk-go/awscdk/v2/awscodepipeline"
	actions "github.com/aws/aws-cdk-go/awscdk/v2/awscodepipelineactions"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	lb "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	kms "github.com/aws/aws-cdk-go/awscdk/v2/awskms"
	logs "github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
)

type ResourceService struct {
	S cdk.Stack
}

type IResourceService interface {
	// alb.go
	NewAlb(name string, vpc ec2.Vpc) lb.ApplicationLoadBalancer
	NewTargetGroup(e NewTargetGroupProps) lb.ApplicationTargetGroup

	// cloudwatch.go
	NewLogGroup(name string, key kms.IKey) logs.LogGroup
	GetLogGroupFromName(name string) logs.ILogGroup

	// codepipeline.go
	NewSourceAction(e NewSourceActionProps) SourceActionReturnValue
	NewBuildAction(e NewBuildActionProps) BuildActionReturnValue
	NewCodePipeline(e NewCodePipelineProps) pipeline.Pipeline

	// container.go
	NewCluster(e NewClusterProps) ecs.Cluster
	NewTaskDefinition(taskName string) ecs.FargateTaskDefinition
	AddContainer(e AddContainerProps) ecs.ContainerDefinition
	NewService(e NewServiceProps) ecs.FargateService
	NewServiceConnection(e NewServiceConnectionProps)

	// ecr.go
	NewEcrRepository(repositoryName string) ecr.Repository

	// iam.go
	NewAssumeRole(name string, actions []string, resources []string) iam.Role

	// kms.go
	NewKey(name string, principal string) kms.Key
	GetKeyFromName(name string) kms.IKey

	//s3.go
	NewBucket(name string) s3.Bucket
	GetBucketFromName(name string) s3.IBucket

	// vpc.go
	NewVpc(vpcName string, cidr string) ec2.Vpc
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
	Subnets        []ec2.ISubnet
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

type NewSourceActionProps struct {
	ActionName    string
	Repository    string
	Owner         string
	Branch        string
	ConnectionArn string
}

type SourceActionReturnValue struct {
	Action   actions.CodeStarConnectionsSourceAction
	Artifact pipeline.Artifact
}

type NewBuildActionProps struct {
	ActionName           string
	EcrRepositoryName    string
	ContainerName        string
	GithubRepositoryName string
	Owner                string
	Branch               string

	BuildRole      iam.IRole
	SourceArtifact pipeline.Artifact
}

type BuildActionReturnValue struct {
	Action   actions.CodeBuildAction
	Artifact pipeline.Artifact
}

type NewCodePipelineProps struct {
	Name   string
	Bucket s3.IBucket
	Stages []struct {
		Name   string
		Action pipeline.IAction
	}
}
