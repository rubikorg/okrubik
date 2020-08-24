package {{ .RouterName }}

import (
	r "github.com/rubikorg/rubik"
)

// Router for /{{ .RouterName }} routes
var Router = r.Create("/{{ .RouterName }}")

func init() {
	indexRoute = r.Route{
		Path:       "/",
		Controller: indexCtl,
	}
	
	Router.Add(indexRoute)
}
