package index

import (
	r "github.com/rubikorg/rubik"
)

func indexCtl(en interface{}) r.ByteResponse {
	return r.Render(r.Type.HTML, nil, "index.html")
}
