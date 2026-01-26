package paser

import (
	// "fmt"
	"io"
	"strconv"
)

func ParseReader(r io.Reader) (*AST, error) {
	l := NewLexer(r)
	ast := &AST{
		Bases: make(map[string]*Basis),
		Vecs:  make(map[string]*Vec),
	}
	var curBasis string
	for tok := l.Next(); tok != nil; tok = l.Next() {
		switch tok.Kind {
		case "BasisAssign":
			// fmt.Println("tok.Args", tok.Args)
			curBasis = tok.Args[0]
			if _, ok := ast.Bases[curBasis]; ok {
				panic("basis redefined: " + curBasis)
			}
			ast.Bases[curBasis] = &Basis{Name: curBasis}
			for _, vn := range tok.Args[1:] {
				key := vn // e.g. "b1"
				if vec, ok := ast.Vecs[key]; !ok {
					panic("basis uses undefined vector: " + key)
				} else {
					ast.Bases[curBasis].Vecs = append(ast.Bases[curBasis].Vecs, vec)
				}
			}

		case "VectorAssign":
			name := ""
			if len(tok.Args[1]) == 0 {
				name = tok.Args[0]
			} else {
				name = tok.Args[0] + tok.Args[1]
			}
			comp := make([]float64, len(tok.Args)-2)
			for i, s := range tok.Args[2:] {
				comp[i], _ = strconv.ParseFloat(s, 64)
			}
			v := &Vec{Name: name, Comp: comp}
			ast.Vecs[name] = v

		case "StmtEval":
			if _, ok := ast.Vecs[tok.Args[0]]; !ok {
				panic("eval uses undefined vector: " + tok.Args[0])
			}
			ast.Eval = &EvalRequest{
				VecName:   tok.Args[0],
				BasisName: tok.Args[1],
			}
		}
	}
	return ast, nil
}
