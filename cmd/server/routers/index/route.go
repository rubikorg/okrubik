package index

import "github.com/rubikorg/rubik"

var Router = rubik.Create("/")

var indexRoute = rubik.Route{
	Path:       "/",
	Controller: indexCtl,
}

func init() {
	Router.Add(indexRoute)
	Router.StorageRoutes("gs.zip")
}
