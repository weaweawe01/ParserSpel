package main

import (
	"fmt"
	"github.com/weaweawe01/ParserSpel/ast"
)

func main() {
	parser := ast.NewSpelExpressionParser()
	// 测试普通表达式
	normalExpressions := []string{
		"T(java.lang.Runtime).getRuntime().exec('id')",
	}
	for i, expr := range normalExpressions {
		fmt.Printf("\n=== 表达式测试 %d ===\n", i+1)
		fmt.Printf("表达式: %s\n", expr)
		result, err := parser.ParseExpressionWithContext(expr, nil)
		if err != nil {
			fmt.Printf("❌ 解析错误: %v\n", err)
		} else {
			fmt.Printf("   ✅ 解析成功!\n")
			ast.PrintASTWithTitle(result.AST, "完整 AST 树形结构")
		}
	}
}
