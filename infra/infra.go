package main

import (
	"fmt"
	resource "infra/resources"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	lb "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfraStackProps struct {
	awscdk.StackProps
}

type Props struct {
	GithubAccessToken string
	HostedZoneId      string
}

const (
	VpcName string = "service-connect"
	VpcCidr string = "192.168.0.0/16"

	KeyName       string = "key"
	LogBucketName string = "service-connect-log-bucket"
	LogGroupName  string = "service-connect-log-group"

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

func NewInfraStack(scope constructs.Construct, id string, props *InfraStackProps, e Props) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	var i resource.IResourceService = &resource.ResourceService{S: stack}

	// VPC
	vpc := i.NewVpc(VpcName, VpcCidr)

	// KMS
	key := i.NewKey(KeyName)

	// LogGroup
	logGroup := i.NewLogGroup(LogGroupName, key)

	// Bucket
	logBucket := i.NewBucket(LogBucketName)

	// ECR
	clientRepository := i.NewEcrRepository(ClientRepositoryName)
	serverRepository := i.NewEcrRepository(ServerRepositoryName)

	// ECS
	cluster := i.NewCluster(resource.NewClusterProps{
		ClusterName: ClusterName,
		NameSpace:   Namespace,
		LogBucket:   logBucket,
		LogGroup:    logGroup,
		Vpc:         vpc,
	})

	// Client Container
	clientTaskDefinitiopn := i.NewTaskDefinition(ClientTaskName)

	i.AddContainer(resource.AddContainerProps{
		ContainerName:   ClientContainerName,
		Port:            ClientPort,
		PortMappingName: ClientServiceName,
		Env: map[string]*string{
			"PORT":           jsii.String(fmt.Sprintf("%g", ClientPort)),
			"CONTAINER_NAME": jsii.String(ClientContainerName),
			"CONTAINER_HOST": jsii.String(ServerContainerName),
			"CONTAINER_PORT": jsii.String(fmt.Sprintf("%g", ServerPort)),
		},
		Image:    awsecs.ContainerImage_FromEcrRepository(clientRepository, jsii.String("v0.1")),
		LogGroup: logGroup,
		Task:     clientTaskDefinitiopn,
	})

	clientService := i.NewService(resource.NewServiceProps{
		ServiceName:    ClientServiceName,
		Port:           ClientPort,
		Cluster:        cluster,
		LogGroup:       logGroup,
		Subnets:        *vpc.PrivateSubnets(),
		TaskDefinition: clientTaskDefinitiopn,
	})

	// Server Container
	serverTaskDefinitiopn := i.NewTaskDefinition(ServerTaskName)

	i.AddContainer(resource.AddContainerProps{
		ContainerName:   ServerContainerName,
		Port:            ServerPort,
		PortMappingName: ServerServiceName,
		Env: map[string]*string{
			"PORT":           jsii.String(fmt.Sprintf("%g", ServerPort)),
			"CONTAINER_NAME": jsii.String(ServerContainerName),
			"CONTAINER_HOST": jsii.String(ClientContainerName),
			"CONTAINER_PORT": jsii.String(fmt.Sprintf("%g", ClientPort)),
		},
		Image:    awsecs.ContainerImage_FromEcrRepository(serverRepository, jsii.String("v0.1")),
		LogGroup: logGroup,
		Task:     serverTaskDefinitiopn,
	})

	serverService := i.NewService(resource.NewServiceProps{
		ServiceName:    ServerServiceName,
		Port:           ServerPort,
		Cluster:        cluster,
		LogGroup:       logGroup,
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

const (
	BootstrapBucketName string = "BBN"
	ConnectionArn       string = "CARN"
	Env                 string = "ENV"
	GithubAccessToken   string = "GHAT"
	GithubOwner         string = "GHO"
	GithubRepository    string = "GHR"
	HostedZoneId        string = "HGI"
	Id                  string = "ID"
	Project             string = "PROJECT"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	var (
		bbn     = app.Node().TryGetContext(jsii.String(BootstrapBucketName))
		carn    = app.Node().TryGetContext(jsii.String(ConnectionArn))
		env     = app.Node().TryGetContext(jsii.String(Env))
		ght     = app.Node().TryGetContext(jsii.String(GithubAccessToken))
		gho     = app.Node().TryGetContext(jsii.String(GithubOwner))
		ghr     = app.Node().TryGetContext(jsii.String(GithubRepository))
		hgi     = app.Node().TryGetContext(jsii.String(HostedZoneId))
		id      = app.Node().TryGetContext(jsii.String(Id))
		project = app.Node().TryGetContext(jsii.String(Project))
	)

	if bbn == nil || carn == nil || env == nil || ght == nil || gho == nil || ghr == nil || hgi == nil || project == nil || id == nil {
		panic("please pass context")
	}

	awscdk.Tags_Of(app).Add(jsii.String("Project"), iToP(project), nil)
	NewInfraStack(app, fmt.Sprintf("%sStack", project),
		&InfraStackProps{
			awscdk.StackProps{
				Env: myenv(),
				Synthesizer: awscdk.NewDefaultStackSynthesizer(
					&awscdk.DefaultStackSynthesizerProps{FileAssetsBucketName: iToP(bbn)},
				),
			},
		},
		Props{},
	)

	app.Synth(nil)
}

func iToP(e interface{}) *string {
	return jsii.String(fmt.Sprintf("%s", e))
}

func myenv() *awscdk.Environment { return nil }
