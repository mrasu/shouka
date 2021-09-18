package configs

type CloudWatchConfig struct {
	GroupName string `json:"group_name"`
}

func (cwc *CloudWatchConfig) RequiresTemplate() bool {
	return cwc.GroupName == ""
}
