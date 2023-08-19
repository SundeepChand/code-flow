package mermaid

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strconv"
	"strings"
	"text/template"

	"github.com/SundeepChand/code-flow/flowchart"
)

type MermaidFlowChart struct {
	sourceTree *flowchart.SourceTree
}

func New() *MermaidFlowChart {
	return &MermaidFlowChart{
		sourceTree: &flowchart.SourceTree{},
	}
}

func (m *MermaidFlowChart) FromAst(astNode ast.Node) {
	fset := token.NewFileSet()

	ast.Inspect(astNode, func(n ast.Node) bool {
		funcType, ok := n.(*ast.FuncDecl)
		if ok {
			m.sourceTree.Title = funcType.Name.String()
			for idx, stmt := range funcType.Body.List {
				buf := new(bytes.Buffer)
				printer.Fprint(buf, fset, stmt)

				m.sourceTree.Nodes = append(m.sourceTree.Nodes, &flowchart.Node{
					Name:  "id" + strconv.Itoa(idx),
					Type:  flowchart.NodeType_Process,
					Stmts: []string{strings.Replace(buf.String(), "\"", "", -1)},
				})

				switch stmt.(type) {
				case *ast.IfStmt:
					fmt.Println("If stmt", stmt)
				case *ast.ForStmt:
					fmt.Println("For stmt", stmt)
				default:
					fmt.Println(stmt)
				}
			}

			return false
		}
		return true
	})
}

func (m *MermaidFlowChart) String() (string, error) {
	// TODO: Move this to constructor
	tmpl, err := template.New("flowchart").Parse(templateStr)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, m.sourceTree)
	return buf.String(), err
}

const templateStr = `---
Title: {{ .Title }}
---
flowchart TD

{{ range .Nodes }}

{{- if eq .Type 1 }}
	{{ .Name }}[/"
		{{ range .Stmts }}
		{{ . -}}
		{{ end }}
	"/]
{{ else if eq .Type 2 }}
	{{ .Name }}["
		{{ range .Stmts }}
		{{ . -}}
		{{ end }}
	"]
{{ else if eq .Type 3 }}
	{{ .Name }}{"
		{{ range .Stmts }}
		{{ . -}}
		{{ end }}
	"}
{{ else if eq .Type 4 }}
	{{ .Name }}(["
		{{ range .Stmts }}
		{{ . -}}
		{{ end }}
	"])
{{ end -}}

{{ end }}

{{- range .Edges -}}
	{{- if .Label }}
	{{ .Start.Name }} --{{ .Label }}--> {{ .End.Name }}
	{{ else }}
	{{ .Start.Name }} --> {{ .End.Name }}
	{{ end }}
{{ end -}}
`
