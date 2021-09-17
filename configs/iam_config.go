package configs

type IamConfig struct {
	EcsTaskExecutionArn string
	CodedeployArn       string
}

func (ic *IamConfig) RequiresTemplate() bool {
	return ic.EcsTaskExecutionArn == "" || ic.CodedeployArn == ""
}
