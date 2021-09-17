package configs

type SubnetConfig struct {
	Subnet1Id string
	Subnet2Id string
}

func (sc *SubnetConfig) RequiresTemplate() bool {
	return sc.Subnet1Id == "" || sc.Subnet2Id == ""
}
