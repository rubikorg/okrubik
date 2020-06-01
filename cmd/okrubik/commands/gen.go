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
	"github.com/rubikorg/okrubik/cmd/okrubik/choose"
	"github.com/rubikorg/okrubik/pkg/entity"
	"github.com/rubikorg/rubik/pkg"
	"golang.org/x/tools/go/ast/astutil"
)

// Gen is code generation method for rubik
// it can generate routers and routes
// and entities
func Gen(args []string) error {
	if len(args) == 0 {
		return errors.New("gen command requires arguments")
	}

	switch args[0] {
	case "bin":
		if len(args) < 3 {
			return errors.New("bin requires a name of the binary and a port")
		}
		return genBin(args[1], args[2])
	case "router":
		if len(args) == 1 {
			return errors.New("router requires a name to initialize")
		}
		// select the project using the project selector
		proj, err := choose.RawProject()
		if err != nil {
			return err
		}
		return genRouter(proj, args[1])
	case "route":
		// check if router name is given as argument
		// use ast to write a new route
		// add it inside init()
		// create a new controller inside controller.go
		break
	case "entity":
		// check if name of entity given
		// loop until user enters text "done"
		break
	}
	return nil
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
