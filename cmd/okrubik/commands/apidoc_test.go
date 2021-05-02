package commands

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestOpenAPIDoc(t *testing.T) {
	b := bytes.NewBufferString("")
	cmd := initApiDocCommand()
	cmd.SetErr(b)
	cmd.Execute()

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
		return
	}

	if string(out) != "" {
		t.Errorf("An error occured while opening api doc url %s", string(out))
	}
}
