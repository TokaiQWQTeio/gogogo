package evaluator

import (
	"gogogo/ast"
	"gogogo/object"
)

func quote(node ast.Node) object.Object {
	return &object.Quote{Node: node}
}
