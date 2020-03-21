package index

import (
	"github.com/rubikorg/rubik"
)

// Router is index's router
var Router = rubik.Create("/")

var indexRoute = rubik.Route{
	Path: "/",
	// Middlewares: []rubik.Middleware{
	// 	guard.JWT(authenticate),
	// },
	Controller: indexCtl,
}

func authenticate(token string) interface{} {
	if token == "ashish" {
		return "success"
	}
	return "failed"
}

func init() {
	Router.Add(indexRoute)
	Router.StorageRoutes("gs.zip")
}
