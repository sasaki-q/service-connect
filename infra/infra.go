package main

import (
	"fmt"
	resource "infra/resources"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscodepipeline"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfraStackProps struct {
	awscdk.StackProps
}

type Props struct {
	ConnectionArn    string
	GithubOwner      string
	GithubRepository string
	Project          string
}

const (
	VpcName string = "service-connect"
	VpcCidr string = "192.168.0.0/16"

	KeyName       string = "service-connect-log-group-key"
	LogBucketName string = "service-connect-log-bucket-2024-10-01"
	LogGroupName  string = "service-connect-log-group"

	ClusterName string = "cluster"
	Namespace   string = "local"

	RepositoryName string = "repository"

	ALBName         string = "alb"
	TargetGroupName string = "target-group"
	ListenerName    string = "listener"

	BlueTargetGroupName  string = "blue-target-group"
	BlueListener         string = "blue-listener"
	GreenTargetGroupName string = "green-target-group"
	GreenListener        string = "green-listener"

	Branch         string = "main"
	PipelineBucket string = "codepipeline-artifact-bucket-2024-02-17"
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
	key := i.NewKey(KeyName, "logs.amazonaws.com")

	// LogGroup
	logGroup := i.NewLogGroup(LogGroupName, key)

	// Bucket
	logBucket := i.NewBucket(LogBucketName)
	pipelineBucket := i.NewBucket(PipelineBucket)

	// ECR
	repository := i.NewEcrRepository(RepositoryName)

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
			"CONTAINER_HOST": jsii.String(fmt.Sprintf("%s.%s", ServerServiceName, Namespace)),
			"CONTAINER_PORT": jsii.String(fmt.Sprintf("%g", ServerPort)),
		},
		Image:    awsecs.ContainerImage_FromEcrRepository(repository, jsii.String("882367fb2ca2760abc041a3d58d9d60dc45818db")),
		LogGroup: logGroup,
		Task:     clientTaskDefinitiopn,
	})

	clientService := i.NewService(resource.NewServiceProps{
		ServiceName:    ClientServiceName,
		Port:           ClientPort,
		DesiredCount:   1,
		MaxCount:       jsii.Number(5),
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
			"CONTAINER_HOST": jsii.String(fmt.Sprintf("%s.%s", ClientServiceName, Namespace)),
			"CONTAINER_PORT": jsii.String(fmt.Sprintf("%g", ClientPort)),
		},
		Image:    awsecs.ContainerImage_FromEcrRepository(repository, jsii.String("882367fb2ca2760abc041a3d58d9d60dc45818db")),
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
	alb := i.NewAlb(ALBName, vpc)
	targetGroup := i.NewTargetGroup(resource.NewTargetGroupProps{
		Name:    TargetGroupName,
		Port:    ClientPort,
		Service: clientService,
		Vpc:     vpc,
	})
	i.AddListener(resource.AddListenerProps{
		Id:          ListenerName,
		Port:        80,
		ALB:         alb,
		TargetGroup: targetGroup,
	})

	/*
		bluetg := i.NewTargetGroup(resource.NewTargetGroupProps{
			Name:    BlueTargetGroupName,
			Port:    ClientPort,
			Service: clientService,
			Vpc:     vpc,
		})

		blueListener := i.AddListener(resource.AddListenerProps{
			Id:          BlueListener,
			Port:        80,
			ALB:         alb,
			TargetGroup: bluetg,
		})

		greentg := i.NewTargetGroup(resource.NewTargetGroupProps{
			Name:    GreenTargetGroupName,
			Port:    ClientPort,
			Service: clientService,
			Vpc:     vpc,
		})

		greenListener := i.AddListener(resource.AddListenerProps{
			Id:          GreenListener,
			Port:        ClientPort,
			ALB:         alb,
			TargetGroup: greentg,
		})
	*/

	// Code Pipeline
	buildRole := i.NewAssumeRole("buildRole", "codebuild.amazonaws.com", []string{"ecr:*", "ecs:*"}, []string{"*"})

	sourceAction := i.NewSourceAction(resource.NewSourceActionProps{
		ActionName:    "SourceAction",
		Repository:    e.GithubRepository,
		Owner:         e.GithubOwner,
		Branch:        Branch,
		ConnectionArn: e.ConnectionArn,
	})

	buildAction := i.NewBuildAction(resource.NewBuildActionProps{
		ActionName:           "BuildAction",
		Path:                 "app/cicd/build.yml",
		ContainerName:        ClientContainerName,
		EcrRepositoryName:    RepositoryName,
		TaskDefinitionArn:    *clientTaskDefinitiopn.TaskDefinitionArn(),
		GithubRepositoryName: e.GithubRepository,
		Owner:                e.GithubOwner,
		Branch:               Branch,
		BuildRole:            buildRole,
		SourceArtifact:       sourceAction.Artifact,
	})

	deployRole := i.NewAssumeRole("deployRole", "codedeploy.amazonaws.com", []string{"ecs:*", "s3:*", "iam:PassRole"}, []string{"*"})
	deployAction := i.NewRollingDeployAction(resource.NewRollingDeployActionProps{
		ActionName:    "DeployAction",
		BuildArtifact: buildAction.Artifact,
		Service:       clientService,
	})

	/*
		deployAction := i.NewBlueGreenDeployAction(resource.NewBlueGreenDeployActionProps{
			ActionName:       "DeployAction",
			Path:             "app/cicd/deploy.yml",
			ALB:              alb,
			BlueTargetGroup:  bluetg,
			BlueListener:     blueListener,
			GreenTargetGroup: greentg,
			GreenListener:    greenListener,
			Service:          clientService,
			SourceArtifact:   sourceAction.Artifact,
			BuildArtifact:    buildAction.Artifact,
		})
	*/

	pipeline := i.NewCodePipeline(resource.NewCodePipelineProps{
		Name:   fmt.Sprintf("%sCodePipeline", e.Project),
		Bucket: pipelineBucket,
		Stages: []struct {
			Name   string
			Action awscodepipeline.IAction
		}{
			{Name: "SourceStage", Action: sourceAction.Action},
			{Name: "BuildStage", Action: buildAction.Action},
			{Name: "DeployStage", Action: deployAction},
		},
	})
	deployRole.GrantAssumeRole(pipeline.Role())

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
		Props{
			ConnectionArn:    fmt.Sprintf("%s", carn),
			GithubOwner:      fmt.Sprintf("%s", gho),
			GithubRepository: fmt.Sprintf("%s", ghr),
			Project:          fmt.Sprintf("%s", project),
		},
	)

	app.Synth(nil)
}

func iToP(e interface{}) *string {
	return jsii.String(fmt.Sprintf("%s", e))
}

func myenv() *awscdk.Environment { return nil }
