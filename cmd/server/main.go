package main

import (
	"github.com/okcherry/cherry"
	"github.com/okcherry/okcherry/cmd/server/config/routers"
)

type ProjectConfig struct {
	Port string `toml:"port"`
}

func main() {
	var config ProjectConfig
	err := cherry.App.Load(&config)

	if err != nil {
		panic(err)
	}

	routers.Import()
	panic(cherry.App.Listen())

}
