/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/cucumber/gherkin-go/v11"
	"github.com/cucumber/messages-go/v10"
	"github.com/iancoleman/orderedmap"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/kubernetes-sigs/ingress-controller-conformance/test/files"
)

func main() {
	var (
		update          bool
		features        []string
		conformancePath string
	)

	flag.BoolVar(&update, "update", false, "update files in place in case of missing steps or method definitions")
	flag.StringVar(&conformancePath, "conformance-path", "test/conformance", "path to conformance test package location")

	flag.Parse()

	// 1. verify flags
	features = flag.CommandLine.Args()
	if len(features) == 0 {
		fmt.Println("Usage: codegen [-update=false] [-conformance-path=test/conformance] [features]")
		fmt.Println()
		fmt.Println("Example: codegen features/default_backend.feature")
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	// 2. parse template
	codeTmpl, err := template.New("template").Funcs(templateFuncs).Parse(goTemplate)
	if err != nil {
		log.Fatalf("Unexpected error parsing template: %v", err)
	}

	// 3. if features is a directory, iterate and search for files with extension .feature
	if len(features) == 1 && files.IsDir(features[0]) {
		root := filepath.Dir(features[0])
		features = []string{}

		err := filepath.Walk(root, visitDir(&features))
		if err != nil {
			log.Fatalf("Unexpected error reading directory %v: %v", root, err)
		}
	}

	// 4. iterate feature files
	for _, path := range features {
		err := processFeature(path, conformancePath, update, codeTmpl)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func processFeature(path, conformance string, update bool, template *template.Template) error {
	// 5. parse feature file
	featureSteps, err := parseFeature(path)
	if err != nil {
		return fmt.Errorf("parsing feature file: %w", err)
	}

	// 6. generate package name to use
	packageName := generatePackage(path)
	// 7. check if go source file exists
	goFile := filepath.Join(conformance, packageName, "feature.go")
	isGoFileOk := files.Exists(goFile)

	mapping := &Mapping{
		Package:      packageName,
		FeatureFile:  path,
		Features:     featureSteps,
		NewFunctions: featureSteps,
		GoFile:       goFile,
	}

	// 8. Extract functions from go source code
	if isGoFileOk {
		goFunctions, err := extractFuncs(goFile)
		if err != nil {
			return fmt.Errorf("extracting go functions: %w", err)
		}

		mapping.GoDefinitions = goFunctions
	}

	// 9. check if feature file is in sync with go code
	isInSync := false

	signatureChanges := []SignatureChange{}

	if isGoFileOk {
		inFeatures := sets.NewString()
		inGo := sets.NewString()

		for _, feature := range mapping.Features {
			inFeatures.Insert(feature.Name)
		}

		for _, gofunc := range mapping.GoDefinitions {
			inGo.Insert(gofunc.Name)
		}

		if newFunctions := inFeatures.Difference(inGo); newFunctions.Len() > 0 {
			log.Printf("Feature file %v contains %v new function/s", mapping.FeatureFile, newFunctions.Len())
			isInSync = false

			var funcs []Function
			for _, f := range newFunctions.List() {
				for _, feature := range mapping.Features {
					if feature.Name == f {
						funcs = append(funcs, feature)
						break
					}
				}
			}

			mapping.NewFunctions = funcs
		} else {
			mapping.NewFunctions = []Function{}
		}

	FeaturesLoop:
		for _, feature := range mapping.Features {
			for _, gofunc := range mapping.GoDefinitions {
				if feature.Name != gofunc.Name {
					continue
				}

				// We need to compare function arguments checking only
				// the number and type. Is not possible to rely in the name
				// in the go code.
				featKeys := feature.Args.Keys()
				goKeys := gofunc.Args.Keys()
				if len(featKeys) != len(goKeys) {
					signatureChanges = append(signatureChanges, SignatureChange{
						Function: gofunc.Name,
						Have:     argsFromMap(gofunc.Args, true),
						Want:     argsFromMap(feature.Args, true),
					})

					break FeaturesLoop
				}

				for index, k := range featKeys {
					fv, _ := feature.Args.Get(k)
					gv, _ := gofunc.Args.Get(goKeys[index])

					if !reflect.DeepEqual(fv, gv) {
						signatureChanges = append(signatureChanges, SignatureChange{
							Function: gofunc.Name,
							Have:     argsFromMap(gofunc.Args, true),
							Want:     argsFromMap(feature.Args, true),
						})

						break FeaturesLoop
					}
				}

			}
		}
	}

	// 10. check signatures are ok
	if len(signatureChanges) != 0 {
		var argBuf bytes.Buffer
		for _, sc := range signatureChanges {
			argBuf.WriteString(fmt.Sprintf(`
function %v
	have %v
	want %v
`, sc.Function, sc.Have, sc.Want))
		}

		return fmt.Errorf("source file %v has a different signature/s:\n %v", mapping.GoFile, argBuf.String())
	}

	// 11. if in sync, nothing to do
	if isInSync {
		return nil
	}

	// 12. New go feature file
	if !isGoFileOk {
		if !update {
			return fmt.Errorf("generated code is out of date (from %v feature, new go file %v)",
				mapping.FeatureFile, mapping.GoFile)
		}

		log.Printf("Generating new go file %v...", mapping.GoFile)
		buf := bytes.NewBuffer(make([]byte, 0))

		err := template.Execute(buf, mapping)
		if err != nil {
			return err
		}

		// 10. if update is set
		if update {
			isDirOk := files.IsDir(mapping.GoFile)
			if !isDirOk {
				err := os.MkdirAll(filepath.Dir(mapping.GoFile), 0755)
				if err != nil {
					return err
				}
			}

			err := ioutil.WriteFile(mapping.GoFile, buf.Bytes(), 0644)
			if err != nil {
				return err
			}

			featFile := filepath.Base(path)
			log.Printf(`Please add '"features/%v": %v.FeatureContext' to features map defined in conformance_test.go (order matters)`,
				featFile, mapping.Package)
		}

		return nil
	}

	if len(mapping.NewFunctions) == 0 {
		return nil
	}

	// 13. if update is set
	if update {
		log.Printf("Updating go file %v...", mapping.GoFile)
		err := updateGoTestFile(mapping.GoFile, mapping.NewFunctions)
		if err != nil {
			return err
		}
	}

	if len(mapping.NewFunctions) != 0 {
		return fmt.Errorf("generated code is out of date")
	}

	return nil
}

// Function holds the definition of a function in a go file or godog step
type Function struct {
	// Name
	Name string
	// Expr Regexp to use in godog Step definition
	Expr string
	// Args function arguments
	// k = name of the argument
	// v = type of the argument
	Args *orderedmap.OrderedMap
}

type Mapping struct {
	Package string

	FeatureFile string
	Features    []Function

	GoFile        string
	GoDefinitions []Function

	NewFunctions []Function
}

// SignatureChange holds information about the definition of a go function
type SignatureChange struct {
	Function string
	Have     string
	Want     string
}

var templateFuncs = template.FuncMap{
	"backticked": func(s string) string {
		return "`" + s + "`"
	},
	"unescape": func(s string) template.HTML {
		return template.HTML(s)
	},
	"argsFromMap": argsFromMap,
}

const goTemplate = `
/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package {{ .Package }}

import (
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"

	tstate "github.com/kubernetes-sigs/ingress-controller-conformance/test/state"
	"github.com/kubernetes-sigs/ingress-controller-conformance/test/kubernetes"
)

var (
	state *tstate.Scenario
)

{{ range .NewFunctions }}
func {{ .Name }}{{ argsFromMap .Args false }} error {
	return godog.ErrPending
}
{{end}}

func FeatureContext(s *godog.Suite) { {{ range .NewFunctions }}
	s.Step({{ backticked .Expr | unescape }}, {{ .Name }}){{end}}

	s.BeforeScenario(func(this *messages.Pickle) {
		state = tstate.New(nil)
	})

	s.AfterScenario(func(*messages.Pickle, error) {
		// delete namespace an all the content
		_ = kubernetes.DeleteNamespace(kubernetes.KubeClient, state.Namespace)
	})
}
`

// parseFeature parses a godog feature file returning the unique
// steps definitions
func parseFeature(path string) ([]Function, error) {
	data, err := files.Read(path)
	if err != nil {
		return nil, err
	}

	gd, err := gherkin.ParseGherkinDocument(bytes.NewReader(data), (&messages.Incrementing{}).NewId)
	if err != nil {
		return nil, err
	}

	scenarios := gherkin.Pickles(*gd, path, (&messages.Incrementing{}).NewId)

	def := []Function{}
	for _, s := range scenarios {
		def = parseSteps(s.Steps, def)
	}

	return def, nil
}

// extractFuncs reads a file containing go source code and returns
// the functions defined in the file.
func extractFuncs(filePath string) ([]Function, error) {
	if !strings.HasSuffix(filePath, ".go") {
		return nil, fmt.Errorf("only files with go extension are valid")
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	funcs := []Function{}

	var printErr error
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		index := 0
		args := orderedmap.New()
		for _, p := range fn.Type.Params.List {
			var typeNameBuf bytes.Buffer

			err := printer.Fprint(&typeNameBuf, fset, p.Type)
			if err != nil {
				printErr = err
				return false
			}

			if len(p.Names) == 0 {
				argName := fmt.Sprintf("arg%d", index+1)
				args.Set(argName, typeNameBuf.String())

				index++
				continue
			}

			for _, ag := range p.Names {
				argName := ag.String()
				args.Set(argName, typeNameBuf.String())
				index++
			}
		}

		// Go functions do not have an expression
		funcs = append(funcs, Function{Name: fn.Name.Name, Args: args})

		return true
	})

	if printErr != nil {
		return nil, printErr
	}

	return funcs, nil
}

func updateGoTestFile(filePath string, newFuncs []Function) error {
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var featureFunc *ast.FuncDecl
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		if fn.Name.Name == "FeatureContext" {
			featureFunc = fn
		}

		return true
	})

	if featureFunc == nil {
		return fmt.Errorf("file %v does not contains a FeatureFunct function", filePath)
	}

	// Add new functions
	astf, err := toAstFunctions(newFuncs)
	if err != nil {
		return err
	}

	node.Decls = append(node.Decls, astf...)

	// Update steps in FeatureContext
	astSteps, err := toContextStepsfuncs(newFuncs)
	if err != nil {
		return err
	}

	featureFunc.Body.List = append(astSteps, featureFunc.Body.List...)

	var buffer bytes.Buffer
	if err = format.Node(&buffer, fileSet, node); err != nil {
		return fmt.Errorf("error formatting file %v: %w", filePath, err)
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %v: %w", filePath, err)
	}

	return ioutil.WriteFile(filePath, buffer.Bytes(), fileInfo.Mode())
}

