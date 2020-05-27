package commands

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/radovskyb/watcher"
	"github.com/rubikorg/rubik/pkg"
)

const (
	sep = string(os.PathSeparator)
)

var cmd *exec.Cmd
var basePath string

// Run is a function for running an app from the rubik.toml file
func Run() error {
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
			fmt.Println("Done!")
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
				pkg.DebugMsg("Restarting rubik server")
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

	survey.AskOne(prompt, &answer)

	basePath = strings.Replace(lookup[answer].Path, "./", pwd+sep, 1)

	cpus := runtime.NumCPU()
	if cfg.MaxProcs > 0 && cfg.MaxProcs <= cpus {
		cpus = cfg.MaxProcs
		os.Setenv("GOMAXPROCS", strconv.Itoa(cpus))
		fmt.Println(fmt.Sprintf("GOMAXPROCS set to %d", cpus))
	}

	if lookup[answer].Watchable {
		go runServer(basePath)
		err := w.AddRecursive(basePath)
		if err != nil {
			panic(err)
		}

		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Fatalln(err)
		}
	} else {
		runServer(basePath)
	}
	return nil
}

func runServer(basePath string) {
	// fmt.Println("Setting new commnd", cmd.Process.Pid)
	os.Chdir(basePath)
	cmd = exec.Command("go", "run", basePath+sep+"main.go")
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
