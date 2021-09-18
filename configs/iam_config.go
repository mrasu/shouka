package configs

const (
	DefaultGithubActionRoleName = "GitHubActionRole"
)

type IamConfig struct {
	EcsTaskExecutionArn string `json:"ecs_task_execution_arn"`
	CodedeployArn       string `json:"codedeploy_arn"`

	GithubActionsArn               string `json:"github_actions_arn"`
	GithubActionsOpenidProviderArn string `json:"github_actions_openid_provider_arn"`
}

func (ic *IamConfig) RequiresTemplate() bool {
	return ic.EcsTaskExecutionArn == "" || ic.CodedeployArn == ""
}
