package boilerplate

import (
	r "github.com/rubikorg/rubik"
)

var genRouterFiles = []string{
	"route.tpl",
	"controller.tpl",
}

func createCtl(en interface{}) r.ByteResponse {
	// return r.ParseDir("create", en, r.Type.Text)
	return r.Success("success")
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
