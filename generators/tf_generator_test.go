package generators_test

import (
	"testing"

	"github.com/mrasu/shouka/tests"

	"github.com/mrasu/shouka/generators"

	"github.com/mrasu/shouka/configs"
	"github.com/mrasu/shouka/generators/templates"
	"github.com/stretchr/testify/assert"
)

func TestTfCodeGenerator_Generate(t *testing.T) {
	for _, td := range []struct {
		title   string
		fn      func(*configs.Config)
		noFiles []string
	}{
		{
			title:   "Nothing exists",
			fn:      func(cnf *configs.Config) {},
			noFiles: []string{},
		},
		{
			title:   "CloudWatch exists",
			fn:      func(cnf *configs.Config) { cnf.Resources.CloudWatch.GroupName = "existing-cloud-watch" },
			noFiles: []string{"terraforms/cloud_watch.tf"},
		},
		{
			title:   "Ecr exists",
			fn:      func(cnf *configs.Config) { cnf.Resources.Ecr.RepositoryUrl = "existing-ecr" },
			noFiles: []string{"terraforms/ecr.tf"},
		},
		{
			title: "IAM exists",
			fn: func(cnf *configs.Config) {
				cnf.Resources.Iam.EcsTaskExecutionArn = "existing-task"
				cnf.Resources.Iam.CodedeployArn = "existing-codedeploy"
				cnf.Resources.Iam.GithubActionsArn = "existing-gh-actions"
				cnf.Resources.Iam.GithubActionsOpenidProviderArn = "existing-gh-actions-provider"
			},
			noFiles: []string{"terraforms/iam.tf"},
		},
		{
			title:   "SecurityGroup exists",
			fn:      func(cnf *configs.Config) { cnf.Resources.SecurityGroup.PublicId = "existing-security-group" },
			noFiles: []string{"terraforms/security_group.tf"},
		},
		{
			title: "Subnet exists",
			fn: func(cnf *configs.Config) {
				cnf.Resources.Subnet.Subnet1Id = "existing-subnet1"
				cnf.Resources.Subnet.Subnet2Id = "existing-subnet2"
			},
			noFiles: []string{"terraforms/subnet.tf"},
		},
		{
			title:   "VPC exists",
			fn:      func(cnf *configs.Config) { cnf.Resources.Vpc.Id = "existing-vpc" },
			noFiles: []string{"terraforms/vpc.tf"},
		},
	} {
		t.Run(td.title, func(t *testing.T) {
			cnf := prepareConfig()
			td.fn(cnf)

			tg, data, df := prepareTf(t, cnf)

			assert.NoError(t, tg.Generate(data))

			expectedFiles := remove(listTemplateFiles(t, "terraforms"), td.noFiles)

			for _, f := range expectedFiles {
				assert.Equal(t, "test-sk-prefix: "+f, df.GetWritten("/tmp/dummy/"+f))
			}
			assert.ElementsMatch(t, joinDir("/tmp/dummy", expectedFiles), df.WrittenFiles())
		})
	}
}

func prepareTf(t *testing.T, config *configs.Config) (*generators.TfGenerator, *generators.Data, *tests.DummyFileReadWriter) {
	t.Helper()

	files := listTemplateFiles(t, "terraforms")
	contents := map[string]string{}
	for _, f := range files {
		contents["templates/"+f+".gotmpl"] = "{{.SkPrefix}}: " + f
	}

	df := tests.NewDummyFileReadWriter(contents)
	file := templates.NewFileSystemWithReadWriter(df)

	tg := generators.NewTfGenerator(file, config)

	data := generators.NewData(config)

	return tg, data, df
}

func prepareConfig() *configs.Config {
	return &configs.Config{
		Directory: "/tmp/dummy",
		SkPrefix:  "test-sk-prefix",
	}
}

func remove(files []string, remFiles []string) (res []string) {
	for _, f := range files {
		match := false
		for _, rf := range remFiles {
			if rf == f {
				match = true
				break
			}
		}
		if !match {
			res = append(res, f)
		}
	}

	return
}
