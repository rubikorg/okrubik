package index

import (
	"github.com/rubikorg/rubik"
)

func indexCtl(en interface{}) (interface{}, error) {
	return rubik.Render("index.html", nil, rubik.Type.HTML)
}