func toContextStepsfuncs(funcs []Function) ([]ast.Stmt, error) {
	astStepsTpl := `
package codegen
func FeatureContext() { {{ range . }}
	s.Step({{ backticked .Expr | unescape }}, {{ .Name }}){{end}}
}
`
	astFile, err := astFromTemplate(astStepsTpl, funcs)
	if err != nil {
		return nil, err
	}

	f := astFile.Decls[0].(*ast.FuncDecl)

	return f.Body.List, nil
}

func toAstFunctions(funcs []Function) ([]ast.Decl, error) {
	astFuncTpl := `
package codegen
{{ range . }}func {{ .Name }}{{ argsFromMap .Args false }} error {
	return godog.ErrPending
}

{{end}}
`
	astFile, err := astFromTemplate(astFuncTpl, funcs)
	if err != nil {
		return nil, err
	}

	return astFile.Decls, nil
}

func astFromTemplate(astFuncTpl string, funcs []Function) (*ast.File, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	astFuncs, err := template.New("ast").Funcs(templateFuncs).Parse(astFuncTpl)
	if err != nil {
		return nil, err
	}

	err = astFuncs.Execute(buf, funcs)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "src.go", buf.String(), parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return astFile, nil
}

// generatePackage returns the name of the
// package to use using the feature filename
func generatePackage(filePath string) string {
	base := path.Base(filePath)
	base = strings.ToLower(base)
	base = strings.ReplaceAll(base, "_", "")
	base = strings.ReplaceAll(base, ".feature", "")

	return base
}

