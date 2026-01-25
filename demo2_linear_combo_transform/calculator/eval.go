package calculator

import (
	"github.com/btsyang/mathlang/paser"
)

func Calculate(ast *paser.AST) []float64 {
	switch e := ast.Eval.(type) {
	case *paser.EvalChangeBasis:
		return evalChangeBasis(e)
	case *paser.EvalTransform:
		return SolveTransform(e, ast.Transforms[e.Transform])
	default:
		panic("unknown eval")
	}
}
