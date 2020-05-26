package index

import (
	r "github.com/rubikorg/rubik"
)

func indexCtl(en interface{}) r.ByteResponse {
	// {
	// 	swagger: "2.0",
	// 	info: {
	// 	description: "This is a sample server Petstore server. You can find out more about Swagger at [http://swagger.io](http://swagger.io) or on [irc.freenode.net, #swagger](http://swagger.io/irc/). For this sample, you can use the api key `special-key` to test the authorization filters.",
	// 	version: "1.0.5",
	// 	title: "Swagger Petstore",
	// 	termsOfService: "http://swagger.io/terms/",
	// 	contact: {
	// 	email: "apiteam@swagger.io"
	// 	},
	// 	license: {
	// 	name: "Apache 2.0",
	// 	url: "http://www.apache.org/licenses/LICENSE-2.0.html"
	// 	}
	// 	},
	im := map[string]string{
		"title":       "server",
		"description": "awesome test",
	}
	var m = map[string]interface{}{
		"swagger": "2.0",
		"info":    im,
	}
	return r.Success(m, r.Type.JSON)
}

func swaggerCtl(en interface{}) r.ByteResponse {
	return r.Render(r.Type.Text, en, "swagger.html")
}

func installCtl(en interface{}) r.ByteResponse {
	return r.Redirect("https://raw.githubusercontent.com/rubikorg/okrubik/master/install")
}

