package main

import (
	"fmt"
	"github.com/weaweawe01/ParserSpel/ast"
)

// AST 分析工具 - 专门用于解析和显示 SpEL 表达式的 AST 结构
func runASTAnalyzer() {
	fmt.Println("\nSpEL AST 分析工具")
	fmt.Println("================")

	// 测试表达式
	testExpressions := []string{
		"42",
		"2 + 3",
		"2 + 3 * 4",
		"(2 + 3) * 4",
		"name == 'John'",
		"age > 18 && active == true",
		"user.name",
		"obj?.property",
		"#root",
		"@myBean",
		"name matches '[A-Z].*'",
	}

	parser := ast.NewSpelExpressionParser()

	for i, expr := range testExpressions {
		fmt.Printf("\n%d. 表达式: %s\n", i+1, expr)
		analyzeExpression(parser, expr)
	}
}

func analyzeExpression(parser *ast.SpelExpressionParser, expression string) {
	fmt.Printf("Token 流:\n")

	// 使用调试解析方法
	expr, err := parser.DoParseExpression(expression)
	if err != nil {
		fmt.Printf("❌ 解析错误: %v\n", err)
		return
	}

	fmt.Printf("✅ 解析成功!\n")
	fmt.Printf("表达式: %s\n", expr.GetExpressionString())
	fmt.Printf("AST 结构: %s\n", expr.ToStringAST())

	// 新增：递归打印 AST 树结构
	ast.PrintASTWithTitle(expr.AST, fmt.Sprintf("AST 树结构 - %s", expression))
}
