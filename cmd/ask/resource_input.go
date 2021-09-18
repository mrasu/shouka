package ask

import (
	survey "github.com/AlecAivazis/survey/v2"
	"github.com/mrasu/shouka/configs"
	"github.com/mrasu/shouka/constants"
	"github.com/pkg/errors"
)

const (
	securityGroup = "security_group"
	subnet        = "subnet"
	cloudWatch    = "cloud_watch"
	ecr           = "ecr"
	iam           = "iam"
	vpc           = "vpc"
)

type resourceInput struct{}

func AskResources() (*configs.ResourceConfig, error) {
	ri := &resourceInput{}

	cnf := &configs.ResourceConfig{}

	if err := ri.askAwsAccountId(cnf); err != nil {
		return nil, err
	}

	if err := ri.askRegion(cnf); err != nil {
		return nil, err
	}

	resources, err := ri.askExistingResources()
	if err != nil {
		return nil, err
	}

	if _, ok := resources[vpc]; !ok {
		az1, az2, err := ri.askAvailabilityZones(cnf.Region)
		if err != nil {
			return nil, err
		}

		cnf.AvailabilityZone = configs.AvailabilityZoneConfig{
			Zone1: az1,
			Zone2: az2,
		}
	} else {
		if err := ri.askVpcIdQuestion(cnf); err != nil {
			return nil, err
		}
	}

	if _, ok := resources[cloudWatch]; ok {
		if err := ri.askCloudWatch(cnf); err != nil {
			return nil, err
		}
	}

	if _, ok := resources[ecr]; ok {
		if err := ri.askEcr(cnf); err != nil {
			return nil, err
		}
	}

	if _, ok := resources[iam]; ok {
		if err := ri.askIam(cnf); err != nil {
			return nil, err
		}
	}

	if _, ok := resources[securityGroup]; ok {
		if err := ri.askSecurityGroup(cnf); err != nil {
			return nil, err
		}
	}

	if _, ok := resources[subnet]; ok {
		if err := ri.askSubnet(cnf); err != nil {
			return nil, err
		}
	}

	return cnf, nil
}

func (ri *resourceInput) askRegion(cnf *configs.ResourceConfig) error {
	var regions []string
	for _, zone := range constants.Zones {
		regions = append(regions, zone.Region)
	}

	prompt := &survey.Select{
		Message:  "Choose aws region",
		Options:  regions,
		PageSize: 20,
	}

	region := ""
	err := survey.AskOne(prompt, &region)
	if err != nil {
		return errors.Wrap(err, "failed to get resourceAnswer for region")
	}
	cnf.Region = region
	return nil
}

func (ri *resourceInput) askAwsAccountId(cnf *configs.ResourceConfig) error {
	prompt := &survey.Input{
		Message: "Tell Id for AWS' account",
		Help:    "You can get it by `aws sts get-caller-identity`",
	}

	id := ""
	if err := survey.AskOne(prompt, &id); err != nil {
		return errors.Wrap(err, "failed to get AwsAccountId")
	}

	cnf.AwsAccountId = id
	return nil
}

func (ri *resourceInput) askExistingResources() (map[string]struct{}, error) {
	prompt := &survey.MultiSelect{
		Message: "Choose existing resources you want to use",
		Options: []string{cloudWatch, ecr, iam, vpc},
	}

	var resources []string
	if err := survey.AskOne(prompt, &resources); err != nil {
		return nil, errors.Wrap(err, "failed to get resourceAnswer for resources")
	}

	for _, r := range resources {
		if r == "vpc" {
			prompt := &survey.MultiSelect{
				Message: "Because you reuse vpc, you can choose more existing resources you want to use",
				Options: []string{securityGroup, subnet},
			}

			var moreResources []string
			if err := survey.AskOne(prompt, &moreResources); err != nil {
				return nil, errors.Wrap(err, "failed to get resourceAnswer for more resources")
			}

			for _, r := range moreResources {
				resources = append(resources, r)
			}
		}
	}

	resourceSet := map[string]struct{}{}
	for _, r := range resources {
		resourceSet[r] = struct{}{}
	}

	return resourceSet, nil
}

func (ri *resourceInput) askAvailabilityZones(region string) (string, string, error) {
	var azOptions []string
	for _, zone := range constants.Zones {
		if zone.Region == region {
			azOptions = zone.AvailabilityZones
			break
		}
	}

	prompt := &survey.MultiSelect{
		Message: "Choose two availability zones you want to use for new vpc",
		Options: azOptions,
	}

	var azs []string
	err := survey.AskOne(prompt, &azs, func(options *survey.AskOptions) error {
		options.Validators = append(options.Validators, survey.MinItems(2))
		options.Validators = append(options.Validators, survey.MaxItems(2))
		return nil
	})

	if err != nil {
		return "", "", errors.Wrap(err, "failed to get resourceAnswer for availability zones")
	}

	return azs[0], azs[1], nil
}

