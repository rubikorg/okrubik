package main

import (
	"{{ .ModulePath }}/cmd/{{ .Bin }}/app"
	"{{ .ModulePath }}/cmd/{{ .Bin }}/routers"
	"{{ .ModulePath }}/pkg/services"
	r "github.com/rubikorg/rubik"
)

func main() {
	var config app.ProjectConfig
	err := r.Load(&config)

	if err != nil {
		panic(err)
	}

	// TODO: set your one-time application level dependency here
	// eg: DB Connection, Logger etc..
	app.SetDep(app.Dependency{
		ProjectConfig: config,
	})

	routers.Import()
	if err := r.Run(services.{{ .CapBin }}Service); err != nil {
		panic(err)
	}
}