func argsFromMap(args *orderedmap.OrderedMap, onlyType bool) string {
	s := "("

	for _, k := range args.Keys() {
		v, ok := args.Get(k)
		if !ok {
			continue
		}

		if onlyType {
			s += fmt.Sprintf("%v, ", v)
		} else {
			s += fmt.Sprintf("%v %v, ", k, v)
		}
	}

	if len(args.Keys()) > 0 {
		s = s[0 : len(s)-2]
	}

	return s + ")"
}

func visitDir(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if filepath.Ext(path) != ".feature" {
			return nil
		}

		*files = append(*files, path)
		return nil
	}
}

// Code below this comment comes from github.com/cucumber/godog
// (code defined in private methods)

const (
	numberGroup = "(\\d+)"
	stringGroup = "\"([^\"]*)\""
)

// parseStepArgs extracts arguments from an expression defined in a step RegExp.
// This code was extracted from
// https://github.com/cucumber/godog/blob/4da503aab2d0b71d380fbe8c48a6af9f729b6f5a/undefined_snippets_gen.go#L41
func parseStepArgs(exp string, argument *messages.PickleStepArgument) *orderedmap.OrderedMap {
	var (
		args      []string
		pos       int
		breakLoop bool
	)

	for !breakLoop {
		part := exp[pos:]
		ipos := strings.Index(part, numberGroup)
		spos := strings.Index(part, stringGroup)

		switch {
		case spos == -1 && ipos == -1:
			breakLoop = true
		case spos == -1:
			pos += ipos + len(numberGroup)
			args = append(args, "int")
		case ipos == -1:
			pos += spos + len(stringGroup)
			args = append(args, "string")
		case ipos < spos:
			pos += ipos + len(numberGroup)
			args = append(args, "int")
		case spos < ipos:
			pos += spos + len(stringGroup)
			args = append(args, "string")
		}
	}

	if argument != nil {
		if argument.GetDocString() != nil {
			args = append(args, "*messages.PickleStepArgument_PickleDocString")
		}

		if argument.GetDataTable() != nil {
			args = append(args, "*messages.PickleStepArgument_PickleTable")
		}
	}

	stepArgs := orderedmap.New()
	for i, v := range args {
		k := fmt.Sprintf("arg%d", i+1)
		stepArgs.Set(k, v)
	}

	return stepArgs
}

