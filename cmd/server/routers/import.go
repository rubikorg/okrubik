package routers

import (
	"github.com/oksketch/oksketch/cmd/server/routers/index"
	"github.com/oksketch/sketch"
)

func Import() {
	// index
	cherry.Use(index.Router)

}
