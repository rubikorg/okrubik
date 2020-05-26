package boilerplate

import (
	"path/filepath"

	r "github.com/rubikorg/rubik"
)

var createFiles = []string{
	"config-decl.go.tpl",
	"routers-index-controller.go.tpl",
	"main.go.tpl",
	"rubik.toml.tpl",
	"templates-index.html.tpl",
	"config-default.toml.tpl",
	"routers-index-route.go.tpl",
	"routers-import.go.tpl",
	"static-main.css.tpl",
	"README.md.tpl",
}

var genRouterFiles = []string{
	"route.tpl",
	"controller.tpl",
}

func createCtl(en interface{}) r.ByteResponse {
	var compiled = make(map[string]string)
	for _, file := range createFiles {
		joinedPath := filepath.Join("create", file)
		result := r.Render(r.Type.Text, en, joinedPath).Data
		b, _ := result.([]byte)
		compiled[file] = string(b)
	}
	return r.Success(compiled, r.Type.JSON)
}

func genRouterCtl(en interface{}) r.ByteResponse {
	var compiled = make(map[string]string)
	for _, file := range genRouterFiles {
		joinedPath := filepath.Join("gen", "router", file)
		result := r.Render(r.Type.Text, en, joinedPath).Data
		b, _ := result.([]byte)
		compiled[file] = string(b)
	}
	return r.Success(compiled, r.Type.JSON)
}

func errorHTMLCtl(en interface{}) r.ByteResponse {
	return r.Render(r.Type.Text, nil, "error.html.tpl")
}
