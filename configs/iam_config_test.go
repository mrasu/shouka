package configs_test

import (
	"testing"

	"github.com/mrasu/shouka/configs"
	"github.com/stretchr/testify/assert"
)

func TestIamConfig_RequiresTemplate(t *testing.T) {
	for _, td := range []struct {
		title            string
		fn               func(*configs.IamConfig)
		requiresTemplate bool
	}{
		{
			title:            "All exists",
			fn:               func(cnf *configs.IamConfig) {},
			requiresTemplate: false,
		},
		{
			title:            "No Task Role",
			fn:               func(cnf *configs.IamConfig) { cnf.EcsTaskExecutionArn = "" },
			requiresTemplate: true,
		},
		{
			title:            "No CodeDeploy Role",
			fn:               func(cnf *configs.IamConfig) { cnf.CodedeployArn = "" },
			requiresTemplate: true,
		},
		{
			title:            "No Gh Role",
			fn:               func(cnf *configs.IamConfig) { cnf.GithubActionsArn = "" },
			requiresTemplate: true,
		},
		{
			title:            "No Gh provider",
			fn:               func(cnf *configs.IamConfig) { cnf.GithubActionsOpenidProviderArn = "" },
			requiresTemplate: true,
		},
	} {
		t.Run(td.title, func(t *testing.T) {
			ic := &configs.IamConfig{
				EcsTaskExecutionArn:            "arn:aws:iam::xxxxx:role/task-role",
				CodedeployArn:                  "arn:aws:iam::xxxxx:role/codedeploy-role",
				GithubActionsArn:               "arn:aws:iam::xxxxx:role/task-role",
				GithubActionsOpenidProviderArn: "arn:aws:iam::xxxxx:oidc-provider/accounts.google.com",
			}
			td.fn(ic)
			assert.Equal(t, ic.RequiresTemplate(), td.requiresTemplate)
		})
	}
}
