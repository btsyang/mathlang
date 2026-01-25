package main

import (
	"fmt"
	"log"
	"os"

	// "mathlang/calculator"
	"github.com/btsyang/mathlang/calculator"
	"github.com/btsyang/mathlang/paser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "用法: %s <文件.org>\n", os.Args[0])
		os.Exit(1)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	stat, _ := file.Stat()
	log.Printf("file size: %d bytes", stat.Size())
	file.Seek(0, 0)
	// =========================================

	// 喂给 phaser
	ast, err := paser.ParseReader(file)
	if err != nil {
		log.Fatal(err)
	}

	// 3. 计算
	if ast.Eval != nil {
		// res := calculator.Calculate(ast)
		res := calculator.Calculate(ast)
		fmt.Printf("[\\vec{v}]_b = (")
		for i, x := range res {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(x)
		}
		fmt.Println(")")
	}
}
