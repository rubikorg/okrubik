package boilerplate

import (
	"github.com/rubikorg/okrubik/pkg/entity"
	"github.com/rubikorg/rubik"
)

// Router for /boilerplate routes
var Router = rubik.Create("/boilerplate")

var createRoute = rubik.Route{
	Path:       "/create",
	Entity:     entity.CreateBoilerplateEntity{},
	Controller: createCtl,
}

func init() {
	Router.Add(createRoute)
}
