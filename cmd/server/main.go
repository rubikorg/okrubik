package main

import (
	"fmt"

	cfg "github.com/rubikorg/okrubik/cmd/server/config"
	"github.com/rubikorg/okrubik/cmd/server/routers"
	"github.com/rubikorg/okrubik/pkg/services"
	"github.com/rubikorg/rubik"

	_ "github.com/rubikorg/blocks/apigen/flutter"
	_ "github.com/rubikorg/blocks/apigen/ts"
	_ "github.com/rubikorg/blocks/guard/jwt"
	_ "github.com/rubikorg/blocks/logger"
	_ "github.com/rubikorg/blocks/swagger"
)

func main() {
	var config cfg.ProjectConfig
	err := rubik.Load(&config)
	if err != nil {
		panic(err)
	}
	routers.Import()
	err = rubik.Run(services.CoreService)
	if err != nil {
		fmt.Println(err.Error())
	}
}
