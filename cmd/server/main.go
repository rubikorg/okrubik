package main

import (
	cfg "github.com/rubikorg/okrubik/cmd/server/config"
	"github.com/rubikorg/okrubik/cmd/server/routers"
	"github.com/rubikorg/rubik"
)

func main() {
	var config cfg.ProjectConfig
	err := rubik.Load(&config)

	if err != nil {
		panic(err)
	}

	routers.Import()
	panic(rubik.Run())

}
