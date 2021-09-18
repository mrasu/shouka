package ask

import (
	survey "github.com/AlecAivazis/survey/v2"
	"github.com/mrasu/shouka/configs"
	"github.com/pkg/errors"
)

type configAnswer struct {
	Directory        string
	GithubRepository string
}

func AskConfig() (*configs.Config, error) {
	cnf := &configs.Config{}
	qs := []*survey.Question{
		{
			Name: "Directory",
			Prompt: &survey.Input{
				Message: "Tell directory to shouka generate codes",
				Default: "shouka_gen",
			},
		}, {
			Name: "GithubRepository",
			Prompt: &survey.Input{
				Message: "Tell repository of GitHub you use",
				Help:    "e.g. \"mrasu/shouka\"",
			},
		},
	}
	ans := configAnswer{}
	if err := survey.Ask(qs, &ans, survey.WithValidator(survey.Required)); err != nil {
		return nil, errors.Wrap(err, "failed to get answers for Config")
	}

	cnf.Directory = ans.Directory
	cnf.GithubRepository = ans.GithubRepository

	rCnf, err := AskResources()
	if err != nil {
		return nil, err
	}
	cnf.Resources = *rCnf

	return cnf, nil
}
