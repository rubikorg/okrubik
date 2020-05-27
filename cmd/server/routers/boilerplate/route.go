package boilerplate

import (
	"github.com/rubikorg/okrubik/pkg/entity"
	"github.com/rubikorg/rubik"
)

// Router for /boilerplate routes
var Router = rubik.Create("/boilerplate")

var createRoute = rubik.Route{
	Path:       "/create",
	Entity:     &entity.CreateBoilerplateEntity{},
	Controller: createCtl,
	ResponseDeclarations: map[int]string{
		200: "object",
	},
}

var genRouterRoute = rubik.Route{
	Path:       "/gen.router",
	Entity:     &entity.GenRouterEntity{},
	Controller: genRouterCtl,
}

var errorHTMLRoute = rubik.Route{
	Path:       "error.html",
	Controller: errorHTMLCtl,
}

func init() {
	Router.Add(createRoute)
	Router.Add(genRouterRoute)
	Router.Add(errorHTMLRoute)
}
