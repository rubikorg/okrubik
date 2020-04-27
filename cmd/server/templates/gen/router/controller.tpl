package {{ .RouterName }}

import r "github.com/rubikorg/rubik"

func indexCtl(en interface{}) r.ByteResponse {
	return r.Success("Hello, this is the {{ .RouterName }} router!")
}