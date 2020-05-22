package index

import (
	"github.com/rubikorg/blocks/cors"
	"github.com/rubikorg/blocks/guard"
	r "github.com/rubikorg/rubik"
)

// Router is index's router
var Router = r.Create("/")
var corsMw = r.GetBlock(cors.BlockName).(cors.BlockCors).MW

var indexRoute = r.Route{
	Guard: guard.BasicGuard{},
	Path:  "/",
	Middlewares: []r.Middleware{
		corsMw(),
	},
	// Entity:     &entity.CreateBoilerplateEntity{},
	Controller: indexCtl,
}

func init() {
	Router.Add(indexRoute)
	// Router.StorageRoutes("gs.zip")
}
