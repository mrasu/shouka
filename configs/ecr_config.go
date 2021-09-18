package configs

import "strings"

type EcrConfig struct {
	RepositoryUrl string `json:"repository_url"`
	Tag           string `json:"tag"`
}

func (ec *EcrConfig) RequiresTemplate() bool {
	return ec.RepositoryUrl == ""
}

func (ec *EcrConfig) RepositoryDomain() string {
	domain, _ := ec.splitRepositoryDomainAndName()
	return domain
}

func (ec *EcrConfig) RepositoryName() string {
	_, name := ec.splitRepositoryDomainAndName()
	return name
}

func (ec *EcrConfig) splitRepositoryDomainAndName() (string, string) {
	words := strings.Split(ec.RepositoryUrl, "/")
	if len(words) != 2 {
		return "", ""
	}
	return words[0], words[1]
}
