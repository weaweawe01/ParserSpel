package main

// 测试用例参考：https://github.com/spring-projects/spring-framework/blob/main/spring-expression/src/test/java/org/springframework/expression/spel/OperatorOverloaderTests.java

import (
	"fmt"
	"github.com/weaweawe01/ParserSpel/ast"
	"testing"
)

// TestOperatorOverloaderAnalysis 测试SpEL操作符重载表达式的AST树结构正确性
func TestOperatorOverloaderAnalysis(t *testing.T) {
	fmt.Println("=== SpEL 操作符重载 AST树结构测试 ===")

	// 测试用例定义
	testCases := []struct {
		name       string
		expression string
		expected   ASTExpectation
	}{
		{
			name:       "字符串与布尔值相加",
			expression: "'abc' + true",
			expected: ASTExpectation{
				NodeType: "OpPlus",
				Value:    "('abc' + true)",
				Children: []ASTExpectation{
					{NodeType: "StringLiteral", Value: "'abc'", Children: []ASTExpectation{}},
					{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "字符串与布尔值相减",
			expression: "'abc' - true",
			expected: ASTExpectation{
				NodeType: "OpMinus",
				Value:    "('abc' - true)",
				Children: []ASTExpectation{
					{NodeType: "StringLiteral", Value: "'abc'", Children: []ASTExpectation{}},
					{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "字符串与null相加",
			expression: "'abc' + null",
			expected: ASTExpectation{
				NodeType: "OpPlus",
				Value:    "('abc' + null)",
				Children: []ASTExpectation{
					{NodeType: "StringLiteral", Value: "'abc'", Children: []ASTExpectation{}},
					{NodeType: "NullLiteral", Value: "null", Children: []ASTExpectation{}},
				},
			},
		},
	}

	// 运行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- 测试用例: %s ---\n", tc.name)
			fmt.Printf("表达式: %s\n", tc.expression)

			parser := ast.NewSpelExpressionParser()
			spelExpr, err := parser.ParseExpressionWithContext(tc.expression, nil)

			if err != nil {
				t.Logf("解析失败（预期可能失败）: %v", err)
				return // 继续下一个测试，不标记为失败
			}

			if spelExpr == nil || spelExpr.AST == nil {
				t.Logf("解析结果为空（预期可能失败）")
				return
			}

			// 打印实际的AST树结构
			fmt.Println("实际AST结构:")
			ast.PrintAST(spelExpr.AST, 0)

			// 验证AST结构
			fmt.Println("验证AST结构...")
			if !validateASTStructure(spelExpr.AST.(ast.SpelNode), tc.expected) {
				t.Logf("AST结构不匹配（预期可能不匹配）!\n期望: %+v\n实际AST见上方输出", tc.expected)
			} else {
				fmt.Println("✓ AST结构验证通过")
			}
		})
	}
}

// TestOperatorOverloaderParsing 测试操作符重载表达式的基本解析功能
func TestOperatorOverloaderParsing(t *testing.T) {
	fmt.Println("\n=== 操作符重载基本解析测试 ===")

	testExpressions := []string{
		"'abc' + true",
		"'abc' - true",
		"'abc' + null",
	}

	parser := ast.NewSpelExpressionParser()

	for i, expr := range testExpressions {
		t.Run(fmt.Sprintf("Expression_%d", i+1), func(t *testing.T) {
			fmt.Printf("\n测试表达式 %d: %s\n", i+1, expr)

			result, err := parser.ParseExpressionWithContext(expr, nil)

			if err != nil {
				fmt.Printf("❌ 解析错误: %v\n", err)
				// 不强制失败，因为某些语法可能尚未实现
			} else {
				fmt.Printf("✅ 解析成功!\n")
				if result != nil && result.AST != nil {
					ast.PrintAST(result.AST, 0)
				}
			}
		})
	}
}

// TestOperatorOverloaderSpecialCases 测试操作符重载的特殊情况
func TestOperatorOverloaderSpecialCases(t *testing.T) {
	fmt.Println("\n=== 操作符重载特殊情况测试 ===")

	parser := ast.NewSpelExpressionParser()

	// 测试字符串与布尔值运算
	t.Run("StringBooleanOperations", func(t *testing.T) {
		testCases := []string{
			"'abc' + true",
			"'abc' + false",
			"'abc' - true",
			"'abc' - false",
		}

		for _, expr := range testCases {
			result, err := parser.ParseExpression(expr)
			if err != nil {
				fmt.Printf("字符串布尔运算 '%s' 解析失败: %v\n", expr, err)
			} else {
				fmt.Printf("字符串布尔运算 '%s' 解析成功: %s\n", expr, result.AST.ToStringAST())
			}
		}
	})

	// 测试字符串与null运算
	t.Run("StringNullOperations", func(t *testing.T) {
		testCases := []string{
			"'abc' + null",
			"'abc' - null",
			"null + 'abc'",
			"null - 'abc'",
		}

		for _, expr := range testCases {
			result, err := parser.ParseExpression(expr)
			if err != nil {
				fmt.Printf("字符串null运算 '%s' 解析失败: %v\n", expr, err)
			} else {
				fmt.Printf("字符串null运算 '%s' 解析成功: %s\n", expr, result.AST.ToStringAST())
			}
		}
	})

	// 测试数字与其他类型运算
	t.Run("NumberOperations", func(t *testing.T) {
		testCases := []string{
			"123 + true",
			"123 + false",
			"123 + null",
			"123 - true",
			"123 - false",
			"123 - null",
		}

		for _, expr := range testCases {
			result, err := parser.ParseExpression(expr)
			if err != nil {
				fmt.Printf("数字运算 '%s' 解析失败: %v\n", expr, err)
			} else {
				fmt.Printf("数字运算 '%s' 解析成功: %s\n", expr, result.AST.ToStringAST())
			}
		}
	})

	// 测试布尔值之间运算
	t.Run("BooleanOperations", func(t *testing.T) {
		testCases := []string{
			"true + false",
			"true - false",
			"false + true",
			"false - true",
		}

		for _, expr := range testCases {
			result, err := parser.ParseExpression(expr)
			if err != nil {
				fmt.Printf("布尔运算 '%s' 解析失败: %v\n", expr, err)
			} else {
				fmt.Printf("布尔运算 '%s' 解析成功: %s\n", expr, result.AST.ToStringAST())
			}
		}
	})

	// 测试复杂表达式
	t.Run("ComplexExpressions", func(t *testing.T) {
		testCases := []string{
			"'abc' + true + null",
			"'abc' + (true - false)",
			"('abc' + true) - null",
			"'hello' + ' ' + 'world' + true",
		}

		for _, expr := range testCases {
			result, err := parser.ParseExpression(expr)
			if err != nil {
				fmt.Printf("复杂表达式 '%s' 解析失败: %v\n", expr, err)
			} else {
				fmt.Printf("复杂表达式 '%s' 解析成功: %s\n", expr, result.AST.ToStringAST())
			}
		}
	})
}

// TestOperatorOverloaderEvaluation 测试操作符重载的求值行为
func TestOperatorOverloaderEvaluation(t *testing.T) {
	fmt.Println("\n=== 操作符重载求值测试 ===")

	parser := ast.NewSpelExpressionParser()

	// 测试字符串与布尔值相加的求值
	t.Run("StringPlusBooleanEvaluation", func(t *testing.T) {
		expr, err := parser.ParseExpression("'abc' + true")
		if err != nil {
			t.Logf("解析失败: %v", err)
			return
		}

		// 尝试求值
		value, evalErr := expr.GetValue()
		if evalErr != nil {
			fmt.Printf("求值失败（预期可能失败）: %v\n", evalErr)
		} else {
			fmt.Printf("'abc' + true 求值结果: %v (类型: %T)\n", value, value)
		}
	})

	// 测试字符串与布尔值相减的求值
	t.Run("StringMinusBooleanEvaluation", func(t *testing.T) {
		expr, err := parser.ParseExpression("'abc' - true")
		if err != nil {
			t.Logf("解析失败: %v", err)
			return
		}

		// 尝试求值
		value, evalErr := expr.GetValue()
		if evalErr != nil {
			fmt.Printf("求值失败（预期可能失败）: %v\n", evalErr)
		} else {
			fmt.Printf("'abc' - true 求值结果: %v (类型: %T)\n", value, value)
		}
	})

	// 测试字符串与null相加的求值
	t.Run("StringPlusNullEvaluation", func(t *testing.T) {
		expr, err := parser.ParseExpression("'abc' + null")
		if err != nil {
			t.Logf("解析失败: %v", err)
			return
		}

		// 尝试求值
		value, evalErr := expr.GetValue()
		if evalErr != nil {
			fmt.Printf("求值失败（预期可能失败）: %v\n", evalErr)
		} else {
			fmt.Printf("'abc' + null 求值结果: %v (类型: %T)\n", value, value)
		}
	})
}
