package resource

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewCluster(e NewClusterProps) ecs.Cluster {
	return ecs.NewCluster(r.S, jsii.String(e.ClusterName),
		&ecs.ClusterProps{
			Vpc:         e.Vpc,
			ClusterName: jsii.String(e.ClusterName),
			ExecuteCommandConfiguration: &ecs.ExecuteCommandConfiguration{
				LogConfiguration: &ecs.ExecuteCommandLogConfiguration{
					CloudWatchLogGroup:          e.LogGroup,
					CloudWatchEncryptionEnabled: jsii.Bool(true),
					S3Bucket:                    e.LogBucket,
					S3EncryptionEnabled:         jsii.Bool(true),
					S3KeyPrefix:                 jsii.String("service-connect"),
				},
				Logging: ecs.ExecuteCommandLogging_OVERRIDE,
			},
			DefaultCloudMapNamespace: &ecs.CloudMapNamespaceOptions{
				Name:                 jsii.String(e.NameSpace),
				UseForServiceConnect: jsii.Bool(true),
				Vpc:                  e.Vpc,
			},
		},
	)
}

func (r *ResourceService) NewTaskDefinition(taskName string) ecs.FargateTaskDefinition {
	return ecs.NewFargateTaskDefinition(r.S, jsii.String(taskName), &ecs.FargateTaskDefinitionProps{
		Cpu:             jsii.Number(256),
		MemoryLimitMiB:  jsii.Number(512),
		Family:          jsii.String(taskName),
		RuntimePlatform: &ecs.RuntimePlatform{CpuArchitecture: ecs.CpuArchitecture_X86_64()},
	})
}

func (r *ResourceService) AddContainer(e AddContainerProps) ecs.ContainerDefinition {
	return e.Task.AddContainer(jsii.String(e.ContainerName),
		&ecs.ContainerDefinitionOptions{
			ContainerName:  jsii.String(e.ContainerName),
			Cpu:            jsii.Number(256),
			MemoryLimitMiB: jsii.Number(512),
			Image:          e.Image,
			PortMappings: &[]*ecs.PortMapping{
				{
					AppProtocol:   ecs.AppProtocol_Http(),
					Name:          jsii.String(e.PortMappingName), // == ServiceConnectConfiguration.Services.PortMappingName
					Protocol:      ecs.Protocol_TCP,
					HostPort:      jsii.Number(e.Port),
					ContainerPort: jsii.Number(e.Port),
				},
			},
			Logging:     ecs.LogDriver_AwsLogs(&ecs.AwsLogDriverProps{StreamPrefix: jsii.String(e.ContainerName), LogGroup: e.LogGroup}),
			Environment: &e.Env,
		},
	)
}

func (r *ResourceService) NewService(e NewServiceProps) ecs.FargateService {
	return ecs.NewFargateService(r.S, jsii.String(e.ServiceName), &ecs.FargateServiceProps{
		Cluster:              e.Cluster,
		DesiredCount:         jsii.Number(1),
		EnableExecuteCommand: jsii.Bool(true),
		ServiceName:          jsii.String(e.ServiceName),
		TaskDefinition:       e.TaskDefinition,
		VpcSubnets:           &awsec2.SubnetSelection{Subnets: &e.Subnets},
		ServiceConnectConfiguration: &ecs.ServiceConnectProps{
			Services: &[]*ecs.ServiceConnectService{
				{
					DiscoveryName:   jsii.String(e.ServiceName),
					Port:            jsii.Number(e.Port),
					PortMappingName: jsii.String(e.ServiceName),
				},
			},
			LogDriver: ecs.LogDriver_AwsLogs(&ecs.AwsLogDriverProps{
				StreamPrefix: jsii.String(fmt.Sprintf("service-connect/%s", e.ServiceName)),
				LogGroup:     e.LogGroup,
			}),
		},
	})
}

func (r *ResourceService) NewServiceConnection(e NewServiceConnectionProps) {
	e.ToConnection.AllowFrom(e.FromConnection, awsec2.Port_Tcp(jsii.Number(e.ToPort)), nil)
}
