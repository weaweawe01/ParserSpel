package main

import (
	"fmt"
	"github.com/weaweawe01/ParserSpel/ast"
	"strings"
)

// 展示一个复杂表达式的详细 AST 分析
func demonstrateComplexAST() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🌳 复杂表达式 AST 树形结构展示")
	fmt.Println(strings.Repeat("=", 60))

	// 选择一个展示算术运算优先级和嵌套结构的表达式
	expression := "2 + 3 * 4 + (5 - 1) * 2"

	fmt.Printf("📝 表达式: %s\n", expression)
	fmt.Printf("📊 预期运算顺序: 2 + (3 * 4) + ((5 - 1) * 2)\n")
	fmt.Printf("🎯 预期结果: 2 + 12 + 8 = 22\n\n")

	parser := ast.NewSpelExpressionParser()

	fmt.Println("🔍 Token 分析:")
	fmt.Println("-------------")

	expr, err := parser.DoParseExpression(expression)
	if err != nil {
		fmt.Printf("❌ 解析失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 解析成功！\n")
	fmt.Printf("📋 线性表示: %s\n", expr.ToStringAST())

	// 显示完整的 AST 树结构
	ast.PrintASTWithTitle(expr.AST, "完整 AST 树形结构")

	fmt.Println("📝 树结构说明:")
	fmt.Println("- 每个缩进层级代表 AST 的深度")
	fmt.Println("- 节点类型显示了语法分析的结果")
	fmt.Println("- 表达式片段显示了每个节点对应的代码部分")
	fmt.Println("- 子节点按从左到右的顺序排列")

	fmt.Println("\n" + strings.Repeat("=", 60))
}
