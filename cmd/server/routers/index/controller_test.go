package index

import (
	"testing"

	"github.com/rubikorg/rubik"
)

var probe rubik.TestProbe

func init() {
	probe = rubik.NewProbe(Router)
}

func TestIndexCtl(t *testing.T) {
	_, rr := probe.Test("GET", "/", nil, nil, indexCtl)
	if rr.Result().StatusCode != 200 {
		t.Error("The status code for indexCtl is not 200:", rr.Body.String())
	}

	if rr.Body.String() != "hello go" {
		t.Errorf("Wrong response, curent resp: %s", rr.Body.String())
	}
}
