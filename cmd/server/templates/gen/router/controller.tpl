package {{ .RouterName }}

import r "github.com/rubikorg/rubik"

func indexCtl(req *r.Request) {
	req.Respond("Hello, this is the {{ .RouterName }} router!")
}