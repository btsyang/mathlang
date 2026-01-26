package calculator

import (
	"fmt"

	"github.com/btsyang/mathlang/parser"
)

// printMatPretty 打印矩阵的漂亮格式
// 参数：
//
//	A: 要打印的矩阵
func printMatPretty(A [][]float64) {
	for i := range A {
		for j := range A[i] {
			fmt.Printf("%8.3f ", A[i][j])
		}
		fmt.Println()
	}
}

// evalChangeBasis 处理基变换计算
// 参数：
//
//	e: 基变换计算请求
//
// 返回：
//
//	[]float64: 计算结果，向量在新基下的坐标
//	error: 计算过程中遇到的错误
func evalChangeBasis(e *parser.EvalChangeBasis) ([]float64, error) {
	basis := e.Basis
	vec := e.Vec.Comp

	dim := len(basis.Vecs)
	if dim == 0 {
		return nil, fmt.Errorf("empty basis: %s", basis.Name)
	}

	// 拼旧基矩阵 B (dim x dim)
	B := make([][]float64, dim)
	for i := 0; i < dim; i++ {
		B[i] = make([]float64, dim)
		for j := 0; j < dim; j++ {
			B[i][j] = basis.Vecs[j].Comp[i]
		}
	}
	return solve(B, vec)

}

// solve 使用高斯消元法求解线性方程组 Bx = v
// 参数：
//
//	B: 系数矩阵
//	v: 右侧向量
//
// 返回：
//
//	[]float64: 解向量 x
//	error: 求解过程中遇到的错误，如矩阵奇异
func solve(B [][]float64, v []float64) ([]float64, error) {
	n := len(v)
	aug := make([][]float64, n)
	for i := 0; i < n; i++ {
		aug[i] = make([]float64, n+1)
		copy(aug[i][:n], B[i])
		aug[i][n] = v[i]
	}
	// 前向消元
	for i := 0; i < n; i++ {
		// 选择主元行
		maxRow := i
		for k := i + 1; k < n; k++ {
			if abs(aug[k][i]) > abs(aug[maxRow][i]) {
				maxRow = k
			}
		}
		// 交换行
		aug[i], aug[maxRow] = aug[maxRow], aug[i]

		// 检查主元是否为零（或接近零），如果是，则矩阵奇异
		if abs(aug[i][i]) < 1e-10 {
			return nil, fmt.Errorf("singular matrix: cannot solve linear system")
		}

		// 消元
		for k := i + 1; k < n; k++ {
			f := aug[k][i] / aug[i][i]
			for j := i; j < n+1; j++ {
				aug[k][j] -= f * aug[i][j]
			}
		}
	}
	// 回代求解
	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		// 检查主元是否为零（或接近零）
		if abs(aug[i][i]) < 1e-10 {
			return nil, fmt.Errorf("singular matrix: cannot solve linear system")
		}
		x[i] = aug[i][n] / aug[i][i]
		for k := i - 1; k >= 0; k-- {
			aug[k][n] -= aug[k][i] * x[i]
		}
	}
	return x, nil
}

// abs 计算浮点数的绝对值
// 参数：
//
//	x: 输入浮点数
//
// 返回：
//
//	float64: x 的绝对值
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// SolveTransform 处理线性变换计算
// 参数：
//
//	eval: 线性变换计算请求
//	tr: 线性变换规则
//
// 返回：
//
//	[]float64: 计算结果，向量经过线性变换后的坐标
//	error: 计算过程中遇到的错误
func SolveTransform(eval *parser.EvalTransform, tr *parser.TransformRule) ([]float64, error) {

	dim := len(tr.Map)
	result := make([]float64, dim)
	for i, bv := range tr.FromBasis.Vecs {
		// 检查映射是否存在
		if _, ok := tr.Map[bv.Name]; !ok {
			return nil, fmt.Errorf("no mapping for input vector: %s", bv.Name)
		}
		vi := eval.Vec.Comp[i]
		for _, term := range tr.Map[bv.Name] {
			j := tr.ToBasis.IndexOf(term.Vec)
			result[j] += vi * term.Coeff
		}
	}
	return result, nil
}
