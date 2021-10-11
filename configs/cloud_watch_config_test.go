package configs_test

import (
	"testing"

	"github.com/mrasu/shouka/configs"
	"github.com/stretchr/testify/assert"
)

func TestCloudWatchConfig_RequiresTemplate(t *testing.T) {
	for _, td := range []struct {
		title            string
		groupName        string
		requiresTemplate bool
	}{
		{
			title:            "No GroupName",
			groupName:        "",
			requiresTemplate: true,
		},
		{
			title:            "GroupName exists",
			groupName:        "hello",
			requiresTemplate: false,
		},
	} {
		t.Run(td.title, func(t *testing.T) {
			cwc := configs.CloudWatchConfig{GroupName: td.groupName}
			assert.Equal(t, cwc.RequiresTemplate(), td.requiresTemplate)
		})
	}
}
