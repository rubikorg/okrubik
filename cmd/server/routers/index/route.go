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

var installRoute = r.Route{
	Path:       "/install",
	Controller: installCtl,
}

func init() {
	Router.Add(indexRoute)
	Router.Add(installRoute)
	// Router.StorageRoutes("gs.zip")
}
