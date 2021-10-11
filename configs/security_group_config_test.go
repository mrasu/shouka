package configs_test

import (
	"testing"

	"github.com/mrasu/shouka/configs"
	"github.com/stretchr/testify/assert"
)

func TestSecurityGroupConfig_RequiresTemplate(t *testing.T) {
	for _, td := range []struct {
		title            string
		publicId         string
		requiresTemplate bool
	}{
		{
			title:            "No SecurityGroup for public internet",
			publicId:         "",
			requiresTemplate: true,
		},
		{
			title:            "SecurityGroup for public internet exists",
			publicId:         "sg-xxxxxxx",
			requiresTemplate: false,
		},
	} {
		t.Run(td.title, func(t *testing.T) {
			sgc := configs.SecurityGroupConfig{PublicId: td.publicId}
			assert.Equal(t, sgc.RequiresTemplate(), td.requiresTemplate)
		})
	}
}
