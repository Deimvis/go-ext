package balg

import (
	"slices"

	"github.com/Deimvis/go-ext/go1.25/ext"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

// ToCNF convert a propositional logical sentence “expr“ to conjunctive normal form.
// This is a very limited implementation. The true implementation can be found here:
// https://github.com/sympy/sympy/blob/b4ce69ad5d40e4e545614b6c76ca9b0be0b98f0b/sympy/logic/boolalg.py#L1654
func ToCNF(e Expression) Expression {
	cnf := distributeANDoverOR(e)
	return cnf
}

func ToCNFCanonical(e Expression) AndExpression {
	cnfMin := ToCNF(e)

	var cnf AndExpression
	if cnfMin.Type() == ET_AndOp {
		cnf = cnfMin.(AndExpression)
	} else {
		// only one clause
		cnf = And(cnfMin)
	}
	cnf = And(ext.Map(cnf.Args(), func(arg Expression) Expression {
		if arg.Type() == ET_OrOp {
			return arg
		}
		// only one argument clause
		return Or(arg)
	})...)

	return cnf
}

// ToDNF convert a propositional logical sentence “expr“ to disjunctive normal form.
// This is a very limited implementation. The true implementation can be found here:
// https://github.com/sympy/sympy/blob/b4ce69ad5d40e4e545614b6c76ca9b0be0b98f0b/sympy/logic/boolalg.py#L1696
func ToDNF(e Expression) Expression {
	panic("not implemented")
}

func distributeANDoverOR(e Expression) Expression {
	if e.Type() == ET_AndOp {
		e := e.(AndExpression)
		return andSimplified(ext.Map(e.Args(), distributeANDoverOR)...)
	} else if e.Type() == ET_OrOp {
		e := e.(OrExpression)
		var andArg AndExpression
		var andArgInd int = -1
		for ind, arg := range e.Args() {
			if arg.Type() == ET_AndOp {
				arg := arg.(AndExpression)
				andArg = arg
				andArgInd = ind
				break
			}
		}
		if andArgInd == -1 {
			return e
		}
		var otherArgs []Expression
		for ind, arg := range e.Args() {
			if ind != andArgInd {
				otherArgs = append(otherArgs, arg)
			}
		}
		return andSimplified(
			ext.Map(
				ext.Map(
					andArg.Args(),
					func(andArgArg Expression) Expression {
						return orSimplified(append(ext.CopyElements(otherArgs), andArgArg)...)
					}),
				distributeANDoverOR)...,
		)
	} else {
		return e
	}
}

func andSimplified(args ...Expression) AndExpression {
	return And(unwrapRecursively(args, ET_AndOp)...)
}

func orSimplified(args ...Expression) AndExpression {
	return Or(unwrapRecursively(args, ET_OrOp)...)
}

func unwrapRecursively(args []Expression, t ExpressionType) []Expression {
	xmust.True((t == ET_AndOp) || (t == ET_OrOp))
	var finalArgs []Expression
	stack := args
	slices.Reverse(stack)
	for len(stack) > 0 {
		el := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if el.Type() == t {
			if t == ET_AndOp {
				el := el.(AndExpression)
				elArgs := el.Args()
				slices.Reverse(elArgs)
				stack = append(stack, elArgs...)
			} else {
				el := el.(OrExpression)
				elArgs := el.Args()
				slices.Reverse(elArgs)
				stack = append(stack, elArgs...)
			}
		} else {
			finalArgs = append(finalArgs, el)
		}
	}
	return finalArgs
}
