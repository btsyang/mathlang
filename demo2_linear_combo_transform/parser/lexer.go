package parser

import (
	"bufio"
	"fmt"

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
	scanner           *bufio.Scanner
	vecAssignRe       *regexp.Regexp
	basisAssignRe     *regexp.Regexp
	evalChangeBasisRe *regexp.Regexp
	transformAssignRe *regexp.Regexp
	evalTransformRe   *regexp.Regexp
	termRe            *regexp.Regexp
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

// NewLexer 创建一个新的词法分析器
// 参数：
//
//	r: 输入流，通常是文件或标准输入
//
// 返回：
//
//	*Lexer: 创建的词法分析器实例
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		scanner:           bufio.NewScanner(r),
		vecAssignRe:       regexp.MustCompile(`\\vec\{([a-zA-Z]+)\}(?:_([0-9]+))?\s*=\s*\\begin\{pmatrix\}\s*([+-]?\d+)\s*\\\\\s*([+-]?\d+)\s*\\end\{pmatrix\}`),
		basisAssignRe:     regexp.MustCompile(`^([a-zA-Z]+)\s*=\s*\\\{\s*(.+)\s*\\\}$`),
		evalChangeBasisRe: regexp.MustCompile(`^\[\s*\\vec\{([a-zA-Z]+)\}\s*\]\s*_\s*([a-zA-Z]+)\s*\\leftarrow\s*\\text\{eval\}\s*$`),
		transformAssignRe: regexp.MustCompile(`([+-]?\s*\d*\.?\d*)\\vec\{([a-zA-Z][a-zA-Z0-9]*)\}(?:_([0-9]+))?`),
		evalTransformRe:   regexp.MustCompile(`^([A-Z])\(\s*\\vec\{([a-zA-Z][a-zA-Z0-9]*)\}(?:_([0-9]+))?\s*\)\s*\\leftarrow\s*\\text\{eval\}\s*$`),
		termRe:            regexp.MustCompile(`\\vec\{([a-zA-Z]+)\}(?:_([0-9]+))`),
	}
}

// Next 从输入流中读取下一个 token
// 返回：
//
//	*Token: 读取的 token
//	error: 读取过程中遇到的错误
func (l *Lexer) Next() (*Token, error) {
	for l.scanner.Scan() {
		line := strings.TrimSpace(l.scanner.Text())

		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "*") {
			continue
		}

		switch classify(line) {
		case StmtVecAssign:
			m := l.vecAssignRe.FindStringSubmatch(line)

			if m == nil {
				return nil, fmt.Errorf("invalid vector assignment: %s", line)
			}
			name := ""
			if len(m[2]) == 0 {
				name = m[1]
			} else {
				name = m[1] + m[2]
			}
			comp := make([]float64, len(m)-3)
			for i, s := range m[3:] {
				var err error
				comp[i], err = strconv.ParseFloat(s, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid component value: %s", s)
				}
			}
			return &Token{Kind: "VectorAssign", Args: &VecAssignArgs{Name: name, Comp: comp}}, nil

		case StmtEvalChangeBasis:
			m := l.evalChangeBasisRe.FindStringSubmatch(line)
			if m == nil {
				return nil, fmt.Errorf("invalid basis change evaluation: %s", line)
			}
			return &Token{Kind: "StmtEvalChangeBasis", Args: &EvalChangeBasisArgs{Vec: m[1], Basis: m[2]}}, nil

		case StmtTransformAssign:
			terms := l.transformAssignRe.FindAllStringSubmatch(line, -1)
			if len(terms) < 2 {
				return nil, fmt.Errorf("invalid transform assignment: %s", line)
			}

			return &Token{
				Kind: "StmtTransformAssign",
				Args: &TransformAssignArgs{Transform: "T", DomainVec: terms[0][2:], RawTerms: terms[1:]},
			}, nil

		case StmtBasisAssign:
			m := l.basisAssignRe.FindStringSubmatch(line)
			if m == nil {
				return nil, fmt.Errorf("invalid basis assignment: %s", line)
			}
			r := make([]string, 0, 3)
			r = append(r, m[1])
			items := strings.Split(m[2], ",")

			for i, _ := range items {
				n := l.termRe.FindStringSubmatch(items[i])
				if n == nil {
					return nil, fmt.Errorf("invalid vector in basis: %s", items[i])
				}
				r = append(r, n[1]+n[2])
			}
			return &Token{Kind: "BasisAssign", Args: &BasisAssignArgs{Name: r[0], Vecs: r[1:]}}, nil

		case StmtEvalTransform:
			// 正则匹配 T(\vec{v}) \leftarrow eval
			m := l.evalTransformRe.FindStringSubmatch(line)
			if m == nil {
				return nil, fmt.Errorf("invalid transform evaluation: %s", line)
			}

			args := &EvalTransformArgs{
				Transform: m[1], //"T"
				VecName:   m[2],
			}

			return &Token{Kind: "StmtEvalTransform", Args: args}, nil
		}

	}
	return nil, nil
}
