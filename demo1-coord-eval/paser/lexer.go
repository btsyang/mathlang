package paser

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

type Token struct {
	Kind string
	Args []string
}

type Lexer struct {
	scanner *bufio.Scanner
}

type StmtKind int

const (
	StmtUnknown StmtKind = iota
	StmtVecAssign
	StmtBasisAssign
	StmtEval
)

func classify(line string) StmtKind {
	switch {
	case strings.Contains(line, "pmatrix"):
		return StmtVecAssign
	case strings.Contains(line, "eval"):
		return StmtEval
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

		if line == "" || strings.HasPrefix(line, "%") {
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
			return &Token{Kind: "VectorAssign", Args: m[1:]}

		case StmtEval:
			re := regexp.MustCompile(
				`^\[\s*\\vec\{([a-zA-Z]+)\}\s*\]\s*_\s*([a-zA-Z]+)\s*\\leftarrow\s*\\text\{eval\}\s*$`,
			)
			m := re.FindStringSubmatch(line)
			return &Token{Kind: "StmtEval", Args: m[1:]}

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

			return &Token{Kind: "BasisAssign", Args: r}
		}
	}
	return nil
}
