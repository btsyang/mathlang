package main

import (
	"fmt"
	"io"
	"log"
	"os"

	// "mathlang/calculator"
	"github.com/btsyang/mathlang/calculator"
	"github.com/btsyang/mathlang/parser"
)

// main 是程序的入口点，处理命令行参数，读取输入，解析表达式，执行计算并输出结果
func main() {
	var file io.Reader
	if len(os.Args) < 2 {
		// 没有提供文件参数，从标准输入读取
		fmt.Println("从标准输入读取输入...")
		file = os.Stdin
	} else {
		// 从文件读取
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		stat, _ := f.Stat()
		log.Printf("file size: %d bytes", stat.Size())
		f.Seek(0, 0)
		file = f
	}
	// =========================================

	// 喂给 parser
	ast, err := parser.ParseReader(file)
	if err != nil {
		log.Fatal(err)
	}

	// 3. 计算
	if ast.Eval != nil {
		// res := calculator.Calculate(ast)
		res, err := calculator.Calculate(ast)
		if err != nil {
			log.Fatal(err)
		}
		// 根据计算类型显示不同的输出格式
		switch e := ast.Eval.(type) {
		case *parser.EvalChangeBasis:
			fmt.Printf("[\\vec{%s}]_%s = (", e.Vec.Name, e.Basis.Name)
		case *parser.EvalTransform:
			fmt.Printf("%s(\\vec{%s}) = (", e.Transform, e.Vec.Name)
		}
		for i, x := range res {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(x)
		}
		fmt.Println(")")
	}
}