// some snippet formatting regexps
var snippetExprCleanup = regexp.MustCompile("([\\/\\[\\]\\(\\)\\\\^\\$\\.\\|\\?\\*\\+\\'])")
var snippetExprQuoted = regexp.MustCompile("(\\W|^)\"(?:[^\"]*)\"(\\W|$)")
var snippetMethodName = regexp.MustCompile("[^a-zA-Z\\_\\ ]")
var snippetNumbers = regexp.MustCompile("(\\d+)")

// parseSteps converts a string step definition in a different one valid as a regular
// expression that can be used in a go Step definition. This original code is located in
// https://github.com/cucumber/godog/blob/4da503aab2d0b71d380fbe8c48a6af9f729b6f5a/fmt.go#L457
func parseSteps(steps []*messages.Pickle_PickleStep, funcDefs []Function) []Function {
	var index int

	for _, step := range steps {
		text := step.Text

		expr := snippetExprCleanup.ReplaceAllString(text, "\\$1")
		expr = snippetNumbers.ReplaceAllString(expr, "(\\d+)")
		expr = snippetExprQuoted.ReplaceAllString(expr, "$1\"([^\"]*)\"$2")
		expr = "^" + strings.TrimSpace(expr) + "$"

		name := snippetNumbers.ReplaceAllString(text, " ")
		name = snippetExprQuoted.ReplaceAllString(name, " ")
		name = strings.TrimSpace(snippetMethodName.ReplaceAllString(name, ""))

		var words []string
		for i, w := range strings.Split(name, " ") {
			switch {
			case i != 0:
				w = strings.Title(w)
			case len(w) > 0:
				w = string(unicode.ToLower(rune(w[0]))) + w[1:]
			}

			words = append(words, w)
		}

		name = strings.Join(words, "")
		if len(name) == 0 {
			index++
			name = fmt.Sprintf("StepDefinitioninition%d", index)
		}

		var found bool
		for _, f := range funcDefs {
			if f.Expr == expr {
				found = true
				break
			}
		}

		if !found {
			args := parseStepArgs(expr, step.Argument)
			funcDefs = append(funcDefs, Function{Name: name, Expr: expr, Args: args})
		}
	}

	return funcDefs
}
