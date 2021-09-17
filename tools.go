package main

//go:generate go run generate.go

import (
	_ "github.com/aws/aws-sdk-go-v2/config"
	_ "github.com/aws/aws-sdk-go-v2/service/ec2"
)
