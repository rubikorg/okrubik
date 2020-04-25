package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rubikorg/okrubik/cmd/okrubik/choose"
	r "github.com/rubikorg/rubik"
	"github.com/rubikorg/rubik/pkg"

	"text/template"
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
		os.MkdirAll(routerPath, 0755)
		// create route file from template
		tplData, err := getOrDownloadTemplateData("route.tpl")
		ctlTplData, err := getOrDownloadTemplateData("controller.tpl")
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		var ctlBuf bytes.Buffer
		t, err := template.New("route").Parse(tplData)
		ctlT, err := template.New("controller").Parse(ctlTplData)
		err = t.Execute(&buf, RouterTemplate{args[1]})
		err = ctlT.Execute(&ctlBuf, RouterTemplate{args[1]})
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filepath.Join(routerPath, "route.go"), buf.Bytes(), 0655)
		// create controller file from template
		err = ioutil.WriteFile(filepath.Join(routerPath, "controller.go"), ctlBuf.Bytes(), 0655)
		if err != nil {
			return err
		}

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

func getOrDownloadTemplateData(templateName string) (string, error) {
	routeTplPath := filepath.Join(pkg.MakeAndGetCacheDirPath(), "templates")
	os.MkdirAll(routeTplPath, 0755)
	routeTplPath = filepath.Join(routeTplPath, templateName)

	var templateStr string
	if f, _ := os.Stat(routeTplPath); f == nil {
		// dowblaod the file
		routeTpleDownload := r.DownloadRequestEntity{
			TargetFilePath: routeTplPath,
		}
		routeTpleDownload.PointTo = "/boilerplate/" + templateName
		raw, err := rubcl.Download(routeTpleDownload)
		if err != nil {
			return "", err
		}
		templateStr = string(raw)
		fmt.Println("rawstring", templateStr)
	} else {
		b, err := ioutil.ReadFile(routeTplPath)
		if err != nil {
			return "", err
		}

		templateStr = string(b)
	}
	return templateStr, nil
}
