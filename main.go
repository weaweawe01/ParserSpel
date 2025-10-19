package main

import (
	"fmt"

	"github.com/weaweawe01/ParserSpel/ast"
)

func main() {
	parser := ast.NewSpelExpressionParser()
	// 测试普通表达式
	normalExpressions := []string{
		"{1,2,3,4,5,6}?.?[#this between {2, 4}]",
	}
	for i, expr := range normalExpressions {
		tokenizer := ast.NewTokenizer(expr)
		tokens, err := tokenizer.Process()
		if err != nil {
			fmt.Errorf("tokenization failed: %v", err)
			return
		}
		fmt.Println("词法序列Token:")
		for count, token := range tokens {
			fmt.Printf("[%d] %s\n", count, token)
		}
		fmt.Printf("\n=== 词法解析: %d ===", i+1)
		result, err := parser.ParseExpressionWithContext(expr, nil)
		if err != nil {
			fmt.Printf("❌ 解析错误: %v\n", err)
		} else {
			ast.PrintASTWithTitle(result.AST, "完整 AST 树形结构")
		}
	}
}
