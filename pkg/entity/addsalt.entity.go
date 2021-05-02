package entity

import "github.com/rubikorg/rubik"

// AddSaltEntity implements rubik.TestableEntity
type AddSaltEntity struct {
	rubik.Entity
}

func (en AddSaltEntity) ComposedEntity() rubik.Entity {
	return en.Entity
}

func (en AddSaltEntity) CoreEntity() interface{} {
	return en
}

func (en AddSaltEntity) Path() string {
	return en.PointTo
}
