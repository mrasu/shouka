package generators

import (
	"fmt"
	"strings"

	"github.com/mrasu/shouka/configs"
)

const (
	DefaultEcrRepositoryName             = "repository"
	DefaultEcrTagName                    = "production"
	DefaultEcrMigrationTagName           = "migration"
	DefaultEcsClusterName                = "cluster"
	DefaultEcsTaskName                   = "task"
	DefaultEcsMigrationTaskName          = "task-migration"
	DefaultEcsServiceName                = "service"
	DefaultCodeDeployApplicationName     = "application"
	DefaultCodeDeployDeploymentGroupName = "deployment-group"
)

type Data struct {
	Resources resourcesData
	ghActions ghActionsData
	appCode   appCodeData
	docs      docsData
}

type resourcesData struct {
	SkPrefix         string
	Region           string
	GithubOwner      string
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
	MigrationTag          string
	DefaultRepositoryName string
}

type ecsData struct {
	DefaultClusterName             string
	DefaultTaskFamilyName          string
	DefaultMigrationTaskFamilyName string
	DefaultServiceName             string
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
	AwsRegion               string
	AwsRoleArn              string
	AwsEcrRegistry          string
	AwsEcrRepository        string
	AwsEcrTag               string
	AwsEcrMigrationTag      string
	AwsEcsClusterName       string
	AwsEcsServiceName       string
	AwsEcsMigrationTaskName string
	AwsApplicationName      string
	AwsDeploymentGroupName  string
}

type appCodeData struct {
	SkPrefix                 string
	AwsTaskDefinitionExample string
}

type docsData struct {
	tldr Combination
}

func NewData(config *configs.Config) *Data {
	return &Data{
		Resources: newResourceData(config),
		ghActions: newGhActionsData(config),
		appCode:   newAppCodeData(config),
		docs:      newDocsData(),
	}
}

func newResourceData(config *configs.Config) resourcesData {
	return resourcesData{
		SkPrefix:         config.SkPrefix,
		Region:           config.Resources.Region,
		GithubOwner:      extractGithubOwner(config.GithubRepository),
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
			MigrationTag:          getEcrMigrationTag(config),
			DefaultRepositoryName: DefaultEcrRepositoryName,
		},
		Ecs: ecsData{
			DefaultClusterName:             DefaultEcsClusterName,
			DefaultTaskFamilyName:          DefaultEcsTaskName,
			DefaultMigrationTaskFamilyName: DefaultEcsMigrationTaskName,
			DefaultServiceName:             DefaultEcsServiceName,
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
		AwsRegion:               config.Resources.Region,
		AwsEcrTag:               getEcrTag(config),
		AwsEcrMigrationTag:      getEcrMigrationTag(config),
		AwsEcsClusterName:       addSkPrefix(config, DefaultEcsClusterName),
		AwsEcsServiceName:       addSkPrefix(config, DefaultEcsServiceName),
		AwsEcsMigrationTaskName: addSkPrefix(config, DefaultEcsMigrationTaskName),
		AwsApplicationName:      addSkPrefix(config, DefaultCodeDeployApplicationName),
		AwsDeploymentGroupName:  addSkPrefix(config, DefaultCodeDeployDeploymentGroupName),
	}

	if config.Resources.Iam.GithubActionsArn == "" && config.Resources.AwsAccountId != "" {
		data.AwsRoleArn = fmt.Sprintf("arn:aws:iam::%s:role/%s", config.Resources.AwsAccountId, configs.DefaultGithubActionRoleName)
	}

	if config.Resources.Ecr.RepositoryUrl == "" {
		if config.Resources.AwsAccountId != "" {
			data.AwsEcrRegistry = fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", config.Resources.AwsAccountId, config.Resources.Region)
		}
		data.AwsEcrRepository = addSkPrefix(config, DefaultEcrRepositoryName)
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
		SkPrefix: config.SkPrefix,
		AwsTaskDefinitionExample: fmt.Sprintf(
			"arn:aws:ecs:%s:%s:task-definition/%s-%s:1",
			config.Resources.Region,
			config.Resources.AwsAccountId,
			config.SkPrefix,
			DefaultEcsTaskName,
		),
	}
}

func newDocsData() docsData {
	return docsData{
		tldr: FargateGithubActions,
	}
}

func addSkPrefix(config *configs.Config, name string) string {
	return config.SkPrefix + "-" + name
}

func extractGithubOwner(repoName string) string {
	return strings.SplitN(repoName, "/", 2)[0]
}

func getEcrTag(config *configs.Config) string {
	if config.Resources.Ecr.Tag == "" {
		return DefaultEcrTagName
	} else {
		return config.Resources.Ecr.Tag
	}
}

func getEcrMigrationTag(config *configs.Config) string {
	if config.Resources.Ecr.MigrationTag == "" {
		return DefaultEcrMigrationTagName
	} else {
		return config.Resources.Ecr.MigrationTag
	}
}
