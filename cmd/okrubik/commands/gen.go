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
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/BurntSushi/toml"
	"github.com/printzero/tint"
	"github.com/rubikorg/okrubik/cmd/okrubik/choose"
	"github.com/rubikorg/okrubik/pkg/entity"
	"github.com/rubikorg/rubik/pkg"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/ast/astutil"
)

var (
	binName     string
	binPort     string
	routerName  string
	routeName   string
	entityName  string
	genAppName  string
	writeEntity bool
)

// ControllerTestTemplate is used to generate a test function
var ControllerTestTemplate = `
func Test{{ .Name }}Route(t *testing.T) {
	// TODO: replace this with your own Entity
	testEntity := struct {
		rubik.Entity
	}
	testEntity.PointTo = "/"
	rr := probe.Test(testEntity)
	if rr.Body.String() != "hello rubik" {
		t.Error("I'm a failing test!!")
	}
}`

// EntityTemplate is used to create an entity file from the okrubik gen command
var EntityTemplate = `package entity

import "github.com/rubikorg/rubik"

// {{ .N }}Entity implements rubik.TestableEntity
type {{ .N }}Entity struct {
	rubik.Entity
	{{- range $k, $v := .M }}
	{{ $k }} {{ $v }}
	{{- end }}
}

func (en {{ .N }}Entity) ComposedEntity() rubik.Entity {
	return en.Entity
}

func (en {{ .N }}Entity) CoreEntity() interface{} {
	return en
}

func (en {{ .N }}Entity) Path() string {
	return en.PointTo
}
`

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
				"Generate needs a subcommand. Please do @(okrubik gen --help) for more info",
				tint.BgMagenta.Bold().Add(tint.White)))
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

	genServiceCmd := &cobra.Command{
		Use:   "service",
		Short: "Generate service binary inside this Rubik workspace",
		Run: func(cmd *cobra.Command, args []string) {
			err := genBin(binName, binPort)
			if err != nil {
				pkg.ErrorMsg(err.Error())
			}
		},
	}

	genServiceCmd.Flags().StringVarP(
		&binName, "name", "n", "", "the binary name of the server you want to generate")
	genServiceCmd.Flags().StringVarP(
		&binPort, "port", "p", "", "the port of the server you want to generate")
	genServiceCmd.MarkFlagRequired("name")
	genServiceCmd.MarkFlagRequired("port")

	genRouteCmd := &cobra.Command{
		Use:   "route",
		Short: "Generate new route for Rubik app",
		Run: func(cmd *cobra.Command, args []string) {
			if routerName == "" || routeName == "" {
				pkg.ErrorMsg("generating a 'route' requires -r and -n arguments. " +
					"do 'okrubik gen route --help' for more info")
				return
			}
			proj, err := choose.RawProject()
			err = genRoute(proj)
			if err != nil {
				pkg.ErrorMsg(err.Error())
			}
		},
	}

	genRouteCmd.Flags().StringVarP(
		&routerName, "router", "r", "", "router name to generate new route")
	genRouteCmd.Flags().StringVarP(
		&routeName, "name", "n", "", "name of the route")
	// TODO: we need to enable this feature in the future :P (now always create new entity)
	// genRouteCmd.Flags().BoolVarP(
	// 	&writeEntity, "entity", "e", false, "flag to generate entity with the route")
	// genRouteCmd.Flags().StringVarP(
	// 		&routeName, "input-entity", "i", "",
	// 		"name of the existing entity present inside the pkg/entity folder")

	genEntityCmd := &cobra.Command{
		Use:   "entity",
		Short: "Generate a new entity for a Rubik API",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 || len(args) < 2 {
				pkg.ErrorMsg("Entity generator requires you to specify Entity name and Entity data")
				return
			}

			err := genEntity(args[0], genAppName, args[1])
			if err != nil {
				pkg.ErrorMsg(err.Error())
			}
		},
	}

	genEntityCmd.Flags().StringVarP(
		&genAppName, "service", "s", "", "the service for which you want to generate the Entity")

	// TODO: uncomment this when you implement adding a base template to the test file
	// genTestCmd := &cobra.Command{
	// 	Use:   "test",
	// 	Short: "Generate a new test for your controller",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		if len(args) == 0 || len(args) < 2 {
	// 			pkg.ErrorMsg("Entity generator requires you to specify Entity name and Entity data")
	// 			return
	// 		}

	// 		err := genTest(args[0], genAppName)
	// 		if err != nil {
	// 			pkg.ErrorMsg(err.Error())
	// 		}
	// 	},
	// }

	// genTestCmd.Flags().StringVarP(
	// 	&routerName, "router", "r", "", "the router for which controller test is to be written")

	// genTestCmd.Flags().StringVarP(
	// 	&genAppName, "service", "s", "", "the service for which test to be generated")

	genCmd.AddCommand(genServiceCmd)
	genCmd.AddCommand(genRouterCmd)
	genCmd.AddCommand(genRouteCmd)
	genCmd.AddCommand(genEntityCmd)
	// genCmd.AddCommand(genTestCmd)

	return genCmd
}

