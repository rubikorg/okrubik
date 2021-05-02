package routers

import (
	"{{ .ModulePath }}/cmd/{{ .Bin }}/routers/index"
	"github.com/rubikorg/rubik"
)

// Import is a single point of contact for all routers into rubik
func Import() {
	// index
	rubik.Use(index.Router)
}
