package boilerplate

import (
	r "github.com/rubikorg/rubik"
)

func createCtl(en interface{}) r.ByteResponse {
	// return r.ParseDir("create", en, r.Type.Text)
	return r.Success("success")
}
