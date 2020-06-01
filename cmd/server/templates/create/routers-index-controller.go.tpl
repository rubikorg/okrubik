package index

import (
	r "github.com/rubikorg/rubik"
)

func indexCtl() r.Controller {
	// Generally r.Render method satisfies rubik.Controller signature
	// so if you want to render any template without dynamic
	// variable you can directly use r.Render method inside
	// Controller field of your rubik Route
	return r.Render(r.Type.HTML, nil, "index.html")
}
