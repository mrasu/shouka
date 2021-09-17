package configs

type SecurityGroupConfig struct {
	PublicId string
}

func (sgc *SecurityGroupConfig) RequiresTemplate() bool {
	return sgc.PublicId == ""
}
