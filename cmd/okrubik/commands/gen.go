package commands

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rubikorg/okrubik/cmd/okrubik/choose"
	r "github.com/rubikorg/rubik"
	"github.com/rubikorg/rubik/pkg"
)

type RouterTemplate struct {
	RouterName string
}

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

		routerPath := filepath.Join(path, "routers", args[1])
		if f, _ := os.Stat(routerPath); f != nil {
			return errors.New("router with name `" + args[1] + "` already exists")
		}
		// create new folder inside routers
		os.MkdirAll(routerPath, 0655)
		// create route file from template
		// TODO: check and download this template if not there
		routeTplPath := filepath.Join(pkg.MakeAndGetCacheDirPath(), "templates", "route.tpl")
		routeByte := r.Render(r.Type.Text, RouterTemplate{args[1]}, routeTplPath)
		if routeByte.Error != nil {
			return routeByte.Error
		}

		fileContent := routeByte.Data.([]byte)
		routePath := filepath.Join(routerPath, "route.go")
		err := ioutil.WriteFile(routePath, fileContent, 0644)
		if err != nil {
			return err
		}
		// create controller file from template
		// use ast to write new router to import.go
		break
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