func genRouter(proj pkg.Project, name string) error {
	pwd, _ := os.Getwd()
	path := strings.ReplaceAll(proj.Path, ".", pwd)
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

	aerr := addRouterToAST(filepath.Join(path, "routers"), name, proj)

	if aerr != nil {
		return aerr
	}

	creationOutput("Generated:", filepath.Join("routers", name))

	return nil
}

func addRouterToAST(path, routerName string, proj pkg.Project) error {
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

func genRoute(proj pkg.Project) error {
	routerPath := path.Join(".", "cmd", proj.Name, "routers", routerName)
	if fi, _ := os.Stat(routerPath); fi == nil {
		msg := fmt.Sprintf("could not find router %s", routerName)
		return errors.New(msg)
	}

	routeFilePath := path.Join(routerPath, "route.go")
	ctlFilePath := path.Join(routerPath, "controller.go")

	// check if a rubik.Route with same name is present inside the given router
	if exists, err := checkIfRouteExists(routeFilePath, routeName); exists {
		if err != nil {
			return err
		}
		return fmt.Errorf("%sRoute already exists", routeName)
	}

	err := genEntity(routeName, proj.Name, "")
	if err != nil {
		return err
	}

	err = addRouteToAST(routeFilePath, routeName)
	if err != nil {
		return err
	}

	err = addControllerToAST(ctlFilePath, routeName)
	if err != nil {
		return err
	}

	colored := t.Exp(fmt.Sprintf("@(%s Route) generated inside @(%s)!", routeName, routerName),
		tint.Magenta, tint.Green.Bold())
	fmt.Println(colored)

	return nil
}

func checkIfRouteExists(routeFilePath, rrouteName string) (bool, error) {
	// TODO: fset and parsing node is same accross all funcs ..extract it!
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, routeFilePath, nil, parser.ParseComments)
	if err != nil {
		return true, err
	}

	var fileBuf bytes.Buffer
	err = printer.Fprint(&fileBuf, fset, node)
	if err != nil {
		return true, err
	}

	fileContents := string(fileBuf.Bytes())

	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if fn.Name.Name == "init" {
			for _, s := range fn.Body.List {
				if expr, ok := s.(*ast.ExprStmt); ok {
					if strings.Contains(fileContents[expr.X.Pos():expr.X.End()],
						fmt.Sprintf("%sRoute", rrouteName)) {
						return true, nil
					}
				}
			}
		}
	}
	return false, nil
}

