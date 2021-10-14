package main

import (
	"embed"
	"encoding/json"

	"github.com/mrasu/shouka/cmd"
	"github.com/mrasu/shouka/configs"
	"github.com/mrasu/shouka/injections"
	"github.com/mrasu/shouka/libs/log"
)

// Embed subdirectories explicitly to include hidden files like .gitignore.gotmpl
//go:embed templates/*
//go:embed templates/terraforms/*
var embedFs embed.FS

func main() {
	injections.EmbedFs = &embedFs

	// runDummy()

	//*
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error happens. error: %+v", err)
	}
	//*/
}

//nolint:deadcode,unused
func runDummy() {
	cnf := &configs.Config{
		// Directory: fmt.Sprintf("tmp/tmp_terraform_%s", time.Now().Format("20060102150405")),
		Directory:        "tmp/dummy_shouka_gen",
		GithubRepository: "mrasu/shouka_gen",
		Resources: configs.ResourceConfig{
			AwsAccountId: "889435949642",
			Region:       "ap-northeast-1",
			AvailabilityZone: configs.AvailabilityZoneConfig{
				Zone1: "ap-northeast-1a",
				Zone2: "ap-northeast-1c",
			},
			CloudWatch: configs.CloudWatchConfig{
				GroupName: "",
			},
			Ecr: configs.EcrConfig{
				RepositoryUrl: "",
				Tag:           "",
			},
			Iam: configs.IamConfig{
				EcsTaskExecutionArn: "",
				CodedeployArn:       "",

				GithubActionsArn:               "b",
				GithubActionsOpenidProviderArn: "c",
			},
			SecurityGroup: configs.SecurityGroupConfig{
				PublicId: "a",
			},
			Subnet: configs.SubnetConfig{
				Subnet1Id: "1",
				Subnet2Id: "2",
			},
			Vpc: configs.VpcConfig{
				Id: "",
			},
		},
	}

	log.Println(cnf)
	data, err := json.Marshal(cnf)
	if err != nil {
		panic(err)
	}
	log.Printf("%s\n", string(data))

	v := configs.Config{}
	err = json.Unmarshal(data, &v)
	if err != nil {
		panic(err)
	}

	log.Println(&v)
}
