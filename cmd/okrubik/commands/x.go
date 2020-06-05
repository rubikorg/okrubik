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
	"github.com/spf13/cobra"
)

func initExecCmd() *cobra.Command {
	execCmd := &cobra.Command{
		Use:     "exec",
		Short:   "Execute a rubik command defined inside rubik.toml under [x] object",
		Aliases: []string{"x"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				pkg.ErrorMsg("exec/x requires a name to execute a command")
				return
			}

			err := x(args[0])
			if err != nil {
				pkg.ErrorMsg(err.Error())
			}
		},
	}
	return execCmd
}

// x executes the `okrubik x:test` commands
func x(arg string) error {
	var rubikConf pkg.Config
	tomlPath := filepath.Join(".", "rubik.toml")
	if f, _ := os.Stat(tomlPath); f == nil {
		return errors.New("not a Rubik workspace")
	}

	_, err := toml.DecodeFile(tomlPath, &rubikConf)
	if err != nil {
		return err
	}

	if rubikConf.X[arg] == nil {
		return errors.New("no such command defined in rubik.toml")
	}

	cmd := rubikConf.X[arg]["command"]
	if cmd == "" {
		return errors.New("a command defined inside x requires `command` variable to be defined")
	}

	var cmds = []string{}
	pwd := rubikConf.X[arg]["pwd"]

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
