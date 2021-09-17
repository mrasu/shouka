package configs

type VpcConfig struct {
	Id string
}

func (vc *VpcConfig) RequiresTemplate() bool {
	return vc.Id == ""
}
