package index

import (
	"github.com/rubikorg/rubik"
)

// Router is index's router
var Router = rubik.Create("/")

var indexRoute = rubik.Route{
	Path:       "/",
	Controller: indexCtl,
}

func init() {
	Router.Add(indexRoute)
}
