package index

import "github.com/okcherry/cherry"

var Router = cherry.App.Create("/")

var indexRoute = cherry.Route{
	Path:       "/",
	Controller: indexCtl,
}

func init() {
	Router.Add(indexRoute)
}
