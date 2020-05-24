package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/rubikorg/rubik/pkg"
)

// Exec executes the `okrubik x:test` commands
func Exec(arg string) error {
	cmdName := strings.Split(arg, ":")
	if len(cmdName) < 2 || cmdName[1] == "" {
		return errors.New("x requires a command after the : symbol")
	}

	var rubikConf pkg.Config
	tomlPath := filepath.Join(".", "rubik.toml")
	if f, _ := os.Stat(tomlPath); f == nil {
		return errors.New("not a Rubik workspace")
	}

	_, err := toml.DecodeFile(tomlPath, &rubikConf)
	if err != nil {
		return err
	}

	if rubikConf.X[cmdName[1]] == nil {
		return errors.New("no such command defined in rubik.toml")
	}

	cmd := rubikConf.X[cmdName[1]]["command"]
	if cmd == "" {
		return errors.New("a command defined inside x requires `command` variable to be defined")
	}

	var cmds = []string{}
	pwd := rubikConf.X[cmdName[1]]["pwd"]

	if strings.Contains(cmd, "&&") {
		cmds = strings.Split(cmd, "&&")
	} else {
		cmds = []string{cmd}
	}

	if pwd != "" {
		pwd = strings.TrimPrefix(pwd, "/")
		joinedFp := []string{"."}
		for _, p := range strings.Split(pwd, "/") {
			joinedFp = append(joinedFp, p)
		}
		pwdPath := filepath.Join(joinedFp...)
		os.Chdir(pwdPath)
	}

	// TODO: fix multiple command working
	for _, c := range cmds {
		fullCmd := strings.Split(strings.Trim(c, ""), " ")
		fmt.Println(fullCmd[0], fullCmd[1:])
		cmdToRun := exec.Command(fullCmd[0], fullCmd[1:]...)
		cmdToRun.Stdout = os.Stdout
		cmdToRun.Stderr = os.Stderr
		cmdToRun.Run()
	}

	return nil
}
