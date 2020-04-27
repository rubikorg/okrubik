package {{ .RouterName }}

import (
	r "github.com/rubikorg/rubik"
)

// Router for /{{ .RouterName }} routes
var Router = r.Create("/{{ .RouterName }}")

var indexRoute = r.Route{
	Path:       "/",
	Controller: indexCtl,
}

func init() {
	Router.Add(indexRoute)
}
