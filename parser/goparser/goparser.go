package goparser

import (
	"go/ast"
	"go/parser"
	"go/token"

	parserInterface "github.com/SundeepChand/code-flow/parser"
)

type GoParser struct {
}

var _ parserInterface.Parser = &GoParser{}

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
