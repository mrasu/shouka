package configs_test

import (
	"testing"

	"github.com/mrasu/shouka/configs"
	"github.com/stretchr/testify/assert"
)

func TestSubnetConfig_RequiresTemplate(t *testing.T) {
	for _, td := range []struct {
		title            string
		subnet1          string
		subnet2          string
		requiresTemplate bool
	}{
		{
			title:            "No Subnet",
			requiresTemplate: true,
		},
		{
			title:            "One Subnet exists",
			subnet1:          "subnet-xxxxxxx",
			requiresTemplate: true,
		},
		{
			title:            "Two Subnets exist",
			subnet1:          "subnet-xxxxxxx",
			subnet2:          "subnet-xxxxxxx",
			requiresTemplate: false,
		},
	} {
		t.Run(td.title, func(t *testing.T) {
			sc := configs.SubnetConfig{Subnet1Id: td.subnet1, Subnet2Id: td.subnet2}
			assert.Equal(t, sc.RequiresTemplate(), td.requiresTemplate)
		})
	}
}
