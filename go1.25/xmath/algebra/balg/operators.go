package balg

func And(exprs ...Expression) AndExpression {
	return &andExpression{
		args: exprs,
	}
}

func Or(exprs ...Expression) OrExpression {
	return &orExpression{
		args: exprs,
	}
}

func Var(name string) VariableExpression {
	return &variableExpression{
		name: name,
	}
}
