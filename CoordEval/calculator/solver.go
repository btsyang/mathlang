package calculator

import (
	"fmt"
	"github.com/btsyang/mathlang/paser"
)

func printMatPretty(A [][]float64) {
	for i := range A {
		for j := range A[i] {
			fmt.Printf("%8.3f ", A[i][j])
		}
		fmt.Println()
	}
}

// SolveEval 只处理 AST.Eval 请求
func SolveEval(ast *paser.AST) []float64 {
	//计算dim
	basis, ok := ast.Bases[ast.Eval.BasisName]
	if !ok {
		panic("undefined basis: " + ast.Eval.BasisName)
	}

	dim := len(basis.Vecs)
	if dim == 0 {
		panic("empty basis: " + basis.Name)
	}
	
	// 拼旧基矩阵 B (dim x dim)
	B := make([][]float64, dim)
	for i := 0; i < dim; i++ {
		B[i] = make([]float64, dim)
		for j := 0; j < dim; j++ {
			B[i][j] = ast.Bases[ast.Eval.BasisName].Vecs[j].Comp[i]
		}
	}
	// 转换向量 v
	v := ast.Vecs[ast.Eval.VecName].Comp

	return solve(B, v)

}

// 极简高斯消元（同前）
func solve(B [][]float64, v []float64) []float64 {
	n := len(v)
	aug := make([][]float64, n)
	for i := 0; i < n; i++ {
		aug[i] = make([]float64, n+1)
		copy(aug[i][:n], B[i])
		aug[i][n] = v[i]
	}
	// 前向 + 回代（代码略，同上一回合）
	for i := 0; i < n; i++ {
		maxRow := i
		for k := i + 1; k < n; k++ {
			if abs(aug[k][i]) > abs(aug[maxRow][i]) {
				maxRow = k
			}
		}
		aug[i], aug[maxRow] = aug[maxRow], aug[i]
		for k := i + 1; k < n; k++ {
			f := aug[k][i] / aug[i][i]
			for j := i; j < n+1; j++ {
				aug[k][j] -= f * aug[i][j]
			}
		}
	}
	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		x[i] = aug[i][n] / aug[i][i]
		for k := i - 1; k >= 0; k-- {
			aug[k][n] -= aug[k][i] * x[i]
		}
	}
	return x
}
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
