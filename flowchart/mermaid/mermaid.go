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
	fset       *token.FileSet
	sourceTree *flowchart.SourceTree
	tmpl       *template.Template
	// Keeps track of the id to assign to the current node
	curId int
}

var _ flowchart.FlowChart = &MermaidFlowChart{}

func New() *MermaidFlowChart {
	// Initialise the template
	tmpl, err := template.New("flowchart").Parse(templateStr)
	if err != nil {
		panic(fmt.Errorf("failed to initalise mermaid flow chart. %w", err))
	}

	return &MermaidFlowChart{
		fset:       token.NewFileSet(),
		sourceTree: &flowchart.SourceTree{},
		tmpl:       tmpl,
	}
}

func (m *MermaidFlowChart) convertAstNodeToString(stmt ast.Node) string {
	buf := new(bytes.Buffer)
	printer.Fprint(buf, m.fset, stmt)
	return buf.String()
}

func (m *MermaidFlowChart) generateNameForStmt(stmt ast.Node) string {
	m.curId++
	return "id:" + strconv.Itoa(m.curId)
}

// fromIfStmt generates the flow-chart nodes by taking an if statement block as input.
func (m *MermaidFlowChart) fromIfStmt(ifStmt *ast.IfStmt) (startNode *flowchart.Node, prevNodes []*flowchart.Node) {
	startNode = &flowchart.Node{
		Name:  m.generateNameForStmt(ifStmt),
		Type:  flowchart.NodeType_Conditional,
		Stmts: []string{m.convertAstNodeToString(ifStmt.Cond)},
	}
	m.sourceTree.Nodes = append(m.sourceTree.Nodes, startNode)

	subStartNode, subPrevNodes := m.fromBlockStmt(ifStmt.Body)
	m.sourceTree.Edges = append(m.sourceTree.Edges, &flowchart.Edge{
		Start: startNode,
		End:   subStartNode,
		Label: "Yes",
	})
	prevNodes = append(prevNodes, subPrevNodes...)

	// If this if statement is the only code block & no else/else-if block exists, return.
	if ifStmt.Else == nil {
		prevNodes = append(prevNodes, startNode)
		return
	}

	switch elseStmt := ifStmt.Else.(type) {
	case *ast.BlockStmt:
		elseStartNode, elsePrevNodes := m.fromBlockStmt(elseStmt)
		prevNodes = append(prevNodes, elsePrevNodes...)
		m.sourceTree.Edges = append(m.sourceTree.Edges, &flowchart.Edge{
			Start: startNode,
			End:   elseStartNode,
			Label: "No",
		})
	case *ast.IfStmt:
		elseIfStartNode, elseIfPrevNodes := m.fromIfStmt(elseStmt)
		prevNodes = append(prevNodes, elseIfPrevNodes...)
		m.sourceTree.Edges = append(m.sourceTree.Edges, &flowchart.Edge{
			Start: startNode,
			End:   elseIfStartNode,
			Label: "No",
		})
	}
	return
}

func (m *MermaidFlowChart) fromForStmt(forStmt *ast.ForStmt) (startNode *flowchart.Node, prevNodes []*flowchart.Node) {
	// Add a node with the initialisation step if it exists.
	if forStmt.Init != nil {
		startNode = &flowchart.Node{
			Name:  m.generateNameForStmt(forStmt.Init),
			Type:  flowchart.NodeType_Process,
			Stmts: []string{strings.Replace(m.convertAstNodeToString(forStmt.Init), "\"", "", -1)},
		}
		m.sourceTree.Nodes = append(m.sourceTree.Nodes, startNode)
	}

	// Add a node for the loop condition.
	forDecisionNode := &flowchart.Node{
		Name:  m.generateNameForStmt(forStmt.Cond),
		Type:  flowchart.NodeType_Conditional,
		Stmts: []string{strings.Replace(m.convertAstNodeToString(forStmt.Cond), "\"", "", -1)},
	}
	m.sourceTree.Nodes = append(m.sourceTree.Nodes, forDecisionNode)
	if startNode == nil {
		// Since, startNode is not set in the previous step,
		// so the forDecisionNode should be the startNode.
		startNode = forDecisionNode
	} else {
		m.sourceTree.Edges = append(m.sourceTree.Edges, &flowchart.Edge{
			Start: startNode,
			End:   forDecisionNode,
		})
	}
	prevNodes = []*flowchart.Node{forDecisionNode}

	forBodyStartNode, forBodyEndNodes := m.fromBlockStmt(forStmt.Body)
	if forBodyStartNode != nil {
		m.sourceTree.Edges = append(m.sourceTree.Edges, &flowchart.Edge{
			Start: forDecisionNode,
			End:   forBodyStartNode,
			Label: "Yes",
		})
	}

	// This keeps track of the node to which all the ending nodes inside the for loop connect to.
	lastNodeInForLoop := forDecisionNode

	// Add a node with post condition if it exists.
	if forStmt.Post != nil {
		postNode := &flowchart.Node{
			Name:  m.generateNameForStmt(forStmt.Post),
			Type:  flowchart.NodeType_Process,
			Stmts: []string{strings.Replace(m.convertAstNodeToString(forStmt.Post), "\"", "", -1)},
		}
		m.sourceTree.Nodes = append(m.sourceTree.Nodes, postNode)

		// Since, the post step exists in the for loop, all the previous nodes should connect to this.
		lastNodeInForLoop = postNode

		m.sourceTree.Edges = append(m.sourceTree.Edges, &flowchart.Edge{
			Start: postNode,
			End:   forDecisionNode,
		})
	}
	for _, node := range forBodyEndNodes {
		m.sourceTree.Edges = append(m.sourceTree.Edges, &flowchart.Edge{
			Start: node,
			End:   lastNodeInForLoop,
		})
	}
	return
}

