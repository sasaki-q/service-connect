package resource

import (
	"fmt"

	build "github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	pipeline "github.com/aws/aws-cdk-go/awscdk/v2/awscodepipeline"
	actions "github.com/aws/aws-cdk-go/awscdk/v2/awscodepipelineactions"
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
			BuildSpec: build.BuildSpec_FromSourceFilename(jsii.String("app/cicd/build.yml")),
			Environment: &build.BuildEnvironment{
				BuildImage: build.LinuxBuildImage_AMAZON_LINUX_2(),
				Privileged: jsii.Bool(true),
			},
			EnvironmentVariables: &map[string]*build.BuildEnvironmentVariable{
				"AWS_DEFAULT_REGION": {Value: r.S.Region()},
				"REPOSITORY_NAME":    {Value: e.EcrRepositoryName},
				"CONTAINER_NAME":     {Value: e.ContainerName},
			},
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
