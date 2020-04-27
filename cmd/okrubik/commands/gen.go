package commands

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rubikorg/okrubik/cmd/okrubik/choose"
	"github.com/rubikorg/okrubik/pkg/entity"
)

// Gen is code generation method for rubik
// it can generate routers and routes
// and entities
func Gen(args []string) error {
	if len(args) == 0 {
		return errors.New("gen command requires arguments")
	}
	// select the project using the project selector
	path, err := choose.Project()
	if err != nil {
		return err
	}

	switch args[0] {
	case "router":
		if len(args) == 1 {
			return errors.New("router requires a name to initialize")
		}

		return genRouter(path, args[1])
		// use ast to write new router to import.go

	case "route":
		// check if router name is given as argument
		// use ast to write a new route
		// add it inside init()
		// create a new controller inside controller.go
		break
	case "entity":
		// check if name of entity given
		// loop until user enters text "done"
		break
	}
	return nil
}

func genRouter(path, name string) error {
	routerPath := filepath.Join(path, "routers", name)
	if f, _ := os.Stat(routerPath); f != nil {
		return errors.New("router with name `" + name + "` already exists")
	}
	// create new folder inside routers
	os.MkdirAll(routerPath, 0755)
	// fetch route and controller
	var files map[string]string
	en := entity.GenRouterEntity{
		RouterName: name,
	}
	en.PointTo = "/boilerplate/gen.router"
	en.Infer = &files
	_, err := rubcl.Get(en)

	if err != nil {
		return err
	}

	// create route and controller file from template
	for k, v := range files {
		fileName := strings.ReplaceAll(k, "tpl", "go")
		err = ioutil.WriteFile(filepath.Join(routerPath, fileName), []byte(v), 0655)
		if err != nil {
			return err
		}
	}

	return nil
}