// fromBlockStmt sets the flow chart data structure by taking a block of
// code as input and converting each of the statements into a node in the flow chart.
func (m *MermaidFlowChart) fromBlockStmt(blockStmt *ast.BlockStmt) (startNode *flowchart.Node, prevNodes []*flowchart.Node) {
	for _, stmt := range blockStmt.List {
		// curNode keeps track of the flowchart node
		// generated from the current statement being processed.
		var curNode *flowchart.Node

		// nextPrev keeps track of all those nodes that will become
		// prev node in the next iteration.
		var nextPrev []*flowchart.Node

		switch stmt := stmt.(type) {
		case *ast.IfStmt:
			ifBlockStartNode, prevNodesFromIf := m.fromIfStmt(stmt)
			curNode = ifBlockStartNode
			nextPrev = prevNodesFromIf
		case *ast.ForStmt:
			forLoopStartNode, prevNodesFromFor := m.fromForStmt(stmt)
			curNode = forLoopStartNode
			nextPrev = prevNodesFromFor
		default:
			curNode = &flowchart.Node{
				Name:  m.generateNameForStmt(stmt),
				Type:  flowchart.NodeType_Process,
				Stmts: []string{strings.Replace(m.convertAstNodeToString(stmt), "\"", "", -1)},
			}
			m.sourceTree.Nodes = append(m.sourceTree.Nodes, curNode)
			nextPrev = []*flowchart.Node{curNode}
		}

		if startNode == nil {
			startNode = curNode
		}
		// TODO: Move this to its own function and handle for appropriate labels from if statements.
		for _, node := range prevNodes {
			m.sourceTree.Edges = append(m.sourceTree.Edges, &flowchart.Edge{
				Start: node,
				End:   curNode,
			})
		}
		prevNodes = nextPrev
	}
	return
}

func (m *MermaidFlowChart) FromAst(astNode ast.Node) {
	startNode := &flowchart.Node{
		Name:  "start",
		Type:  flowchart.NodeType_Terminal,
		Stmts: []string{"START"},
	}
	endNode := &flowchart.Node{
		Name:  "stop",
		Type:  flowchart.NodeType_Terminal,
		Stmts: []string{"STOP"},
	}

	m.sourceTree.Nodes = append(m.sourceTree.Nodes, startNode, endNode)

	// Iterate on the AST and find process the first function declaration.
	ast.Inspect(astNode, func(n ast.Node) bool {
		funcType, ok := n.(*ast.FuncDecl)
		if ok {
			m.sourceTree.Title = funcType.Name.String()

			subStartNode, prevNodes := m.fromBlockStmt(funcType.Body)
			m.sourceTree.Edges = append(m.sourceTree.Edges, &flowchart.Edge{
				Start: startNode,
				End:   subStartNode,
			})

			for _, node := range prevNodes {
				m.sourceTree.Edges = append(m.sourceTree.Edges, &flowchart.Edge{
					Start: node,
					End:   endNode,
				})
			}

			return false
		}
		return true
	})
}

func (m *MermaidFlowChart) String() (string, error) {
	buf := new(bytes.Buffer)
	err := m.tmpl.Execute(buf, m.sourceTree)
	return buf.String(), err
}

func (m *MermaidFlowChart) Clear() {
	m.sourceTree = &flowchart.SourceTree{}
}

const templateStr = `---
Title: {{ .Title }}
---
flowchart TD

{{ range .Nodes }}

{{- if eq .Type 1 }}
	{{ .Name }}(["
		{{ range .Stmts }}
		{{- . -}}
		{{ end }}
	"])
{{ else if eq .Type 2 }}
	{{ .Name }}[/"
		{{ range .Stmts }}
		{{- . -}}
		{{ end }}
	"/]
{{ else if eq .Type 3 }}
	{{ .Name }}["
		{{ range .Stmts }}
		{{- . -}}
		{{ end }}
	"]
{{ else if eq .Type 4 }}
	{{ .Name }}{"
		{{ range .Stmts }}
		{{- . -}}
		{{ end }}
	"}
{{ end -}}

{{ end }}

{{- range .Edges -}}
	{{- if .Label }}
	{{ .Start.Name }} --{{ .Label }}--> {{ .End.Name }}
	{{ else }}
	{{ .Start.Name }} --> {{ .End.Name }}
	{{ end -}}
{{ end -}}
`
