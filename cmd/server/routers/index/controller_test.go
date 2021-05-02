package index

import (
	"testing"

	"github.com/rubikorg/rubik"
)

var probe *rubik.TestProbe

func init() {
	probe = rubik.NewProbe(Router)
}

func TestIndexCtl(t *testing.T) {
	en := iEn{
		Name: "ashish",
		Age:  22,
	}
	en.PointTo = "/"
	rr := probe.Test(en)
	resp := rr.Body.String()
	if resp != "hello: ashish of age: 22" {
		t.Error("Response")
	}
}
