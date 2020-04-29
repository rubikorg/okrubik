package entity

import (
	"github.com/rubikorg/rubik"
)

// CreateBoilerplateEntity is used to get all the contents of files required for
// `okrubik create` depending upon the data given by the user
type CreateBoilerplateEntity struct {
	rubik.Entity
	ModulePath string
	Port       string
	Name       string
	Done       bool
}

// GenRouterEntity is used when user executes `okrubik gen router [name]`
type GenRouterEntity struct {
	rubik.Entity
	RouterName string
}
