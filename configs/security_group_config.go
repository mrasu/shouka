package configs

type SecurityGroupConfig struct {
	PublicId string `json:"public_id"`
}

func (sgc *SecurityGroupConfig) RequiresTemplate() bool {
	return sgc.PublicId == ""
}
