package parser

import (
	"go/ast"
)

type Parser interface {
	Parse(input string) (ast.Node, error)
}
