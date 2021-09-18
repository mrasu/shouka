package configs

type SubnetConfig struct {
	Subnet1Id string `json:"subnet_1_id"`
	Subnet2Id string `json:"subnet_2_id"`
}

func (sc *SubnetConfig) RequiresTemplate() bool {
	return sc.Subnet1Id == "" || sc.Subnet2Id == ""
}
