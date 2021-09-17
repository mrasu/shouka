package configs

type ResourceConfig struct {
	Region           string
	AvailabilityZone AvailabilityZoneConfig
	CloudWatch       CloudWatchConfig
	Ecr              EcrConfig
	Iam              IamConfig
	SecurityGroup    SecurityGroupConfig
	Subnet           SubnetConfig
	Vpc              VpcConfig
}
