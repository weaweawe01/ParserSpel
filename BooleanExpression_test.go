package main

// 当前测试用例 参考:https://github.com/spring-projects/spring-framework/blob/main/spring-expression/src/test/java/org/springframework/expression/spel/BooleanExpressionTests.java
// 已全部测试完毕

import (
	"fmt"
	"github.com/weaweawe01/ParserSpel/ast"
	"testing"
)

// TestBooleanExpressions 测试布尔表达式的解析和求值
func TestBooleanExpressions(t *testing.T) {
	fmt.Println("=== 布尔表达式测试 ===")

	// 测试用例定义
	testCases := []struct {
		name       string
		expression string
		expected   ASTExpectation
	}{
		{
			name:       "布尔值 true",
			expression: "true",
			expected: ASTExpectation{
				NodeType: "BooleanLiteral",
				Value:    "true",
				Children: []ASTExpectation{},
			},
		},
		{
			name:       "布尔值 false",
			expression: "false",
			expected: ASTExpectation{
				NodeType: "BooleanLiteral",
				Value:    "false",
				Children: []ASTExpectation{},
			},
		},
		{
			name:       "OR 运算 - false or false",
			expression: "false or false",
			expected: ASTExpectation{
				NodeType: "OpOr",
				Value:    "(false or false)",
				Children: []ASTExpectation{
					{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
					{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "OR 运算 - false or true",
			expression: "false or true",
			expected: ASTExpectation{
				NodeType: "OpOr",
				Value:    "(false or true)",
				Children: []ASTExpectation{
					{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
					{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "OR 运算 - true or false",
			expression: "true or false",
			expected: ASTExpectation{
				NodeType: "OpOr",
				Value:    "(true or false)",
				Children: []ASTExpectation{
					{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
					{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "OR 运算 - true or true",
			expression: "true or true",
			expected: ASTExpectation{
				NodeType: "OpOr",
				Value:    "(true or true)",
				Children: []ASTExpectation{
					{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
					{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "AND 运算 - false and false",
			expression: "false and false",
			expected: ASTExpectation{
				NodeType: "OpAnd",
				Value:    "(false and false)",
				Children: []ASTExpectation{
					{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
					{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "AND 运算 - false and true",
			expression: "false and true",
			expected: ASTExpectation{
				NodeType: "OpAnd",
				Value:    "(false and true)",
				Children: []ASTExpectation{
					{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
					{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "AND 运算 - true and false",
			expression: "true and false",
			expected: ASTExpectation{
				NodeType: "OpAnd",
				Value:    "(true and false)",
				Children: []ASTExpectation{
					{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
					{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "AND 运算 - true and true",
			expression: "true and true",
			expected: ASTExpectation{
				NodeType: "OpAnd",
				Value:    "(true and true)",
				Children: []ASTExpectation{
					{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
					{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "NOT 运算 - !false",
			expression: "!false",
			expected: ASTExpectation{
				NodeType: "OperatorNot",
				Value:    "!false",
				Children: []ASTExpectation{
					{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "NOT 运算 - !true",
			expression: "!true",
			expected: ASTExpectation{
				NodeType: "OperatorNot",
				Value:    "!true",
				Children: []ASTExpectation{
					{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "NOT 运算 - not false",
			expression: "not false",
			expected: ASTExpectation{
				NodeType: "OperatorNot",
				Value:    "!false",
				Children: []ASTExpectation{
					{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "NOT 运算 - NoT true (大小写不敏感)",
			expression: "NoT true",
			expected: ASTExpectation{
				NodeType: "OperatorNot",
				Value:    "!true",
				Children: []ASTExpectation{
					{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "组合运算 - false and false or true",
			expression: "false and false or true",
			expected: ASTExpectation{
				NodeType: "OpOr",
				Value:    "((false and false) or true)",
				Children: []ASTExpectation{
					{
						NodeType: "OpAnd",
						Value:    "(false and false)",
						Children: []ASTExpectation{
							{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
							{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
						},
					},
					{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "组合运算 - true and false or true",
			expression: "true and false or true",
			expected: ASTExpectation{
				NodeType: "OpOr",
				Value:    "((true and false) or true)",
				Children: []ASTExpectation{
					{
						NodeType: "OpAnd",
						Value:    "(true and false)",
						Children: []ASTExpectation{
							{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
							{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
						},
					},
					{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "组合运算 - true and false or false",
			expression: "true and false or false",
			expected: ASTExpectation{
				NodeType: "OpOr",
				Value:    "((true and false) or false)",
				Children: []ASTExpectation{
					{
						NodeType: "OpAnd",
						Value:    "(true and false)",
						Children: []ASTExpectation{
							{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
							{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
						},
					},
					{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "复杂布尔表达式 - (true or false) and !false",
			expression: "(true or false) and !false",
			expected: ASTExpectation{
				NodeType: "OpAnd",
				Value:    "((true or false) and !false)",
				Children: []ASTExpectation{
					{
						NodeType: "OpOr",
						Value:    "(true or false)",
						Children: []ASTExpectation{
							{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
							{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
						},
					},
					{
						NodeType: "OperatorNot",
						Value:    "!false",
						Children: []ASTExpectation{
							{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "三元条件运算符 - true ? 'foo' : 'bar'",
			expression: "true ? 'foo' : 'bar'",
			expected: ASTExpectation{
				NodeType: "Ternary",
				Value:    "(true ? 'foo' : 'bar')",
				Children: []ASTExpectation{
					{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
					{NodeType: "StringLiteral", Value: "'foo'", Children: []ASTExpectation{}},
					{NodeType: "StringLiteral", Value: "'bar'", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "三元条件运算符 - false ? 'foo' : 'bar'",
			expression: "false ? 'foo' : 'bar'",
			expected: ASTExpectation{
				NodeType: "Ternary",
				Value:    "(false ? 'foo' : 'bar')",
				Children: []ASTExpectation{
					{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
					{NodeType: "StringLiteral", Value: "'foo'", Children: []ASTExpectation{}},
					{NodeType: "StringLiteral", Value: "'bar'", Children: []ASTExpectation{}},
				},
			},
		},
	}

	// 执行测试
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- 测试用例: %s ---\n", tc.name)
			fmt.Printf("表达式: %s\n", tc.expression)

			// 创建解析器
			parser := ast.NewSpelExpressionParser()

			// 解析表达式
			spelExpr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				t.Fatalf("解析表达式失败: %v", err)
			}

			// 检查AST结构
			fmt.Printf("AST-> %v\n", spelExpr.AST)
			fmt.Printf("%s\n", spelExpr.GetExpressionString())

			// 输出实际AST结构以便调试
			fmt.Println("实际AST结构:")
			ast.PrintAST(spelExpr.AST, 0)

			// 验证AST结构
			fmt.Println("验证AST结构...")
			if !validateASTStructure(spelExpr.AST, tc.expected) {
				t.Errorf("AST结构验证失败")
			} else {
				fmt.Println("✅ AST结构验证通过")
			}
		})
	}
}

// TestBooleanExpressionEvaluation 测试布尔表达式的求值
func TestBooleanExpressionEvaluation(t *testing.T) {
	fmt.Println("\n=== 布尔表达式求值测试 ===")

	// 求值测试用例
	evaluationTests := []struct {
		name       string
		expression string
		expected   interface{}
	}{
		{"true", "true", true},
		{"false", "false", false},
		{"false or false", "false or false", false},
		{"false or true", "false or true", true},
		{"true or false", "true or false", true},
		{"true or true", "true or true", true},
		{"false and false", "false and false", false},
		{"false and true", "false and true", false},
		{"true and false", "true and false", false},
		{"true and true", "true and true", true},
		{"!false", "!false", true},
		{"!true", "!true", false},
		{"not false", "not false", true},
		{"not true", "not true", false},
		{"false and false or true", "false and false or true", true},
		{"true and false or true", "true and false or true", true},
		{"true and false or false", "true and false or false", false},
	}

	parser := ast.NewSpelExpressionParser()

	for _, tc := range evaluationTests {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- 求值测试: %s ---\n", tc.expression)

			expr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				t.Fatalf("解析表达式失败: %v", err)
			}

			result, err := expr.GetValue()
			if err != nil {
				t.Fatalf("求值失败: %v", err)
			}

			fmt.Printf("表达式: %s\n", tc.expression)
			fmt.Printf("期望结果: %v\n", tc.expected)
			fmt.Printf("实际结果: %v\n", result)

			if result != tc.expected {
				t.Errorf("求值结果不匹配: 期望 %v, 实际 %v", tc.expected, result)
			} else {
				fmt.Println("✅ 求值结果正确")
			}
		})
	}
}

// TestBooleanErrorCases 测试布尔表达式的错误情况
func TestBooleanErrorCases(t *testing.T) {
	fmt.Println("\n=== 布尔表达式错误情况测试 ===")

	// 错误情况测试用例
	errorTests := []struct {
		name       string
		expression string
		shouldFail bool
		errorMsg   string
	}{
		{
			name:       "数字与布尔值OR运算",
			expression: "1.0 or false",
			shouldFail: true,
			errorMsg:   "类型转换错误",
		},
		{
			name:       "布尔值与数字OR运算",
			expression: "false or 39.4",
			shouldFail: true,
			errorMsg:   "类型转换错误",
		},
		{
			name:       "布尔值与字符串AND运算",
			expression: "true and 'hello'",
			shouldFail: true,
			errorMsg:   "类型转换错误",
		},
		{
			name:       "字符串之间AND运算",
			expression: "'hello' and 'goodbye'",
			shouldFail: true,
			errorMsg:   "类型转换错误",
		},
		{
			name:       "数字NOT运算",
			expression: "!35.2",
			shouldFail: true,
			errorMsg:   "类型转换错误",
		},
		{
			name:       "字符串NOT运算",
			expression: "!'foob'",
			shouldFail: true,
			errorMsg:   "类型转换错误",
		},
	}

	parser := ast.NewSpelExpressionParser()

	for _, tc := range errorTests {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- 错误测试: %s ---\n", tc.name)
			fmt.Printf("表达式: %s\n", tc.expression)

			expr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				if tc.shouldFail {
					fmt.Printf("✅ 预期的解析错误: %v\n", err)
					return
				} else {
					t.Fatalf("意外的解析错误: %v", err)
				}
			}

			// 如果解析成功，尝试求值
			result, err := expr.GetValue()

			if tc.shouldFail {
				if err != nil {
					fmt.Printf("✅ 预期的求值错误: %v\n", err)
				} else {
					t.Errorf("应该失败但成功了，结果: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("不应该失败但失败了: %v", err)
				} else {
					fmt.Printf("✅ 成功求值: %v\n", result)
				}
			}
		})
	}
}
