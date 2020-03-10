package main

import (
	cfg "github.com/oksketch/oksketch/cmd/server/config"
	"github.com/oksketch/oksketch/cmd/server/routers"
	"github.com/oksketch/sketch"
)

func main() {
	var config cfg.ProjectConfig
	err := sketch.Load(&config)

	if err != nil {
		panic(err)
	}

	routers.Import()
	panic(sketch.Listen())

}
