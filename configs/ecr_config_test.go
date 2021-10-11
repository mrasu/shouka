package configs_test

import (
	"testing"

	"github.com/mrasu/shouka/configs"
	"github.com/stretchr/testify/assert"
)

func TestEcrConfig_RequiresTemplate(t *testing.T) {
	for _, td := range []struct {
		title            string
		repositoryUrl    string
		requiresTemplate bool
	}{
		{
			title:            "No RepositoryUrl",
			requiresTemplate: true,
		},
		{
			title:            "RepositoryUrl exists",
			repositoryUrl:    "hello",
			requiresTemplate: false,
		},
	} {
		t.Run(td.title, func(t *testing.T) {
			ec := configs.EcrConfig{RepositoryUrl: td.repositoryUrl}
			assert.Equal(t, ec.RequiresTemplate(), td.requiresTemplate)
		})
	}
}

func TestEcrConfig_RepositoryDomain(t *testing.T) {
	for _, td := range []struct {
		title         string
		repositoryUrl string
		domain        string
	}{
		{
			title:         "No RepositoryUrl",
			repositoryUrl: "",
			domain:        "",
		},
		{
			title:         "RepositoryUrl exists",
			repositoryUrl: "xxxx.dkr.ecr.us-east-1.amazonaws.com/sk-repository",
			domain:        "xxxx.dkr.ecr.us-east-1.amazonaws.com",
		},
		{
			title:         "RepositoryUrl contains only domain",
			repositoryUrl: "xxxx.dkr.ecr.us-east-1.amazonaws.com",
			domain:        "",
		},
	} {
		t.Run(td.title, func(t *testing.T) {
			ec := configs.EcrConfig{RepositoryUrl: td.repositoryUrl}
			assert.Equal(t, ec.RepositoryDomain(), td.domain)
		})
	}
}

func TestEcrConfig_RepositoryName(t *testing.T) {
	for _, td := range []struct {
		title         string
		repositoryUrl string
		repoName      string
	}{
		{
			title:         "No RepositoryUrl",
			repositoryUrl: "",
			repoName:      "",
		},
		{
			title:         "RepositoryUrl exists",
			repositoryUrl: "xxxx.dkr.ecr.us-east-1.amazonaws.com/sk-repository",
			repoName:      "sk-repository",
		},
		{
			title:         "RepositoryUrl contains only domain",
			repositoryUrl: "xxxx.dkr.ecr.us-east-1.amazonaws.com",
			repoName:      "",
		},
	} {
		t.Run(td.title, func(t *testing.T) {
			ec := configs.EcrConfig{RepositoryUrl: td.repositoryUrl}
			assert.Equal(t, ec.RepositoryName(), td.repoName)
		})
	}
}