func addRouteToAST(routeFilePath, rrouteName string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, routeFilePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	rconf := pkg.GetRubikConfig()
	importStmt := fmt.Sprintf("%s/pkg/entity", rconf.Module)
	astutil.AddImport(fset, node, importStmt)

	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Name.Name == "init" {
			routeDecl := fmt.Sprintf(
				`r.Route{Path: "/%s", Entity: entity.%sEntity{}, Controller: %sCtl}`,
				rrouteName, capitalize(rrouteName), rrouteName)

			expr, err := parser.ParseExpr(routeDecl)
			routeAssignStmt := ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: fmt.Sprintf("%sRoute", rrouteName),
					},
				},
				Rhs: []ast.Expr{expr},
			}

			addExpr, err := parser.ParseExpr(fmt.Sprintf("Router.Add(%sRoute)", rrouteName))
			if err != nil {
				return err
			}
			rubikAddStmt := ast.ExprStmt{
				X: addExpr,
			}

			fn.Body.List = append(fn.Body.List, &routeAssignStmt)
			fn.Body.List = append(fn.Body.List, &rubikAddStmt)
			if err != nil {
				return err
			}
			var buf bytes.Buffer
			err = printer.Fprint(&buf, fset, node)
			if err != nil {
				return err
			}

			formatted, err := format.Source(buf.Bytes())
			err = ioutil.WriteFile(routeFilePath, formatted, 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func addControllerToAST(controllerFilePath, ctlRouteName string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, controllerFilePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	ctlParamIdent := []*ast.Ident{
		{
			Name: "req",
		},
	}

	reqTypeExpr, err := parser.ParseExpr("*r.Request")
	if err != nil {
		return err
	}

	fnBodyExpr, err := parser.ParseExpr(fmt.Sprintf(`req.Respond("I am %s controller!")`,
		ctlRouteName))

	ctlIncomingFields := []*ast.Field{
		{
			Names: ctlParamIdent,
			Type:  reqTypeExpr,
		},
	}
	ctlDecl := ast.FuncDecl{
		Name: &ast.Ident{Name: fmt.Sprintf("%sCtl", ctlRouteName)},
		Recv: nil,
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: ctlIncomingFields,
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: fnBodyExpr,
				},
			},
		},
	}

	node.Decls = append(node.Decls, &ctlDecl)
	var buf bytes.Buffer
	err = printer.Fprint(&buf, fset, node)
	if err != nil {
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	err = ioutil.WriteFile(controllerFilePath, formatted, 0755)
	if err != nil {
		return err
	}
	return nil
}

func genEntity(en, app, data string) error {
	var proj pkg.Project
	// TODO: this whole if-else can be abstracted out into a function
	if app == "" {
		var err error
		proj, err = choose.RawProject()
		if err != nil {
			return err
		}
	} else {
		config := pkg.GetRubikConfig()
		if config.Module == "" {
			return errors.New("not a valid Rubik project")
		}

		for _, a := range config.App {
			if a.Name == app {
				proj = a
			}
		}

		if proj.Name == "" {
			return fmt.Errorf("Project: %s not found in this Rubik workspace", app)
		}
	}

	pwd, _ := os.Getwd()
	entityPath := filepath.Join(pwd, "pkg", "entity", fmt.Sprintf(
		"%s.entity.go", strings.ToLower(en)))
	structData := parseEntityData(data)
	tpl, err := template.New("entity").Parse(EntityTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	type tplData = struct {
		N string
		M map[string]string
	}
	err = tpl.Execute(&buf, tplData{N: capitalize(en), M: structData})
	if err != nil {
		return err
	}

	if f, _ := os.Stat(filepath.Join(pwd, "pkg", "entity")); f == nil {
		os.MkdirAll(filepath.Join(pwd, "pkg", "entity"), 0755)
	}

	err = ioutil.WriteFile(entityPath, buf.Bytes(), 0755)
	if err != nil {
		return err
	}

	return nil
}

func genTest(testName, app string) error {
	return nil
}

func parseEntityData(dataArg string) map[string]string {
	data := make(map[string]string)
	if strings.Contains(dataArg, ",") {
		members := strings.Split(dataArg, ",")
		for _, m := range members {
			// TODO: add check/validation if given type is a Go type
			d := strings.Split(strings.Trim(m, " "), " ")
			member := capitalize(d[0])
			data[member] = d[1]
		}
	} else if strings.Contains(dataArg, " ") {
		d := strings.Split(dataArg, " ")
		member := capitalize(d[0])
		data[member] = d[1]
	}
	return data
}

// TODO: keep this function somewhere common
func capitalize(target string) string {
	r := []rune(target)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
