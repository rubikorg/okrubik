package routers

import (
	"github.com/oksketch/oksketch/cmd/server/routers/index"
	"github.com/oksketch/sketch"
)

func Import() {
	// index
	sketch.Use(index.Router)

}
