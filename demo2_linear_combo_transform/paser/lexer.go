package paser

import (
	"bufio"
	// "fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Token struct {
	Kind string
	Args any
}

type VecAssignArgs struct {
	Name string
	Comp []float64
}

type BasisAssignArgs struct {
	Name string
	Vecs []string
}

type EvalChangeBasisArgs struct {
	Vec   string
	Basis string
}

// func (*EvalChangeBasis) evalKind() {} //这个是AST的东西

// type LinearTermArg struct {
// 	LinearTermLit map[string][]LinearTerm
// }

type TransformAssignArgs struct {
	Transform string   // "T"
	DomainVec []string // "b2"
	RawTerms  [][]string
}

type EvalTransformArgs struct {
	Transform string // "T"
	VecName   string
	// ToBasis   string 如果在规则部分定义 输入基和输出基， 这里先删掉，表达式里目前没有这个部分
}

// func (*EvalTransform) evalKind() {}

type Lexer struct {
	scanner *bufio.Scanner
}

type StmtKind int

const (
	StmtUnknown StmtKind = iota
	StmtVecAssign
	StmtBasisAssign
	StmtTransformAssign
	StmtEvalChangeBasis
	StmtEvalTransform
)

func classify(line string) StmtKind {
	switch {
	case strings.Contains(line, "pmatrix"):
		return StmtVecAssign
	case strings.Contains(line, "eval") && strings.Contains(line, "[\\vec"):
		return StmtEvalChangeBasis
	case strings.Contains(line, "eval") && strings.Contains(line, "T(\\vec"):
		return StmtEvalTransform
	case strings.Contains(line, "T(\\vec"):
		return StmtTransformAssign
	case strings.Contains(line, "{"):
		return StmtBasisAssign
	default:
		return StmtUnknown
	}
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{scanner: bufio.NewScanner(r)}
}

func (l *Lexer) Next() *Token {
	for l.scanner.Scan() {
		line := strings.TrimSpace(l.scanner.Text())

		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "*") {
			continue
		}

		switch classify(line) {
		case StmtVecAssign:
			re := regexp.MustCompile(
				`\\vec\{([a-zA-Z]+)\}(?:_([0-9]+))?\s*=\s*\\begin\{pmatrix\}\s*([+-]?\d+)\s*\\\\\s*([+-]?\d+)\s*\\end\{pmatrix\}`,
			)
			m := re.FindStringSubmatch(line)

			if m == nil {
				panic("invalid vector assignment: " + line)
			}
			name := ""
			if len(m[2]) == 0 {
				name = m[1]
			} else {
				name = m[1] + m[2]
			}
			comp := make([]float64, len(m)-3)
			for i, s := range m[3:] {
				comp[i], _ = strconv.ParseFloat(s, 64)
			}
			return &Token{Kind: "VectorAssign", Args: &VecAssignArgs{Name: name, Comp: comp}}

		case StmtEvalChangeBasis:
			re := regexp.MustCompile(
				`^\[\s*\\vec\{([a-zA-Z]+)\}\s*\]\s*_\s*([a-zA-Z]+)\s*\\leftarrow\s*\\text\{eval\}\s*$`,
			)
			m := re.FindStringSubmatch(line)
			return &Token{Kind: "StmtEvalChangeBasis", Args: &EvalChangeBasisArgs{Vec: m[1], Basis: m[2]}}

		case StmtTransformAssign:
			termRe := regexp.MustCompile(
				`([+-]?\s*\d*\.?\d*)\\vec\{([a-zA-Z][a-zA-Z0-9]*)\}(?:_([0-9]+))?`,
			)
			terms := termRe.FindAllStringSubmatch(line, -1)

			return &Token{
				Kind: "StmtTransformAssign",
				Args: &TransformAssignArgs{Transform: "T", DomainVec: terms[0][2:], RawTerms: terms[1:]},
			}

		case StmtBasisAssign:
			re := regexp.MustCompile(
				`^([a-zA-Z]+)\s*=\s*\\\{\s*(.+)\s*\\\}$`,
			)
			r := make([]string, 0, 3)
			m := re.FindStringSubmatch(line)
			if m == nil {
				panic("invalid basis assignment: " + line)
			}
			r = append(r, m[1])
			items := strings.Split(m[2], ",")
			re = regexp.MustCompile(
				`\\vec\{([a-zA-Z]+)\}(?:_([0-9]+))`,
			)

			for i, _ := range items {
				n := re.FindStringSubmatch(items[i])
				r = append(r, n[1]+n[2])
			}
			return &Token{Kind: "BasisAssign", Args: &BasisAssignArgs{Name: r[0], Vecs: r[1:]}}

		case StmtEvalTransform:
			// 正则匹配 T(\vec{v}) \leftarrow eval
			re := regexp.MustCompile(`^([A-Z])\(\s*\\vec\{([a-zA-Z][a-zA-Z0-9]*)\}(?:_([0-9]+))?\s*\)\s*\\leftarrow\s*\\text\{eval\}\s*$`)
			m := re.FindStringSubmatch(line)
			if m == nil {
				panic("invalid EvalTransform statement")
			}

			args := &EvalTransformArgs{
				Transform: m[1], //"T"
				VecName:   m[2],
			}

			return &Token{Kind: "StmtEvalTransform", Args: args}
		}

	}
	return nil
}
