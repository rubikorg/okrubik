package main

import (
	"fmt"

	_ "github.com/rubikorg/blocks/logger"
	cfg "github.com/rubikorg/okrubik/cmd/server/config"
	"github.com/rubikorg/okrubik/cmd/server/routers"
	"github.com/rubikorg/rubik"

	_ "github.com/rubikorg/blocks/apigen/flutter"
	_ "github.com/rubikorg/blocks/apigen/ts"
	_ "github.com/rubikorg/blocks/guard/jwt"
	_ "github.com/rubikorg/blocks/swagger"
)

func main() {
	var config cfg.ProjectConfig
	err := rubik.Load(&config)
	if err != nil {
		panic(err)
	}
	routers.Import()
	err = rubik.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
