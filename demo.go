package main

import (
	"fmt"
	"github.com/weaweawe01/ParserSpel/ast"
	"log"
)

// 演示 SpEL 解析器的高级功能
func demonstrateAdvancedFeatures() {
	fmt.Println("高级 SpEL 功能演示")
	fmt.Println("==================")

	parser := ast.NewSpelExpressionParser()

	// 测试各种表达式类型
	examples := map[string]string{
		"数值计算":   "((10 + 5) * 2) / 3",
		"字符串操作": "'Hello' + ' ' + 'World'",
		"布尔逻辑":   "(true || false) && !false",
		"比较运算":   "42 > 30 && 'abc' == 'abc'",
		"正则匹配":   "'Test123' matches '[A-Z][a-z]+[0-9]+'",
		"混合表达式": "5 > 3 && 'test' == 'test'",
		"复杂逻辑":   "(age > 18 && age < 65) || isVip == true",
		"属性访问":   "user.profile.name",
		"安全导航":   "user?.profile?.email",
		"变量引用":   "#root.getValue()",
		"Bean引用":   "@userService.findById(123)",
	}

	for description, expression := range examples {
		fmt.Printf("\n【%s】\n", description)
		fmt.Printf("表达式: %s\n", expression)

		expr, err := parser.ParseExpression(expression)
		if err != nil {
			fmt.Printf("解析错误: %v\n", err)
			continue
		}

		fmt.Printf("AST结构: %s\n", expr.ToStringAST())

		value, err := expr.GetValue()
		if err != nil {
			fmt.Printf("求值错误: %v\n", err)
		} else {
			fmt.Printf("计算结果: %v (%T)\n", value, value)
		}
	}
}

// 演示语法分析的详细过程
func demonstrateParsingDetails() {
	fmt.Println("\n\n语法分析详细过程")
	fmt.Println("==================")

	testExpression := "2 + 3 * 4 > 10 && true"

	fmt.Printf("测试表达式: %s\n\n", testExpression)

	// 1. 词法分析
	fmt.Println("1. 词法分析结果:")
	tokenizer := ast.NewTokenizer(testExpression)
	tokens, err := tokenizer.Process()
	if err != nil {
		log.Printf("词法分析错误: %v\n", err)
		return
	}

	for i, token := range tokens {
		fmt.Printf("   Token[%d]: %s\n", i, token.String())
	}

	// 2. 语法分析
	fmt.Println("\n2. 语法分析结果:")
	parser := ast.NewSpelExpressionParser()
	expr, err := parser.ParseExpression(testExpression)
	if err != nil {
		fmt.Printf("语法分析错误: %v\n", err)
		return
	}

	fmt.Printf("   AST: %s\n", expr.ToStringAST())

	// 3. 表达式求值
	fmt.Println("\n3. 表达式求值:")
	value, err := expr.GetValue()
	if err != nil {
		fmt.Printf("求值错误: %v\n", err)
	} else {
		fmt.Printf("   结果: %v (%T)\n", value, value)
	}

	// 4. 解析步骤说明
	fmt.Println("\n4. 解析步骤说明:")
	fmt.Println("   ① 2 + 3 * 4 → 2 + (3 * 4) → 2 + 12 → 14")
	fmt.Println("   ② 14 > 10 → true")
	fmt.Println("   ③ true && true → true")
}

// 验证运算符优先级
func demonstrateOperatorPrecedence() {
	fmt.Println("\n\n运算符优先级验证")
	fmt.Println("==================")

	parser := ast.NewSpelExpressionParser()

	testCases := []struct {
		expression  string
		expected    interface{}
		description string
	}{
		{"2 + 3 * 4", int64(14), "乘法优先于加法"},
		{"(2 + 3) * 4", int64(20), "括号改变优先级"},
		{"10 - 6 / 2", int64(7), "除法优先于减法"},
		{"true || false && false", true, "AND 优先于 OR"},
		{"!false && true", true, "NOT 优先于 AND"},
		{"5 > 3 && 2 < 4", true, "比较优先于逻辑"},
	}

	for i, tc := range testCases {
		fmt.Printf("\n测试 %d: %s\n", i+1, tc.description)
		fmt.Printf("表达式: %s\n", tc.expression)

		expr, err := parser.ParseExpression(tc.expression)
		if err != nil {
			fmt.Printf("解析错误: %v\n", err)
			continue
		}

		value, err := expr.GetValue()
		if err != nil {
			fmt.Printf("求值错误: %v\n", err)
			continue
		}

		fmt.Printf("实际结果: %v\n", value)
		fmt.Printf("预期结果: %v\n", tc.expected)

		if fmt.Sprintf("%v", value) == fmt.Sprintf("%v", tc.expected) {
			fmt.Println("✅ 测试通过")
		} else {
			fmt.Println("❌ 测试失败")
		}
	}
}

func runAdvancedDemo() {
	demonstrateAdvancedFeatures()
	demonstrateParsingDetails()
	demonstrateOperatorPrecedence()
}
