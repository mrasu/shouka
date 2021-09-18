package configs

type ResourceConfig struct {
	AwsAccountId     string                 `json:"aws_account_id"`
	Region           string                 `json:"region"`
	AvailabilityZone AvailabilityZoneConfig `json:"availability_zone"`
	CloudWatch       CloudWatchConfig       `json:"cloud_watch"`
	Ecr              EcrConfig              `json:"ecr"`
	Iam              IamConfig              `json:"iam"`
	SecurityGroup    SecurityGroupConfig    `json:"security_group"`
	Subnet           SubnetConfig           `json:"subnet"`
	Vpc              VpcConfig              `json:"vpc"`
}
