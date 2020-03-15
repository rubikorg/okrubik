package routers

import (
	"github.com/rubikorg/okrubik/cmd/server/routers/index"
	"github.com/rubikorg/rubik"
)

func Import() {
	// index
	rubik.Use(index.Router)

}
