package calculator

import (
	"fmt"

	"github.com/btsyang/mathlang/parser"
)

// Calculate 根据抽象语法树执行相应的计算
// 参数：
//
//	ast: 抽象语法树，包含所有定义和计算请求
//
// 返回：
//
//	[]float64: 计算结果，通常是向量的分量
//	error: 计算过程中遇到的错误
func Calculate(ast *parser.AST) ([]float64, error) {
	switch e := ast.Eval.(type) {
	case *parser.EvalChangeBasis:
		return evalChangeBasis(e)
	case *parser.EvalTransform:
		return SolveTransform(e, ast.Transforms[e.Transform])
	default:
		return nil, fmt.Errorf("unknown eval type")
	}
}
