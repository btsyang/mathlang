package paser

type AST struct {
	Bases map[string]*Basis
	Vecs  map[string]*Vec
	Eval  *EvalRequest
}

type Basis struct {
	Name string
	Vecs []*Vec // 顺序即列
}

type Vec struct {
	Name string
	Comp []float64
}

type EvalRequest struct {
	VecName   string
	BasisName string
}
