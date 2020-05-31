package index

import (
	r "github.com/rubikorg/rubik"
)

func indexCtl(req *r.Request) {
	req.Respond("hello go")
}
