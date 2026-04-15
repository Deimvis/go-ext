package balg

import (
	"fmt"
	"strings"

	"github.com/Deimvis/go-ext/go1.25/ext"
)

type Expression interface {
	Type() ExpressionType
	Variables() []VariableExpression
	Eval(vars map[string]bool) bool
}

type LiteralExpression interface {
	Expression
	Value() bool
}

type VariableExpression interface {
	Expression
	Name() string
}

type AndExpression interface {
	Expression
	Args() []Expression
}

type OrExpression interface {
	Expression
	Args() []Expression
}

var (
	L_TRUE  LiteralExpression = &literalExpression{v: true}
	L_FALSE LiteralExpression = &literalExpression{v: false}
)

type ExpressionType int

var (
	ET_Literal  ExpressionType = 0
	ET_Variable ExpressionType = 1
	ET_AndOp    ExpressionType = 2
	ET_OrOp     ExpressionType = 3
	// ET_NotOp = 4
)

type literalExpression struct {
	v bool
}

func (l *literalExpression) Type() ExpressionType {
	return ET_Literal
}

func (l *literalExpression) Variables() []VariableExpression {
	return nil
}

func (l *literalExpression) Eval(vars map[string]bool) bool {
	return l.v
}

func (e *literalExpression) Value() bool {
	return e.v
}

type variableExpression struct {
	name string
}

func (e *variableExpression) Type() ExpressionType {
	return ET_Variable
}

func (e *variableExpression) Variables() []VariableExpression {
	return []VariableExpression{e}
}

func (e *variableExpression) Eval(vars map[string]bool) bool {
	return vars[e.name]
}

func (e *variableExpression) Name() string {
	return e.name
}

type andExpression struct {
	args []Expression
}

func (e *andExpression) Type() ExpressionType {
	return ET_AndOp
}

func (e *andExpression) Variables() []VariableExpression {
	seen := make(map[string]struct{})
	var vars []VariableExpression
	for _, arg := range e.args {
		for _, v := range arg.Variables() {
			if _, alreadySeen := seen[v.Name()]; alreadySeen {
				continue
			}
			seen[v.Name()] = struct{}{}
			vars = append(vars, v)
		}
	}
	return vars
}

func (e *andExpression) Eval(vars map[string]bool) bool {
	for _, arg := range e.args {
		if !arg.Eval(vars) {
			return false
		}
	}
	return true
}

func (e *andExpression) Args() []Expression {
	return ext.CopyElements(e.args)
}

type orExpression struct {
	args []Expression
}

func (e *orExpression) Type() ExpressionType {
	return ET_OrOp
}

func (e *orExpression) Variables() []VariableExpression {
	seen := make(map[string]struct{})
	var vars []VariableExpression
	for _, arg := range e.args {
		for _, v := range arg.Variables() {
			if _, alreadySeen := seen[v.Name()]; alreadySeen {
				continue
			}
			seen[v.Name()] = struct{}{}
			vars = append(vars, v)
		}
	}
	return vars
}

func (e *orExpression) Eval(vars map[string]bool) bool {
	for _, arg := range e.args {
		if arg.Eval(vars) {
			return true
		}
	}
	return false
}

func (e *orExpression) Args() []Expression {
	return ext.CopyElements(e.args)
}

func ExpressionToString(e Expression) string {
	if e.Type() == ET_Literal {
		e := e.(LiteralExpression)
		if e.Value() {
			return "TRUE"
		} else {
			return "FALSE"
		}
	} else if e.Type() == ET_Variable {
		e := e.(VariableExpression)
		return e.Name()
	} else if e.Type() == ET_AndOp {
		e := e.(AndExpression)
		argsStr := ext.Map(e.Args(), ExpressionToString)
		return "(" + strings.Join(argsStr, " ∧ ") + ")"
	} else if e.Type() == ET_OrOp {
		e := e.(AndExpression)
		argsStr := ext.Map(e.Args(), ExpressionToString)
		return "(" + strings.Join(argsStr, " v ") + ")"
	}
	panic(fmt.Errorf("unsupported expression type: %d", e.Type()))
}
