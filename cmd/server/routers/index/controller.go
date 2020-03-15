package index

import (
	config2 "github.com/rubikorg/okrubik/cmd/server/config"
	"github.com/rubikorg/rubik"
)

func indexCtl(en interface{}) (interface{}, error) {
	config := rubik.GetConfig().(config2.ProjectConfig)
	return config.Bench, nil
}
