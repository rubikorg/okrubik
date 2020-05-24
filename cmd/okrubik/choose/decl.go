package choose

import (
	"errors"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rubikorg/rubik/pkg"
)

// Project selection prompt
func Project() (string, error) {
	pwd, _ := os.Getwd()
	proj, _ := RawProject()
	return strings.Replace(proj.Path, "./", pwd+string(os.PathSeparator), 1), nil
}

func RawProject() (pkg.Project, error) {
	cfg := pkg.GetRubikConfig()
	if cfg.ProjectName == "" {
		return pkg.Project{}, errors.New("not a valid rubik config")
	}

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

	return lookup[answer], nil
}
