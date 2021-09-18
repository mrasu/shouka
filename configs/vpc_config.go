package configs

type VpcConfig struct {
	Id string `json:"id"`
}

func (vc *VpcConfig) RequiresTemplate() bool {
	return vc.Id == ""
}
