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
	cfg := pkg.GetRubikConfig()
	if cfg.ProjectName == "" {
		return "", errors.New("not a valid rubik config")
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

	return strings.Replace(lookup[answer].Path, "./", pwd+string(os.PathSeparator), 1), nil
}
