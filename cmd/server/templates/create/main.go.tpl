package main

import (
	cfg "{{ .ModulePath }}/cmd/server/config"
	"{{ .ModulePath }}/cmd/server/routers"
	r "github.com/rubikorg/rubik"
)

func main() {
	var config cfg.ProjectConfig
	err := r.Load(&config)

	if err != nil {
		panic(err)
	}

	routers.Import()
	panic(r.Run())

}
