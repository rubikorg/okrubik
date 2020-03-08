package routers

import (
	"github.com/okcherry/cherry"
	"github.com/okcherry/okcherry/cmd/server/routers/index"
)

func Import() {
	// index
	cherry.App.Use(index.Router)

}
