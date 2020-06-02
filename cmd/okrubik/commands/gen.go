package commands

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/printzero/tint"
	"github.com/rubikorg/okrubik/cmd/okrubik/choose"
	"github.com/rubikorg/okrubik/pkg/entity"
	"github.com/rubikorg/rubik/pkg"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/ast/astutil"
)

var (
	binName    string
	binPort    string
	routerName string
)

// initGenCmd is code generation method for rubik
// it can generate routers and routes
// and entities
func initGenCmd() *cobra.Command {
	var genCmd = &cobra.Command{
		Use:     "gen",
		Short:   "Generates project code for your Rubik server",
		Aliases: []string{"g", "generate"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(t.Exp(
				"@(Generate needs a subcommand.) Please do okrubik gen --help for more info",
				tint.Magenta))
		},
	}

	genRouterCmd := &cobra.Command{
		Use:   "router",
		Short: "Generate router for an app inside this Rubik workspace",
		Run: func(cmd *cobra.Command, args []string) {
			proj, err := choose.RawProject()
			if err != nil {
				pkg.ErrorMsg(err.Error())
			}

			err = genRouter(proj, routerName)
			if err != nil {
				pkg.ErrorMsg(err.Error())
			}
		},
	}

	genRouterCmd.Flags().StringVarP(
		&routerName, "name", "n", "", "the name of the router/domain you want to generate")
	genRouterCmd.MarkFlagRequired("name")

	genBinCmd := &cobra.Command{
		Use:   "bin",
		Short: "Generate binary inside this Rubik workspace",
		Run: func(cmd *cobra.Command, args []string) {
			err := genBin(binName, binPort)
			if err != nil {
				pkg.ErrorMsg(err.Error())
			}
		},
	}

	genBinCmd.Flags().StringVarP(
		&binName, "name", "n", "", "the binary name of the server you want to generate")
	genBinCmd.Flags().StringVarP(
		&binPort, "port", "p", "", "the port of the server you want to generate")
	genBinCmd.MarkFlagRequired("name")
	genBinCmd.MarkFlagRequired("port")

	genCmd.AddCommand(genBinCmd)
	genCmd.AddCommand(genRouterCmd)

	return genCmd
}

func genRouter(proj pkg.Project, name string) error {
	pwd, _ := os.Getwd()
	path := strings.ReplaceAll(proj.Path, "./", pwd+sep)
	routerPath := filepath.Join(path, "routers", name)
	if f, _ := os.Stat(routerPath); f != nil {
		return errors.New("router with name `" + name + "` already exists")
	}
	// create new folder inside routers
	os.MkdirAll(routerPath, 0755)
	// fetch route and controller
	var files map[string]string
	en := entity.GenRouterEntity{
		RouterName: name,
	}
	en.PointTo = "/boilerplate/gen.router"
	en.Infer = &files
	_, err := rubcl.Get(en)

	if err != nil {
		return err
	}

	// create route and controller file from template
	for k, v := range files {
		fileName := strings.ReplaceAll(k, "tpl", "go")
		err = ioutil.WriteFile(filepath.Join(routerPath, fileName), []byte(v), 0655)
		if err != nil {
			return err
		}
	}

	aerr := addAstRouter(filepath.Join(path, "routers"), name, proj)

	if aerr != nil {
		return aerr
	}

	creationOutput("Generated:", filepath.Join("routers", name))

	return nil
}

func addAstRouter(path, routerName string, proj pkg.Project) error {
	rubikToml := filepath.Join(".", "rubik.toml")
	if f, _ := os.Stat(rubikToml); f == nil {
		return errors.New("Not a rubik project. Cannot find rubik.toml")
	}

	var config pkg.Config
	_, err := toml.DecodeFile(rubikToml, &config)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	importFilePath := filepath.Join(path, "import.go")
	node, err := parser.ParseFile(
		fset, importFilePath, nil, parser.ParseComments)

	if err != nil {
		return err
	}

	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Name.Name == "Import" {
			// we found the import function
			// add import of the newly created router
			importStmt := fmt.Sprintf("%s/cmd/%s/routers/%s", config.Module, proj.Name, routerName)
			astutil.AddImport(fset, node, importStmt)

			expr, err := parser.ParseExpr(fmt.Sprintf("rubik.Use(%s.Router)", routerName))
			if err != nil {
				return err
			}

			callStmt := ast.ExprStmt{
				X: expr,
			}
			fn.Body.List = append(fn.Body.List, &callStmt)

			var buf bytes.Buffer
			err = printer.Fprint(&buf, fset, node)
			if err != nil {
				return err
			}

			formatted, err := format.Source(buf.Bytes())
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(importFilePath, formatted, 0755)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func genBin(name, port string) error {
	rubikToml := filepath.Join(".", "rubik.toml")
	if f, _ := os.Stat(rubikToml); f == nil {
		return errors.New("Not a rubik project. Cannot find rubik.toml")
	}

	var config pkg.Config
	_, err := toml.DecodeFile(rubikToml, &config)
	if err != nil {
		return err
	}

	// TODO: check if bin with same name exists

	newApp := pkg.Project{
		Name:      name,
		Path:      "./cmd/" + name,
		Watchable: true,
	}

	config.App = append(config.App, newApp)
	var files map[string]string
	cbe := entity.CreateBoilerplateEntity{
		Name:       name,
		ModulePath: config.Module,
		Port:       port,
		Bin:        name,
	}
	cbe.PointTo = "/boilerplate/create"
	cbe.Infer = &files

	if _, err := rubcl.Get(cbe); err != nil {
		return err
	}

	basePath := filepath.Join(".", "cmd", name)
	os.MkdirAll(basePath, 0755)

	for name, content := range files {
		truePath := basePath
		namePath := strings.Split(name, "-")
		if !strings.Contains(name, "rubik.toml") {
			for _, p := range namePath {
				if !strings.Contains(p, ".tpl") {
					truePath = filepath.Join(truePath, p)
				}
			}
			os.MkdirAll(truePath, 0755)
			fileName := strings.ReplaceAll(namePath[len(namePath)-1], ".tpl", "")
			filePath := filepath.Join(truePath, fileName)
			creationOutput("creating", filePath)
			err := ioutil.WriteFile(filePath, []byte(content), 0755)
			if err != nil {
				return err
			}
		}
	}

	runTidyCommand(name)

	creationOutput("configuring", rubikToml)

	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(&config); err != nil {
		return err
	}

	if err := ioutil.WriteFile(rubikToml, buf.Bytes(), 0755); err != nil {
		return err
	}

	return nil
}
