package index

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	jwtBlock "github.com/rubikorg/blocks/guard/jwt"
	r "github.com/rubikorg/rubik"
)

type tokenStruct struct {
	Token string `json:"token"`
}

func indexCtl(req *r.Request) {
	en := req.Entity.(iEn)
	r.Log.I <- "this is an info message"
	req.Respond(fmt.Sprintf("hello: %s of age: %d", en.Name, en.Age))
}

func getTokenCtl(req *r.Request) {
	claims := make(jwt.MapClaims)
	claims["uid"] = 1
	token, err := jwtBlock.CreateToken(claims, true)
	if err != nil {
		req.Throw(400, err)
		return
	}

	req.Respond(tokenStruct{token}, r.Type.JSON)
}

func someCtl(req *r.Request) {
	req.Respond("I am some controller!")
}
func facebookCtl(req *r.Request) {
	req.Respond("I am facebook controller!")
}
func addSaltCtl(req *r.Request) {
	req.Respond("I am addSalt controller!")
}
