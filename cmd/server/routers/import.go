package routers

import (
	"github.com/rubikorg/okrubik/cmd/server/routers/boilerplate"
	"github.com/rubikorg/okrubik/cmd/server/routers/index"
	"github.com/rubikorg/rubik"
)

// Import all routers to rubik
func Import() {
	// index
	rubik.Use(index.Router)
	// boilerplate
	rubik.Use(boilerplate.Router)
}
