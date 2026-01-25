package paser

type AST struct {
	Bases      map[string]*Basis
	Vecs       map[string]*Vec
	Transforms map[string]*TransformRule
	Eval       EvalStmt
}

type EvalStmt interface {
	evalKind()
}

type Basis struct {
	Name string
	Vecs []*Vec // 顺序即列
}

type BasisEnv map[string]*Basis

type Vec struct {
	Name string
	Comp []float64
}

type LinearTerm struct {
	Coeff float64 // 系数
	Vec   string  // 向量名（符号引用）
}

type TransformRule struct {
	Name      string // "T"
	FromBasis *Basis // "b"
	ToBasis   *Basis // "c"
	Map       map[string][]LinearTerm
}

type EvalChangeBasis struct {
	Vec   *Vec   // 已绑定的向量
	Basis *Basis // 已绑定的基
}

func (*EvalChangeBasis) evalKind() {}

type EvalTransform struct {
	Transform string         // "T"
	Rule      *TransformRule // 已绑定规则
	Vec       *Vec           // 输入向量（在 FromBasis 下）
}

func (*EvalTransform) evalKind() {}

func (b *Basis) IndexOf(vecName string) int {
	for i, v := range b.Vecs {
		if v.Name == vecName {
			return i
		}
	}
	panic("vector " + vecName + " not found in basis " + b.Name)
}
