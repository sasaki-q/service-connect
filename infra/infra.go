package main

import (
	resource "infra/resources"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	lb "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	logs "github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfraStackProps struct {
	awscdk.StackProps
}

type Props struct {
	lbt *string
	lg  *string
}

const (
	VpcName string = "service-connect"
	VpcCidr string = "192.168.0.0/16"

	ClusterName string = "cluster"
	Namespace   string = "local"
)

const (
	ClientRepositoryName string  = "client_repository"
	ClientTaskName       string  = "client_task_definition"
	ClientContainerName  string  = "client_container"
	ClientPort           float64 = 8000
	ClientServiceName    string  = "client_service"
)

const (
	ServerRepositoryName string  = "server_repository"
	ServerTaskName       string  = "server_task_definition"
	ServerContainerName  string  = "server_container"
	ServerPort           float64 = 8001
	ServerServiceName    string  = "server_service"
)

func NewInfraStack(scope constructs.Construct, id string, props *InfraStackProps, arg Props) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	lbt := awss3.Bucket_FromBucketName(stack, arg.lbt, arg.lbt)
	lg := logs.LogGroup_FromLogGroupName(stack, arg.lg, arg.lg)
	var i resource.IResourceService = &resource.ResourceService{S: stack}

	// VPC
	vpc := i.NewVpc(VpcName, VpcCidr)

	// ECR
	clientRepository := i.NewEcrRepository(ClientRepositoryName)
	serverRepository := i.NewEcrRepository(ServerRepositoryName)

	// ECS
	cluster := i.NewCluster(resource.NewClusterProps{
		ClusterName: ClusterName,
		NameSpace:   Namespace,
		LogBucket:   lbt,
		LogGroup:    lg,
		Vpc:         vpc,
	})

	// Client Container
	clientTaskDefinitiopn := i.NewTaskDefinition(ClientTaskName)

	i.AddContainer(resource.AddContainerProps{
		ContainerName:   ClientContainerName,
		Port:            ClientPort,
		PortMappingName: ClientServiceName,
		Env:             map[string]*string{},
		Image:           awsecs.ContainerImage_FromEcrRepository(clientRepository, jsii.String("v0.1")),
		LogGroup:        lg,
		Task:            clientTaskDefinitiopn,
	})

	clientService := i.NewService(resource.NewServiceProps{
		ServiceName:    ClientServiceName,
		Port:           ClientPort,
		Cluster:        cluster,
		LogGroup:       lg,
		Subnets:        *vpc.PrivateSubnets(),
		TaskDefinition: clientTaskDefinitiopn,
	})

	// Server Container
	serverTaskDefinitiopn := i.NewTaskDefinition(ServerTaskName)

	i.AddContainer(resource.AddContainerProps{
		ContainerName:   ServerContainerName,
		Port:            ServerPort,
		PortMappingName: ServerServiceName,
		Env:             map[string]*string{},
		Image:           awsecs.ContainerImage_FromEcrRepository(serverRepository, jsii.String("v0.1")),
		LogGroup:        lg,
		Task:            serverTaskDefinitiopn,
	})

	serverService := i.NewService(resource.NewServiceProps{
		ServiceName:    ServerServiceName,
		Port:           ServerPort,
		Cluster:        cluster,
		LogGroup:       lg,
		Subnets:        *vpc.PrivateSubnets(),
		TaskDefinition: serverTaskDefinitiopn,
	})

	// Service Connect
	i.NewServiceConnection(resource.NewServiceConnectionProps{
		ToConnection:   serverService.Connections(),
		ToPort:         ServerPort,
		FromConnection: clientService.Connections(),
	})

	// Load Balancer
	alb := i.NewAlb("alb", vpc)
	tg := i.NewTargetGroup(resource.NewTargetGroupProps{
		Name:    "tg",
		Port:    ClientPort,
		Service: clientService,
		Vpc:     vpc,
	})

	alb.AddListener(jsii.String("listener"), &lb.BaseApplicationListenerProps{
		Protocol:            lb.ApplicationProtocol_HTTP,
		DefaultTargetGroups: &[]lb.IApplicationTargetGroup{tg},
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	envs, ok := app.Node().TryGetContext(jsii.String("dev")).(map[string]interface{})
	if !ok {
		panic("ERROR")
	}

	var (
		bt  = jsii.String(envs["bt"].(string))
		lbt = jsii.String(envs["lbt"].(string))
		lg  = jsii.String(envs["lg"].(string))
	)
	if bt == nil || lbt == nil || lg == nil {
		panic("ERROR")
	}

	awscdk.Tags_Of(app).Add(jsii.String("Project"), jsii.String("Service-Connect"), nil)
	NewInfraStack(app, "InfraStack",
		&InfraStackProps{
			awscdk.StackProps{
				Env: env(),
				Synthesizer: awscdk.NewDefaultStackSynthesizer(
					&awscdk.DefaultStackSynthesizerProps{FileAssetsBucketName: bt},
				),
			},
		},
		Props{lbt: lbt, lg: lg},
	)

	app.Synth(nil)
}

func env() *awscdk.Environment { return nil }
