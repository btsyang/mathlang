package paser

import (
	"io"
	"strconv"
	"strings"
	// "strconv"
)

func ParseReader(r io.Reader) (*AST, error) {
	l := NewLexer(r)
	ast := &AST{
		Bases:      make(map[string]*Basis),
		Vecs:       make(map[string]*Vec),
		Transforms: make(map[string]*TransformRule),
	}
	var curBasis string
	for tok := l.Next(); tok != nil; tok = l.Next() {
		switch tok.Kind {
		case "BasisAssign":
			// fmt.Println("tok.Args", tok.Args)
			args := tok.Args.(*BasisAssignArgs)
			curBasis = args.Name
			if _, ok := ast.Bases[curBasis]; ok {
				panic("basis redefined: " + curBasis)
			}
			ast.Bases[curBasis] = &Basis{Name: curBasis}
			for _, vn := range args.Vecs {
				key := vn // e.g. "b1"
				if vec, ok := ast.Vecs[key]; !ok {
					panic("basis uses undefined vector: " + key)
				} else {
					ast.Bases[curBasis].Vecs = append(ast.Bases[curBasis].Vecs, vec)
				}
			}

		case "VectorAssign":
			args := tok.Args.(*VecAssignArgs)
			v := &Vec{Name: args.Name, Comp: args.Comp}
			ast.Vecs[args.Name] = v

		case "StmtTransformAssign":
			args := tok.Args.(*TransformAssignArgs)
			linearTerms := make([]LinearTerm, 0)
			toBasis := args.RawTerms[0][2]
			toBasisChecker := true
			domainVec := args.DomainVec[0] + args.DomainVec[1]

			for _, t := range args.RawTerms {
				if toBasisChecker {
					toBasisChecker = (toBasis == t[2])
				} else {
					panic("to Basis error ")
				}

				coeffStr := strings.ReplaceAll(t[1], " ", "")
				coeff := 1.0
				if coeffStr == "" || coeffStr == "+" {
					coeff = 1
				} else if coeffStr == "-" {
					coeff = -1
				} else {
					coeff, _ = strconv.ParseFloat(coeffStr, 64)
				}
				vec := t[2] + t[3]
				linearTerms = append(linearTerms, LinearTerm{Coeff: coeff, Vec: vec})
			}
			// 1. 查或建 TransformRule
			tr, ok := ast.Transforms[args.Transform]
			if !ok {
				tr = &TransformRule{
					Name:      args.Transform,
					FromBasis: ast.Bases[args.DomainVec[0]], // 默认
					ToBasis:   ast.Bases[toBasis],
					Map:       make(map[string][]LinearTerm),
				}
				ast.Transforms[args.Transform] = tr
			}

			// 2. 写入一行规则：T(b2) = ...
			tr.Map[domainVec] = linearTerms

		case "StmtEvalChangeBasis":
			args := tok.Args.(*EvalChangeBasisArgs)

			v, ok := ast.Vecs[args.Vec]
			if !ok {
				panic("eval uses undefined vector: " + args.Vec)
			}

			basis, ok := ast.Bases[args.Basis]
			if !ok {
				panic("eval uses undefined basis: " + args.Basis)
			}

			ast.Eval = &EvalChangeBasis{
				Vec:   v,
				Basis: basis,
			}
		case "StmtEvalTransform":
			args := tok.Args.(*EvalTransformArgs)
			v, ok := ast.Vecs[args.VecName]
			if !ok {
				panic("eval uses undefined vector: " + args.VecName)
			}

			t, ok := ast.Transforms[args.Transform]
			if !ok {
				panic("eval uses undefined transform: " + args.Transform)
			}

			ast.Eval = &EvalTransform{
				Transform: args.Transform,
				Rule:      t,
				Vec:       v,
			}
		}
	}
	return ast, nil
}
