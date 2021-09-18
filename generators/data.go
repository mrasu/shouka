package generators

import (
	"fmt"

	"github.com/mrasu/shouka/configs"
)

const (
	DefaultEcrRepositoryName             = "sk-repository"
	DefaultEcrTagName                    = "production"
	DefaultEcsTaskName                   = "sk-task"
	DefaultCodeDeployApplicationName     = "sk-application"
	DefaultCodeDeployDeploymentGroupName = "sk-deployment-group"
)

type data struct {
	Resources resourcesData
	ghActions ghActionsData
	appCode   appCodeData
}

type resourcesData struct {
	Region           string
	GithubRepository string
	AvailabilityZone availabilityZoneData
	CloudWatch       cloudWatchData
	Ecr              ecrData
	Ecs              ecsData
	Iam              iamData
	SecurityGroup    securityGroupData
	Subnet           subnetData
	Vpc              vpcData
}

type availabilityZoneData struct {
	Zone1 string
	Zone2 string
}

type cloudWatchData struct {
	GroupName string
}

type ecrData struct {
	RepositoryUrl         string
	Tag                   string
	DefaultRepositoryName string
}

type ecsData struct {
	DefaultTaskFamilyName string
}

type iamData struct {
	EcsTaskExecutionArn string
	CodedeployArn       string

	GithubActionsArn               string
	GithubActionsOpenidProviderArn string
	GithubRepository               string
}

type securityGroupData struct {
	PublicId string
}

type subnetData struct {
	Subnet1Id string
	Subnet2Id string
}

type vpcData struct {
	Id string
}

type ghActionsData struct {
	AwsRoleArn             string
	AwsEcrRegistry         string
	AwsEcrRepository       string
	AwsEcrTag              string
	AwsApplicationName     string
	AwsDeploymentGroupName string
}

type appCodeData struct {
	AwsTaskDefinitionExample string
}

func newData(config *configs.Config) *data {
	return &data{
		Resources: newResourceData(config),
		ghActions: newGhActionsData(config),
		appCode:   newAppCodeData(config),
	}
}

func newResourceData(config *configs.Config) resourcesData {
	return resourcesData{
		Region:           config.Resources.Region,
		GithubRepository: config.GithubRepository,
		AvailabilityZone: availabilityZoneData{
			Zone1: config.Resources.AvailabilityZone.Zone1,
			Zone2: config.Resources.AvailabilityZone.Zone2,
		},
		CloudWatch: cloudWatchData{
			GroupName: config.Resources.CloudWatch.GroupName,
		},
		Ecr: ecrData{
			RepositoryUrl:         config.Resources.Ecr.RepositoryUrl,
			Tag:                   getEcrTag(config),
			DefaultRepositoryName: DefaultEcrRepositoryName,
		},
		Ecs: ecsData{
			DefaultTaskFamilyName: DefaultEcsTaskName,
		},
		Iam: iamData{
			EcsTaskExecutionArn:            config.Resources.Iam.EcsTaskExecutionArn,
			CodedeployArn:                  config.Resources.Iam.CodedeployArn,
			GithubActionsArn:               config.Resources.Iam.GithubActionsArn,
			GithubActionsOpenidProviderArn: config.Resources.Iam.GithubActionsOpenidProviderArn,
		},
		SecurityGroup: securityGroupData{
			PublicId: config.Resources.SecurityGroup.PublicId,
		},
		Subnet: subnetData{
			Subnet1Id: config.Resources.Subnet.Subnet1Id,
			Subnet2Id: config.Resources.Subnet.Subnet2Id,
		},
		Vpc: vpcData{
			Id: config.Resources.Vpc.Id,
		},
	}
}

func newGhActionsData(config *configs.Config) ghActionsData {
	data := ghActionsData{
		AwsEcrTag:              getEcrTag(config),
		AwsApplicationName:     DefaultCodeDeployApplicationName,
		AwsDeploymentGroupName: DefaultCodeDeployDeploymentGroupName,
	}

	if config.Resources.Iam.GithubActionsArn == "" && config.Resources.AwsAccountId != "" {
		data.AwsRoleArn = fmt.Sprintf("arn:aws:iam::%s:role/%s", config.Resources.AwsAccountId, configs.DefaultGithubActionRoleName)
	}

	if config.Resources.Ecr.RepositoryUrl == "" {
		if config.Resources.AwsAccountId != "" {
			data.AwsEcrRegistry = fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", config.Resources.AwsAccountId, config.Resources.Region)
		}
		data.AwsEcrRepository = DefaultEcrRepositoryName
	} else {
		if config.Resources.Ecr.RepositoryDomain() != "" {
			data.AwsEcrRegistry = config.Resources.Ecr.RepositoryDomain()
		}
		if config.Resources.Ecr.RepositoryName() != "" {
			data.AwsEcrRepository = config.Resources.Ecr.RepositoryName()
		}
	}

	return data
}

func newAppCodeData(config *configs.Config) appCodeData {
	return appCodeData{
		AwsTaskDefinitionExample: fmt.Sprintf(
			"arn:aws:ecs:%s:%s:task-definition/%s:1",
			config.Resources.Region,
			config.Resources.AwsAccountId,
			DefaultEcsTaskName,
		),
	}
}

func getEcrTag(config *configs.Config) string {
	if config.Resources.Ecr.Tag == "" {
		return DefaultEcrTagName
	} else {
		return config.Resources.Ecr.Tag
	}
}