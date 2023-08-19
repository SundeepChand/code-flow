package goparser

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type GoParser struct {
}

func New() *GoParser {
	return &GoParser{}
}

func (g *GoParser) Parse(input string) (ast.Node, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "demo.go", input, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return file, nil
}
