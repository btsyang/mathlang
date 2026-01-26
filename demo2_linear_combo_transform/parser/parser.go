package parser

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	// "strconv"
)

// ParseReader 从输入流中解析线性代数表达式，构建抽象语法树
// 参数：
//
//	r: 输入流，通常是文件或标准输入
//
// 返回：
//
//	*AST: 构建的抽象语法树
//	error: 解析过程中遇到的错误
func ParseReader(r io.Reader) (*AST, error) {
	l := NewLexer(r)
	ast := &AST{
		Bases:      make(map[string]*Basis),
		Vecs:       make(map[string]*Vec),
		Transforms: make(map[string]*TransformRule),
	}
	var curBasis string
	for {
		tok, err := l.Next()
		if err != nil {
			return nil, err
		}
		if tok == nil {
			break
		}
		switch tok.Kind {
		case "BasisAssign":
			// fmt.Println("tok.Args", tok.Args)
			args := tok.Args.(*BasisAssignArgs)
			curBasis = args.Name

			// 检查基名称是否符合规范（单个小写字母）
			if !regexp.MustCompile(`^[a-z]$`).MatchString(curBasis) {
				return nil, fmt.Errorf("invalid basis name: %s, basis name should be a single lowercase letter", curBasis)
			}

			if _, ok := ast.Bases[curBasis]; ok {
				return nil, fmt.Errorf("basis redefined: %s", curBasis)
			}
			ast.Bases[curBasis] = &Basis{Name: curBasis}

			// 检查分量名称是否符合规范（基名称加上数字下标）
			for _, vn := range args.Vecs {
				key := vn // e.g. "b1"

				// 检查分量名称是否以基名称开头，后跟数字
				if !regexp.MustCompile(`^` + curBasis + `\d+$`).MatchString(key) {
					return nil, fmt.Errorf("invalid vector name in basis %s: %s, vector name should be %s followed by number", curBasis, key, curBasis)
				}

				if vec, ok := ast.Vecs[key]; !ok {
					return nil, fmt.Errorf("basis uses undefined vector: %s", key)
				} else {
					ast.Bases[curBasis].Vecs = append(ast.Bases[curBasis].Vecs, vec)
					// 设置向量的 Basis 字段为当前基
					vec.Basis = ast.Bases[curBasis]
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
					return nil, fmt.Errorf("to Basis error: inconsistent basis in linear combination")
				}

				coeffStr := strings.ReplaceAll(t[1], " ", "")
				coeff := 1.0
				if coeffStr == "" || coeffStr == "+" {
					coeff = 1
				} else if coeffStr == "-" {
					coeff = -1
				} else {
					var err error
					coeff, err = strconv.ParseFloat(coeffStr, 64)
					if err != nil {
						return nil, fmt.Errorf("invalid coefficient: %s", coeffStr)
					}
				}
				vec := t[2] + t[3]
				linearTerms = append(linearTerms, LinearTerm{Coeff: coeff, Vec: vec})
			}
			// 1. 查或建 TransformRule
			tr, ok := ast.Transforms[args.Transform]
			if !ok {
				fromBasis := args.DomainVec[0]
				if _, exists := ast.Bases[fromBasis]; !exists {
					return nil, fmt.Errorf("undefined basis: %s", fromBasis)
				}
				if _, exists := ast.Bases[toBasis]; !exists {
					return nil, fmt.Errorf("undefined basis: %s", toBasis)
				}
				tr = &TransformRule{
					Name:      args.Transform,
					FromBasis: ast.Bases[fromBasis], // 默认
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
				return nil, fmt.Errorf("eval uses undefined vector: %s", args.Vec)
			}

			basis, ok := ast.Bases[args.Basis]
			if !ok {
				return nil, fmt.Errorf("eval uses undefined basis: %s", args.Basis)
			}

			ast.Eval = &EvalChangeBasis{
				Vec:   v,
				Basis: basis,
			}
		case "StmtEvalTransform":
			args := tok.Args.(*EvalTransformArgs)
			v, ok := ast.Vecs[args.VecName]
			if !ok {
				return nil, fmt.Errorf("eval uses undefined vector: %s", args.VecName)
			}

			t, ok := ast.Transforms[args.Transform]
			if !ok {
				return nil, fmt.Errorf("eval uses undefined transform: %s", args.Transform)
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
