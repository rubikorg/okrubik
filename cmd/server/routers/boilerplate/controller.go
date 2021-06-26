package boilerplate

import (
	"path/filepath"
	"strings"

	"github.com/rubikorg/okrubik/pkg/entity"
	r "github.com/rubikorg/rubik"
)

var genRouterFiles = []string{
	"route.tpl",
	"controller.tpl",
	"controller_test.tpl",
}

func createCtl(req *r.Request) {
	createFiles := []string{
		"app-config.go.tpl",
		"app-dep.go.tpl",
		"routers-index-controller.go.tpl",
		"main.go.tpl",
		"#rubik.toml.tpl",
		"templates-index.html.tpl",
		"config-default.toml.tpl",
		"routers-index-route.go.tpl",
		"routers-import.go.tpl",
		"static-main.css.tpl",
		"#README.md.tpl",
	}
	var compiled = make(map[string]string)

	createEn := req.Entity.(*entity.CreateBoilerplateEntity)
	// if it is a new workspace create this file for user
	if createEn.IsNew {
		createFiles = append(createFiles, "#pkg-services-list.go.tpl")
	}

	for _, file := range createFiles {
		var cleanFileName = file
		if strings.HasPrefix(file, "#") {
			cleanFileName = strings.ReplaceAll(file, "#", "")
		}
		joinedPath := filepath.Join("create", cleanFileName)
		result := r.RenderContent(r.Type.Text, createEn, joinedPath).Data
		b, _ := result.([]byte)
		compiled[file] = string(b)
	}

	req.Respond(compiled, r.Type.JSON)
}

func genRouterCtl(req *r.Request) {
	var compiled = make(map[string]string)
	for _, file := range genRouterFiles {
		joinedPath := filepath.Join("gen", "router", file)
		result := r.RenderContent(r.Type.Text, req.Entity, joinedPath).Data
		b, _ := result.([]byte)
		compiled[file] = string(b)
	}
	req.Respond(compiled, r.Type.JSON)
}

func errorHTMLCtl(req *r.Request) {
	cacheStore, err := r.Storage.Access("cache")
	if err != nil {
		req.Throw(500, r.E("Error accessing storage cache"))
		return
	}

	b := cacheStore.Get("error.html")
	if b != nil {
		req.Respond(string(b), r.Type.Text)
		return
	}
	req.Throw(500, r.E("Error accessing cache/error.html"))
}
