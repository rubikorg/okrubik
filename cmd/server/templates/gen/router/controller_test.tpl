package {{ .RouterName }}

import (
	"testing"

	"github.com/rubikorg/rubik"
)

var probe *rubik.TestProbe

func init() {
	probe = rubik.NewProbe(Router)
}

func TestIndexCtl(t *testing.T) {
    // TODO: create your owm entity here for the route and pass it to Test func
	rr := probe.Test(nil)
	if rr.Body.String() != "hello: ashish of age: 22" {
		t.Error("Response: is not `hello: ashish of age: 22`")
	}
}
