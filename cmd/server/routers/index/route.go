package index

import "github.com/oksketch/sketch"

var Router = cherry.Create("/")

var indexRoute = cherry.Route{
	Path:       "/",
	Controller: indexCtl,
}

func init() {
	Router.Add(indexRoute)
	Router.StorageRoutes("gs.zip")
}
