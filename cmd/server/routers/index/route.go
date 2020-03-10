package index

import "github.com/oksketch/sketch"

var Router = sketch.Create("/")

var indexRoute = sketch.Route{
	Path:       "/",
	Controller: indexCtl,
}

func init() {
	Router.Add(indexRoute)
	Router.StorageRoutes("gs.zip")
}
