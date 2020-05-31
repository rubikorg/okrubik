package index

import (
	r "github.com/rubikorg/rubik"
)

// Router is index's router
var Router = r.Create("/")

// var corsMw = r.GetBlock(cors.BlockName).(cors.BlockCors).MW

// type iEn struct {
// 	rubik.Entity
// 	Name string
// }

// var indexRoute = r.Route{
// 	Path:   "/",
// 	Entity: &iEn{},
// 	// Entity:     &entity.CreateBoilerplateEntity{},
// 	Controller: indexCtl,
// }

// var testRoute = r.Route{
// 	Path:       "/print/:id",
// 	Entity:     &entity.TestEntity{},
// 	Controller: printCtl,
// }

var installRoute = r.Route{
	Path:       "/install",
	Controller: r.Proxy("https://raw.githubusercontent.com/rubikorg/okrubik/master/install"),
}

func init() {
	// Router.Add(indexRoute)
	// Router.Add(testRoute)
	Router.Add(installRoute)
	// Router.StorageRoutes("gs.zip")
}
