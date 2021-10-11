package configs_test

import (
	"testing"

	"github.com/mrasu/shouka/configs"
	"github.com/stretchr/testify/assert"
)

func TestVpcConfig_RequiresTemplate(t *testing.T) {
	for _, td := range []struct {
		title            string
		id               string
		requiresTemplate bool
	}{
		{
			title:            "No SecurityGroup for public internet",
			id:               "",
			requiresTemplate: true,
		},
		{
			title:            "SecurityGroup for public internet exists",
			id:               "vpc-xxxxxxx",
			requiresTemplate: false,
		},
	} {
		t.Run(td.title, func(t *testing.T) {
			vc := configs.VpcConfig{Id: td.id}
			assert.Equal(t, vc.RequiresTemplate(), td.requiresTemplate)
		})
	}
}
