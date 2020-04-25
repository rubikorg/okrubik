package {{ .RouterName }}

import r "github.com/rubikrog/rubik"

func indexCtl(en interface{}) r.ByteResponse {
	return r.Success("Hello, this is the {{ .RouterName }} router!")
}