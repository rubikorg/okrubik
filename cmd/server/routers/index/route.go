package index

import (
	"fmt"
	"net/http"

	"github.com/rubikorg/okrubik/pkg/entity"
	r "github.com/rubikorg/rubik"
)

// Router is index's router
var Router = r.Create("/")

type iEn struct {
	r.Entity
	Name string `rubik:"awesome|body"`
	Age  int
}

func (ien iEn) ComposedEntity() r.Entity {
	return ien.Entity
}

func (ien iEn) CoreEntity() interface{} {
	return ien
}

func (ien iEn) Path() string {
	return ien.Entity.PointTo
}

func (i *iEn) getNamePlusAge() string {
	return fmt.Sprintf("%s+%d", i.Name, i.Age)
}

func printUid(req *r.Request) {
	fmt.Println(req.Claims)
}

func init() {
	indexRoute := r.Route{
		Path:       "/",
		Controller: indexCtl,
		// Validation: r.Validation{
		// 	"Name": []r.Assertion{
		// 		checker.MustExist,
		// 		checker.StrIsOneOf("ash", "ashish"),
		// 	},
		// },
	}
	Router.Add(indexRoute)

	installRoute := r.Route{
		Path:       "/install",
		Controller: r.Proxy("https://raw.githubusercontent.com/rubikorg/okrubik/master/install"),
	}
	Router.Add(installRoute)

	getTokenRoute := r.Route{
		Path:       "/token",
		Controller: getTokenCtl,
	}
	Router.Add(getTokenRoute)

	someRoute := r.Route{
		Method:     http.MethodPost,
		Path:       "/some",
		Controller: someCtl,
	}
	Router.Add(someRoute)
	facebookRoute := r.Route{
		Path:       "/facebook",
		Controller: facebookCtl,
	}
	Router.Add(facebookRoute)
	addSaltRoute := r.Route{
		Path:       "/addSalt",
		Entity:     entity.AddSaltEntity{},
		Controller: addSaltCtl,
	}
	Router.Add(addSaltRoute)

	// Router.StorageRoutes("gs.zip")
}
