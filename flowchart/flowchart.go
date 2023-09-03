package flowchart

import "go/ast"

type FlowChart interface {
	Clear()
	FromAst(astNode ast.Node)
	String() (string, error)
}

type NodeType uint8

const (
	NodeType_Unspecified NodeType = iota
	NodeType_Terminal
	NodeType_IO
	NodeType_Process
	NodeType_Conditional
)

type Node struct {
	Type  NodeType
	Name  string
	Stmts []string
}

type Edge struct {
	Start, End *Node
	Label      string
}

type SourceTree struct {
	Title string
	Nodes []*Node
	Edges []*Edge
}
