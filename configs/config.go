package configs

type Config struct {
	Directory        string         `json:"directory"`
	GithubRepository string         `json:"github_repository"`
	Resources        ResourceConfig `json:"resources"`
}
