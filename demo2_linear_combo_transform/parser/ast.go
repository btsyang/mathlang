package parser

// AST 是抽象语法树的根节点，包含所有定义和计算请求
type AST struct {
	Bases      map[string]*Basis         // 基的映射，键为基名
	Vecs       map[string]*Vec           // 向量的映射，键为向量名
	Transforms map[string]*TransformRule // 线性变换规则的映射，键为变换名
	Eval       EvalStmt                  // 计算请求
}

// EvalStmt 是计算请求的接口，有两个实现：EvalChangeBasis 和 EvalTransform
type EvalStmt interface {
	evalKind() // 接口方法，用于类型断言
}

// Vec 表示一个向量，包含名称、基和分量
type Vec struct {
	Name  string    // 向量名称
	Basis *Basis    // 向量所属的基
	Comp  []float64 // 向量的分量
}

// Basis 表示一个基，包含名称和向量列表
type Basis struct {
	Name string // 基的名称
	Vecs []*Vec // 基中的向量列表，顺序即列序
}

// IndexOf 查找向量在基中的索引
func (b *Basis) IndexOf(vecName string) int {
	for i, v := range b.Vecs {
		if v.Name == vecName {
			return i
		}
	}
	panic("vector " + vecName + " not found in basis " + b.Name)
}

// BasisEnv 是基的环境映射
type BasisEnv map[string]*Basis

// LinearTerm 表示线性组合中的一项，包含系数和向量名
type LinearTerm struct {
	Coeff float64 // 系数
	Vec   string  // 向量名（符号引用）
}

// TransformRule 表示线性变换规则，包含名称、输入基、输出基和映射
type TransformRule struct {
	Name      string                  // 变换名称，如 "T"
	FromBasis *Basis                  // 输入基
	ToBasis   *Basis                  // 输出基
	Map       map[string][]LinearTerm // 映射，键为输入基中的向量名，值为输出基中的线性组合
}

// EvalChangeBasis 表示基变换计算请求
type EvalChangeBasis struct {
	Vec   *Vec   // 已绑定的向量
	Basis *Basis // 已绑定的基
}

func (*EvalChangeBasis) evalKind() {}

// EvalTransform 表示线性变换计算请求
type EvalTransform struct {
	Transform string         // 变换名称，如 "T"
	Rule      *TransformRule // 已绑定的变换规则
	Vec       *Vec           // 输入向量（在 FromBasis 下）
}

func (*EvalTransform) evalKind() {}
