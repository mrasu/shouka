package cli

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

type resourceAnswer struct {
	VpcId                  string
	CloudWatchGroupName    string
	EcrRepositoryUrl       string
	EcrTag                 string
	IamEcsTaskExecutionArn string
	IamCodedeployArn       string
	SecurityGroupPublicId  string
	Subnet1Id              string
	Subnet2Id              string
}

func AskResources() (*configs.ResourceConfig, error) {
	ri := &resourceInput{}

	cnf := &configs.ResourceConfig{}
	region, err := ri.askRegion()
	if err != nil {
		return nil, err
	}
	cnf.Region = region

	resources, err := ri.askExistingResources()
	if err != nil {
		return nil, err
	}

	ans := resourceAnswer{}
	var inputQs []*survey.Question
	var setters []func(resourceAnswer)

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
		q, setter := ri.createVpcIdQuestion(cnf)

		inputQs = append(inputQs, q)
		setters = append(setters, setter)
	}

	if _, ok := resources[cloudWatch]; ok {
		q, setter := ri.createCloudWatchQuestion(cnf)

		inputQs = append(inputQs, q)
		setters = append(setters, setter)
	}

	if _, ok := resources[ecr]; ok {
		qs, setter := ri.createEcrQuestion(cnf)

		for _, q := range qs {
			inputQs = append(inputQs, q)
		}
		setters = append(setters, setter)
	}

	if _, ok := resources[iam]; ok {
		qs, setter := ri.createIamQuestion(cnf)

		for _, q := range qs {
			inputQs = append(inputQs, q)
		}
		setters = append(setters, setter)
	}

	if _, ok := resources[securityGroup]; ok {
		q, setter := ri.createSecurityGroupQuestion(cnf)

		inputQs = append(inputQs, q)
		setters = append(setters, setter)
	}

	if _, ok := resources[subnet]; ok {
		qs, setter := ri.createSubnetQuestion(cnf)

		for _, q := range qs {
			inputQs = append(inputQs, q)
		}
		setters = append(setters, setter)
	}

	if err := survey.Ask(inputQs, &ans, survey.WithValidator(survey.Required)); err != nil {
		return nil, errors.Wrap(err, "failed to get answers for ResourceConfig")
	}
	for _, s := range setters {
		s(ans)
	}

	return cnf, nil
}

func (ri *resourceInput) askRegion() (string, error) {
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
		return "", errors.Wrap(err, "failed to get resourceAnswer for region")
	}

	return region, nil
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

func (ri *resourceInput) createVpcIdQuestion(cnf *configs.ResourceConfig) (*survey.Question, func(resourceAnswer)) {
	return &survey.Question{
			Name: "VpcId",
			Prompt: &survey.Input{
				Message: "Tell VPC ID you use",
			},
		}, func(ans resourceAnswer) {
			cnf.Vpc = configs.VpcConfig{
				Id: ans.VpcId,
			}
		}
}

func (ri *resourceInput) createCloudWatchQuestion(cnf *configs.ResourceConfig) (*survey.Question, func(resourceAnswer)) {
	return &survey.Question{
			Name: "CloudWatchGroupName",
			Prompt: &survey.Input{
				Message: "Tell Name of CloudWatch's log group you use",
			},
		}, func(ans resourceAnswer) {
			cnf.CloudWatch = configs.CloudWatchConfig{
				GroupName: ans.CloudWatchGroupName,
			}
		}
}

func (ri *resourceInput) createEcrQuestion(cnf *configs.ResourceConfig) ([]*survey.Question, func(resourceAnswer)) {
	return []*survey.Question{
			{
				Name: "EcrRepositoryUrl",
				Prompt: &survey.Input{
					Message: "Tell URL of ECR you use",
				},
			}, {
				Name: "EcrTag",
				Prompt: &survey.Input{
					Message: "Tell Tag of Docker image you use",
				},
			},
		}, func(ans resourceAnswer) {
			cnf.Ecr = configs.EcrConfig{
				RepositoryUrl: ans.EcrRepositoryUrl,
				Tag:           ans.EcrTag,
			}
		}
}

func (ri *resourceInput) createIamQuestion(cnf *configs.ResourceConfig) ([]*survey.Question, func(resourceAnswer)) {
	return []*survey.Question{
			{
				Name: "IamEcsTaskExecutionArn",
				Prompt: &survey.Input{
					Message: "Tell Arn of IAM role you use to execute task of ECS",
				},
			}, {
				Name: "IamCodedeployArn",
				Prompt: &survey.Input{
					Message: "Tell Arn of IAM role you use for Codedeploy",
				},
			},
		}, func(ans resourceAnswer) {
			cnf.Iam = configs.IamConfig{
				EcsTaskExecutionArn: ans.IamEcsTaskExecutionArn,
				CodedeployArn:       ans.IamCodedeployArn,
			}
		}
}

func (ri *resourceInput) createSecurityGroupQuestion(cnf *configs.ResourceConfig) (*survey.Question, func(resourceAnswer)) {
	return &survey.Question{
			Name: "SecurityGroupPublicId",
			Prompt: &survey.Input{
				Message: "Tell Security group ID you use accessible from public",
			},
		}, func(ans resourceAnswer) {
			cnf.SecurityGroup = configs.SecurityGroupConfig{
				PublicId: ans.SecurityGroupPublicId,
			}
		}
}

func (ri *resourceInput) createSubnetQuestion(cnf *configs.ResourceConfig) ([]*survey.Question, func(resourceAnswer)) {
	return []*survey.Question{
			{
				Name: "Subnet1Id",
				Prompt: &survey.Input{
					Message: "Tell Subnet you use for an availability zone",
				},
			}, {
				Name: "Subnet2Id",
				Prompt: &survey.Input{
					Message: "Tell one more Subnet you use for a different availability zone",
				},
			},
		}, func(ans resourceAnswer) {
			cnf.Subnet = configs.SubnetConfig{
				Subnet1Id: ans.Subnet1Id,
				Subnet2Id: ans.Subnet2Id,
			}
		}
}