func (ri *resourceInput) askVpcIdQuestion(cnf *configs.ResourceConfig) error {
	prompt := &survey.Input{
		Message: "Tell VPC ID you use",
	}

	res := ""
	if err := askRequiredOne(prompt, &res); err != nil {
		return errors.Wrap(err, "failed to get VPC data")
	}

	cnf.Vpc = configs.VpcConfig{
		Id: res,
	}
	return nil
}

func (ri *resourceInput) askCloudWatch(cnf *configs.ResourceConfig) error {
	prompt := &survey.Input{
		Message: "Tell Name of CloudWatch's log group you use",
	}

	res := ""
	if err := askRequiredOne(prompt, &res); err != nil {
		return errors.Wrap(err, "failed to get CloudWatch data")
	}

	cnf.CloudWatch = configs.CloudWatchConfig{
		GroupName: res,
	}
	return nil
}

func (ri *resourceInput) askEcr(cnf *configs.ResourceConfig) error {
	qs := []*survey.Question{
		{
			Name: "Url",
			Prompt: &survey.Input{
				Message: "Tell URL of ECR you use",
			},
			Validate: survey.Required,
		}, {
			Name: "Tag",
			Prompt: &survey.Input{
				Message: "Tell Tag of Docker image you use",
			},
			Validate: survey.Required,
		},
	}

	res := struct {
		Url string
		Tag string
	}{}

	if err := survey.Ask(qs, &res); err != nil {
		return errors.Wrap(err, "failed to get ECR data")
	}

	cnf.Ecr = configs.EcrConfig{
		RepositoryUrl: res.Url,
		Tag:           res.Tag,
	}

	return nil
}

func (ri *resourceInput) askIam(cnf *configs.ResourceConfig) error {
	qs := []*survey.Question{
		{
			Name: "EcsTaskExecutionArn",
			Prompt: &survey.Input{
				Message: "Tell Arn of IAM role you use to execute task of ECS",
			},
			Validate: survey.Required,
		}, {
			Name: "CodedeployArn",
			Prompt: &survey.Input{
				Message: "Tell Arn of IAM role you use for Codedeploy",
			},
			Validate: survey.Required,
		}, {
			Name: "GithubActionsArn",
			Prompt: &survey.Input{
				Message: "Tell Arn of IAM role you use to connect GitHub Actions",
				Help:    "Keep empty if you don't have",
			},
		},
	}

	res := struct {
		EcsTaskExecutionArn string
		CodedeployArn       string
		GithubActionsArn    string
	}{}

	if err := survey.Ask(qs, &res); err != nil {
		return errors.Wrap(err, "failed to get IAM data")
	}

	cnf.Iam = configs.IamConfig{
		EcsTaskExecutionArn: res.EcsTaskExecutionArn,
		CodedeployArn:       res.CodedeployArn,
		GithubActionsArn:    res.GithubActionsArn,
	}

	if res.GithubActionsArn == "" {
		prompt := &survey.Input{
			Message: "Tell Arn of IAM Identity provider you use to connect GitHub Actions",
			Help:    "Keep empty if you don't have",
		}

		res := ""
		if err := askRequiredOne(prompt, &res); err != nil {
			return errors.Wrap(err, "failed to get GithubActionsOpenidProviderArn")
		}

		cnf.Iam.GithubActionsOpenidProviderArn = res
	}

	return nil
}

func (ri *resourceInput) askSecurityGroup(cnf *configs.ResourceConfig) error {
	prompt := &survey.Input{
		Message: "Tell Security group ID you use accessible from public",
	}

	res := ""
	if err := askRequiredOne(prompt, &res); err != nil {
		return errors.Wrap(err, "failed to get SecurityGroup data")
	}

	cnf.SecurityGroup = configs.SecurityGroupConfig{
		PublicId: res,
	}
	return nil
}

func (ri *resourceInput) askSubnet(cnf *configs.ResourceConfig) error {
	qs := []*survey.Question{
		{
			Name: "Subnet1Id",
			Prompt: &survey.Input{
				Message: "Tell Subnet you use for an availability zone",
			},
			Validate: survey.Required,
		}, {
			Name: "Subnet2Id",
			Prompt: &survey.Input{
				Message: "Tell one more Subnet you use for a different availability zone",
			},
			Validate: survey.Required,
		},
	}

	res := struct {
		Subnet1Id string
		Subnet2Id string
	}{}

	if err := survey.Ask(qs, &res); err != nil {
		return errors.Wrap(err, "failed to get Subnet data")
	}

	cnf.Subnet = configs.SubnetConfig{
		Subnet1Id: res.Subnet1Id,
		Subnet2Id: res.Subnet2Id,
	}

	return nil
}

func askRequiredOne(p survey.Prompt, res interface{}) error {
	return survey.AskOne(p, res, survey.WithValidator(survey.Required))
}
