// +build linux darwin

package commands

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/printzero/tint"
	"github.com/radovskyb/watcher"
	"github.com/rubikorg/okrubik/cmd/okrubik/choose"
	"github.com/rubikorg/rubik/pkg"
	"github.com/spf13/cobra"
)

var (
	cmd                *exec.Cmd
	basePath           string
	appName            string
	build              bool
	pluginName         string
	runExtBeforeServer bool
	args               string
)

func initRunCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:     "run",
		Short:   "Runs the app created under this workspace",
		Aliases: []string{"r"},
		Run: func(cmd *cobra.Command, args []string) {
			if pluginName != "" {
				err := runPlugins()
				if err != nil {
					pkg.ErrorMsg(err.Error())
				}

				return
			}
			err := run("")
			if err != nil {
				pkg.ErrorMsg(err.Error())
			}
		},
	}

	runCmd.Flags().StringVarP(&appName, "app", "a", "", "use this flag to run the app/service")
	// TODO: add this feature in the future release
	// runCmd.Flags().BoolVarP(&build, "build", "b", false,
	// 	"build a target binary and run the app/service")
	runCmd.Flags().StringVarP(&pluginName, "plugin", "p", "",
		"use this flags to run Rubik extension blocks")
	runCmd.Flags().StringVarP(&args, "args", "v", "",
		"arguments passed to the plugins/blocks to make it flexible")

	// runCmd.Flags().BoolVarP(&runExtBeforeServer, "run-ext", "", false,
	// 	"use this flags to run Rubik extentions first and start the server ")

	return runCmd

}

// run is a function for running an app from the rubik.toml file
func run(basePath string) error {
	pwd, _ := os.Getwd()
	cfg := pkg.GetRubikConfig()
	if cfg.ProjectName == "" {
		return errors.New("not a valid rubik config")
	}

	w := watcher.New()
	w.SetMaxEvents(1)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	go func() {
		for range done {
			fmt.Println("Shutting down rubik server...")
			killServer()
			fmt.Println("Bye!")
			os.Exit(0)
		}
	}()

	go func() {
		for {
			select {
			case _, ok := <-w.Event:
				if !ok {
					return
				}
				// sleep for 1 sec for changes to get written
				pkg.DebugMsg("Restarting rubik server")
				fmt.Println(
					t.Exp("@(waiting for a second for changes to complete...)", tint.Yellow))
				time.Sleep(time.Second)
				killServer()
				go runServer(basePath)
			case err, ok := <-w.Error:
				if !ok {
					fmt.Println("Something went wrong")
					return
				}
				fmt.Println(err.Error())
			}

		}
	}()

	cpus := runtime.NumCPU()
	if cfg.MaxProcs > 0 && cfg.MaxProcs <= cpus {
		cpus = cfg.MaxProcs
		os.Setenv("GOMAXPROCS", strconv.Itoa(cpus))
		fmt.Println(fmt.Sprintf("GOMAXPROCS set to %d", cpus))
	}

	if basePath == "" {
		var lookup = make(map[string]pkg.Project)
		var options = []string{}
		var answer string

		for _, a := range cfg.App {
			lookup[a.Name] = a
			options = append(options, a.Name)
		}

		prompt := &survey.Select{
			Message: "Select app:",
			Options: options,
		}

		if appName != "" {
			if lookup[appName].Name == "" {
				pkg.ErrorMsg("No such app. Please choose")
				survey.AskOne(prompt, &answer)
			}
			answer = appName
		} else {
			survey.AskOne(prompt, &answer)
		}

		basePath = strings.Replace(lookup[answer].Path, "./", pwd+string(os.PathSeparator), 1)

		fmt.Println(lookup[answer])

		if lookup[answer].Watchable {
			go runServer(basePath)
			startWatcher(w, basePath)
		} else if lookup[answer].RunCommand != "" {
			os.Chdir(basePath)
			customCmd := strings.Split(lookup[answer].RunCommand, " ")
			cmd = exec.Command(customCmd[0], customCmd[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.SysProcAttr = &syscall.SysProcAttr{
				Setpgid: true,
			}
			cmd.Run()
		} else {
			runServer(basePath)
		}

	} else {
		go runServer(basePath)
		startWatcher(w, basePath)
	}

	return nil
}

func startWatcher(w *watcher.Watcher, basePath string) {
	err := w.AddRecursive(basePath)
	if err != nil {
		panic(err)
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func runServer(basePath string) {
	os.Chdir(basePath)
	cmd = exec.Command("go", "run", filepath.Join(basePath, "main.go"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Run()
}

func killServer() {
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		fmt.Println("Cannot kill process", err.Error())
	}
	syscall.Kill(-pgid, syscall.SIGINT)
}

func runPlugins() error {
	var proj pkg.Project
	os.Setenv("RUBIK_ENV", "plugin")
	if appName == "" {
		var err error
		proj, err = choose.RawProject()
		if err != nil || proj.Name == "" {
			return err
		}
	} else {
		cfg := pkg.GetRubikConfig()
		if cfg.ProjectName == "" {
			return errors.New("not a valid rubik config")
		}

		for _, a := range cfg.App {
			if a.Name == appName {
				proj = a
			}
		}

		if proj.Name == "" {
			return fmt.Errorf("%s does not exists in this workspace", appName)
		}
	}

	os.Setenv("RUBIK_PROJ", proj.Name)
	os.Setenv("RUBIK_ARGS", args)
	os.Setenv("RUBIK_PLUGIN", pluginName)
	pwd, _ := os.Getwd()
	path := strings.ReplaceAll(proj.Path, ".", pwd)
	runServer(path)

	os.Setenv("RUBIK_ENV", "")
	os.Setenv("RUBIK_PROJ", "")

	if runExtBeforeServer {
		err := run(path)
		if err != nil {
			return err
		}
	}

	return nil
}
