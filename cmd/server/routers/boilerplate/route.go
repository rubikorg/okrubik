package boilerplate

import (
	"github.com/rubikorg/rubik"
)

// Router for /boilerplate routes
var Router = rubik.Create("/boilerplate")

var createRoute = rubik.Route{
	Path:       "/create",
	Controller: createCtl,
}

func init() {
	Router.Add(createRoute)
	Router.StorageRoutes("route.tpl", "controller.tpl")
}
