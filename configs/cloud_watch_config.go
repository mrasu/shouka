package configs

type CloudWatchConfig struct {
	GroupName string
}

func (cwc *CloudWatchConfig) RequiresTemplate() bool {
	return cwc.GroupName == ""
}
