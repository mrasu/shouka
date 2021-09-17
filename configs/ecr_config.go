package configs

type EcrConfig struct {
	RepositoryUrl string
	Tag           string
}

func (ec *EcrConfig) RequiresTemplate() bool {
	return ec.RepositoryUrl == "" || ec.Tag == ""
}
