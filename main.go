package main

import (
	"embed"
	"log"

	"github.com/mrasu/shouka/cli"

	"github.com/mrasu/shouka/configs"

	"github.com/mrasu/shouka/generators"
)

//go:embed templates/*
var embedFs embed.FS

func main() {
	cnf, err := cli.AskConfig()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	/*
		cnf := &configs.Config{
			// Directory: fmt.Sprintf("tmp/tmp_terraform_%s", time.Now().Format("20060102150405")),
			Directory: "shouka_gen",
			Resources: configs.ResourceConfig{
				Region: "ap-northeast-1",
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
				},
				SecurityGroup: configs.SecurityGroupConfig{
					PublicId: "",
				},
				Subnet: configs.SubnetConfig{
					Subnet1Id: "",
					Subnet2Id: "",
				},
				Vpc: configs.VpcConfig{
					Id: "",
				},
			},
		}
	*/

	if err := generate(cnf); err != nil {
		log.Fatalf("%+v", err)
	}
}

func generate(cnf *configs.Config) error {
	// if _, err := os.Stat(dir); err == nil {
	// 	panic("already exists")
	// }
	//
	// if err := os.MkdirAll(dir, 0755); err != nil {
	// 	panic(err)
	// }
	g := generators.NewGenerator(&embedFs, cnf)
	return g.Generate()
}
