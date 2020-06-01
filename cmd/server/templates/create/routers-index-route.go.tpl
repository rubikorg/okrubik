package index

import (
	r "github.com/rubikorg/rubik"
)

// Router is index's router
var Router = r.Create("/")

var indexRoute = r.Route{
	Path:       "/",
	Controller: indexCtl(),
}

func init() {
	Router.Add(indexRoute)
}
