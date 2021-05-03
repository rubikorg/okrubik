package commands

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rubikorg/rubik/pkg"
	"github.com/spf13/cobra"
)

var (
	allFlag         bool
	specificService string
)

func initTestCmd() *cobra.Command {
	testCmd := cobra.Command{
		Use:     "test",
		Short:   "Run the tests inside your services",
		Aliases: []string{"t"},
		Run: func(cmd *cobra.Command, args []string) {
			if allFlag {
				err := runAllTests()
				if err != nil {
					pkg.ErrorMsg(err.Error())
				}
			} else if specificService != "" {
				err := runTestsForSpecificService()
				if err != nil {
					pkg.ErrorMsg(err.Error())
				}
			}
		},
	}

	testCmd.Flags().BoolVarP(&allFlag, "all", "a", false,
		"Use this flag to run tests for all the services")
	testCmd.Flags().StringVarP(&specificService, "service", "s", "",
		"Use this flag to run tests for a specific service")
	return &testCmd
}

func runAllTests() error {
	rconf, err := pkg.GetRubikConfig()
	if err != nil {
		return err
	}

	for _, service := range rconf.App {
		pkg.EmojiMsg("🚀", service.Name)
		servicePath := filepath.Join(service.Path, "...")
		cmd := exec.Command("go", "test", "-v", "./"+servicePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	return nil
}

func runTestsForSpecificService() error {
	rconf, err := pkg.GetRubikConfig()
	if err != nil {
		return err
	}
	var proj pkg.Project
	for _, service := range rconf.App {
		if service.Name == specificService {
			proj = service
		}
	}

	if proj.Name == "" {
		return errors.New("No such service")
	}

	pkg.EmojiMsg("🚀", proj.Name)
	servicePath := filepath.Join(proj.Path, "...")
	cmd := exec.Command("go", "test", "-v", "./"+servicePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	return nil
}
