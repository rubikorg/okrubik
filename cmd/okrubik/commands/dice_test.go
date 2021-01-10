package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var diceConfigPath string
var modelsPath string

func setupDiceTests() error {
	testpath, err := filepath.Abs(filepath.Join("."))
	if err != nil {
		return err
	}

	srcDiceConfig := filepath.Join("..", "..", "..", "dice.yaml")
	diceConfigPath = filepath.Join(testpath, "dice.yaml")
	modelsPath = filepath.Join(testpath, "models")
	out, err := ioutil.ReadFile(srcDiceConfig)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(diceConfigPath, out, 0744)
	if err != nil {
		return err
	}

	return nil
}

func teardownDiceTests() {
	os.RemoveAll(modelsPath)
	os.Remove(diceConfigPath)
}

func TestMain(m *testing.M) {
	// SETUP
	err := setupDiceTests()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	exitCode := m.Run()
	// TEARDOWN
	teardownDiceTests()
	os.Exit(exitCode)
}

func TestDiceCommand_NoModelPassed(t *testing.T) {
	cmd := initDiceCommand()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.Execute()

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Error(err)
		return
	}

	if string(out) != "insufficient arguments. nothing to do!\n" {
		t.Error("Dice command is working without the model name argument, pleace check!")
	}
}

func TestDiceCommand_ModelPassed(t *testing.T) {
	cmd := initDiceCommand()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--model", "games"})
	cmd.Execute()

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
		return
	}

	if string(out) != "" {
		t.Errorf("there is some out put %s when executing with model name", string(out))
	}
}
