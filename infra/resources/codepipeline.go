package resource

import (
	"fmt"

	build "github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	deploy "github.com/aws/aws-cdk-go/awscdk/v2/awscodedeploy"
	pipeline "github.com/aws/aws-cdk-go/awscdk/v2/awscodepipeline"
	actions "github.com/aws/aws-cdk-go/awscdk/v2/awscodepipelineactions"
	lb "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewSourceAction(e NewSourceActionProps) SourceActionReturnValue {
	artifact := pipeline.NewArtifact(jsii.String(e.ActionName))

	action := actions.NewCodeStarConnectionsSourceAction(&actions.CodeStarConnectionsSourceActionProps{
		ActionName:    jsii.String(e.ActionName),
		RunOrder:      jsii.Number(1),
		Repo:          jsii.String(e.Repository),
		Owner:         jsii.String(e.Owner),
		Branch:        jsii.String(e.Branch),
		TriggerOnPush: jsii.Bool(true),
		Output:        artifact,
		ConnectionArn: jsii.String(e.ConnectionArn),
	})

	return SourceActionReturnValue{
		Action:   action,
		Artifact: artifact,
	}
}

func (r *ResourceService) NewBuildAction(e NewBuildActionProps) BuildActionReturnValue {
	artifact := pipeline.NewArtifact(jsii.String(e.ActionName))

	project := build.NewProject(r.S, jsii.String(fmt.Sprintf("%sProject", e.ActionName)),
		&build.ProjectProps{
			BuildSpec: build.BuildSpec_FromSourceFilename(jsii.String(e.Path)),
			Environment: &build.BuildEnvironment{
				BuildImage: build.LinuxBuildImage_AMAZON_LINUX_2(),
				Privileged: jsii.Bool(true),
			},
			EnvironmentVariables: &map[string]*build.BuildEnvironmentVariable{
				"AWS_DEFAULT_REGION":   {Value: r.S.Region()},
				"REPOSITORY_NAME":      {Value: e.EcrRepositoryName},
				"CONTAINER_NAME":       {Value: e.ContainerName},
				"TASK_DEFINITION_ARN":  {Value: e.TaskDefinitionArn},
				"BUILD_IMAGE_ARN":      {Value: fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/go:1.21.0-bullseye", *r.S.Account(), *r.S.Region())},
				"PRODUCTION_IMAGE_ARN": {Value: fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/debian:bullseye", *r.S.Account(), *r.S.Region())},
			},
			ProjectName: jsii.String(fmt.Sprintf("%sProject", e.ActionName)),
			Source: build.Source_GitHub(&build.GitHubSourceProps{
				Identifier:  jsii.String(fmt.Sprintf("ID_%s", e.ActionName)),
				Repo:        jsii.String(e.GithubRepositoryName),
				Owner:       jsii.String(e.Owner),
				BranchOrRef: jsii.String(e.Branch),
			}),
			Role: e.BuildRole,
		},
	)

	action := actions.NewCodeBuildAction(
		&actions.CodeBuildActionProps{
			ActionName: jsii.String(e.ActionName),
			Type:       actions.CodeBuildActionType_BUILD,
			Project:    project,
			RunOrder:   jsii.Number(1),
			Input:      e.SourceArtifact,
			Outputs:    &[]pipeline.Artifact{artifact},
		},
	)

	return BuildActionReturnValue{
		Action:   action,
		Artifact: artifact,
	}
}

func (r *ResourceService) NewRollingDeployAction(e NewRollingDeployActionProps) actions.EcsDeployAction {
	return actions.NewEcsDeployAction(&actions.EcsDeployActionProps{
		ActionName: jsii.String(e.ActionName),
		Service:    e.Service,
		Input:      e.BuildArtifact,
	})
}

// https://repost.aws/questions/QUWtGsiusrRPaAhOI4xRRvpw/ecs-fargate-with-ecs-connect-and-code-deploy
func (r *ResourceService) NewBlueGreenDeployAction(e NewDeployActionProps) actions.CodeDeployEcsDeployAction {
	deploymentGroup := deploy.NewEcsDeploymentGroup(r.S, jsii.String(e.ActionName),
		&deploy.EcsDeploymentGroupProps{
			BlueGreenDeploymentConfig: &deploy.EcsBlueGreenDeploymentConfig{
				BlueTargetGroup: lb.ApplicationTargetGroup_FromTargetGroupAttributes(
					r.S,
					e.BlueTargetGroup.TargetGroupName(),
					&lb.TargetGroupAttributes{
						TargetGroupArn:   e.BlueTargetGroup.TargetGroupArn(),
						LoadBalancerArns: e.ALB.LoadBalancerArn(),
					},
				),
				GreenTargetGroup: lb.ApplicationTargetGroup_FromTargetGroupAttributes(
					r.S,
					e.GreenTargetGroup.TargetGroupName(),
					&lb.TargetGroupAttributes{
						TargetGroupArn:   e.GreenTargetGroup.TargetGroupArn(),
						LoadBalancerArns: e.ALB.LoadBalancerArn(),
					},
				),
				Listener:     e.BlueListener,
				TestListener: e.GreenListener,
			},
			DeploymentConfig:    deploy.EcsDeploymentConfig_CANARY_10PERCENT_5MINUTES(),
			DeploymentGroupName: jsii.String(fmt.Sprintf("%sDeploymentGroup", e.ActionName)),
			Service:             e.Service,
		},
	)

	return actions.NewCodeDeployEcsDeployAction(
		&actions.CodeDeployEcsDeployActionProps{
			ActionName:                 jsii.String(e.ActionName),
			RunOrder:                   jsii.Number(1),
			DeploymentGroup:            deploymentGroup,
			AppSpecTemplateFile:        pipeline.NewArtifactPath(e.SourceArtifact, jsii.String(e.Path)),
			TaskDefinitionTemplateFile: pipeline.NewArtifactPath(e.BuildArtifact, jsii.String("taskdef.json")),
			ContainerImageInputs:       &[]*actions.CodeDeployEcsContainerImageInput{{Input: e.BuildArtifact}},
		},
	)
}

func (r *ResourceService) NewCodePipeline(e NewCodePipelineProps) pipeline.Pipeline {
	stages := []*pipeline.StageProps{}

	for _, v := range e.Stages {
		stages = append(
			stages,
			&pipeline.StageProps{StageName: jsii.String(v.Name), Actions: &[]pipeline.IAction{v.Action}},
		)
	}

	return pipeline.NewPipeline(r.S, jsii.String(e.Name),
		&pipeline.PipelineProps{
			ArtifactBucket: e.Bucket,
			PipelineName:   jsii.String(e.Name),
			Stages:         &stages,
		},
	)
}
