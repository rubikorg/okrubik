package main

import (
	cfg "{{ .ModulePath }}/cmd/{{ .Bin }}/config"
	"{{ .ModulePath }}/cmd/{{ .Bin }}/routers"
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
