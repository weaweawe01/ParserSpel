package main

// 当前测试用例参考：https://github.com/spring-projects/spring-framework/blob/main/spring-expression/src/test/java/org/springframework/expression/spel/ConstructorInvocationTests.java
// 测试构造函数调用的解析和AST结构

import (
	"fmt"
	"github.com/weaweawe01/ParserSpel/ast"
	"testing"
)

// TestConstructorInvocations 测试构造函数调用的解析
func TestConstructorInvocations(t *testing.T) {
	fmt.Println("=== 构造函数调用测试 ===")

	// 测试用例定义
	testCases := []struct {
		name       string
		expression string
		expected   ASTExpectation
	}{
		{
			name:       "带参数的构造函数",
			expression: "new String('hello world')",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new String('hello world')",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "String",
						Children: []ASTExpectation{
							{
								NodeType: "Identifier",
								Value:    "String",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "StringLiteral",
						Value:    "'hello world'",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "无参数构造函数",
			expression: "new String()",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new String()",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "String",
						Children: []ASTExpectation{
							{
								NodeType: "Identifier",
								Value:    "String",
								Children: []ASTExpectation{},
							},
						},
					},
				},
			},
		},
		{
			name:       "多参数构造函数",
			expression: "new String('hello', 'world')",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new String('hello', 'world')",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "String",
						Children: []ASTExpectation{
							{
								NodeType: "Identifier",
								Value:    "String",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "StringLiteral",
						Value:    "'hello'",
						Children: []ASTExpectation{},
					},
					{
						NodeType: "StringLiteral",
						Value:    "'world'",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "带数字参数的构造函数",
			expression: "new Double(3)",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new Double(3)",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "Double",
						Children: []ASTExpectation{
							{
								NodeType: "Identifier",
								Value:    "Double",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "IntLiteral",
						Value:    "3",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "带浮点数参数的构造函数",
			expression: "new Double(3.0)",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new Double(3.0)",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "Double",
						Children: []ASTExpectation{
							{
								NodeType: "Identifier",
								Value:    "Double",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "RealLiteral",
						Value:    "3.0",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "构造函数参数转换",
			expression: "new String(3.0)",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new String(3.0)",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "String",
						Children: []ASTExpectation{
							{
								NodeType: "Identifier",
								Value:    "String",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "RealLiteral",
						Value:    "3.0",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "包全限定名构造函数",
			expression: "new java.lang.String('test')",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new java.lang.String('test')",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "java.lang.String",
						Children: []ASTExpectation{
							{
								NodeType: "Identifier",
								Value:    "java",
								Children: []ASTExpectation{},
							},
							{
								NodeType: "Identifier",
								Value:    "lang",
								Children: []ASTExpectation{},
							},
							{
								NodeType: "Identifier",
								Value:    "String",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "StringLiteral",
						Value:    "'test'",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "带变量参数的构造函数",
			expression: "new String(#var)",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new String(#var)",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "String",
						Children: []ASTExpectation{
							{
								NodeType: "Identifier",
								Value:    "String",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "VariableReference",
						Value:    "#var",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "嵌套构造函数调用",
			expression: "new String(new String('nested'))",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new String(new String('nested'))",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "String",
						Children: []ASTExpectation{
							{
								NodeType: "Identifier",
								Value:    "String",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "ConstructorReference",
						Value:    "new String('nested')",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "String",
								Children: []ASTExpectation{
									{
										NodeType: "Identifier",
										Value:    "String",
										Children: []ASTExpectation{},
									},
								},
							},
							{
								NodeType: "StringLiteral",
								Value:    "'nested'",
								Children: []ASTExpectation{},
							},
						},
					},
				},
			},
		},
		{
			name:       "带运算表达式参数的构造函数",
			expression: "new Integer(1 + 2)",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new Integer((1 + 2))",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "Integer",
						Children: []ASTExpectation{
							{
								NodeType: "Identifier",
								Value:    "Integer",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "OpPlus",
						Value:    "(1 + 2)",
						Children: []ASTExpectation{
							{
								NodeType: "IntLiteral",
								Value:    "1",
								Children: []ASTExpectation{},
							},
							{
								NodeType: "IntLiteral",
								Value:    "2",
								Children: []ASTExpectation{},
							},
						},
					},
				},
			},
		},
	}

	// 执行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- 测试: %s ---\n", tc.name)
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

// TestComplexConstructorCalls 测试复杂的构造函数调用场景
func TestComplexConstructorCalls(t *testing.T) {
	fmt.Println("\n=== 复杂构造函数调用测试 ===")

	testCases := []struct {
		name       string
		expression string
		expected   ASTExpectation
	}{
		{
			name:       "可变参数构造函数",
			expression: "new Fruit('a','b','c')",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new Fruit('a', 'b', 'c')",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "Fruit",
						Children: []ASTExpectation{
							{
								NodeType: "Identifier",
								Value:    "Fruit",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "StringLiteral",
						Value:    "'a'",
						Children: []ASTExpectation{},
					},
					{
						NodeType: "StringLiteral",
						Value:    "'b'",
						Children: []ASTExpectation{},
					},
					{
						NodeType: "StringLiteral",
						Value:    "'c'",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "混合类型参数构造函数",
			expression: "new Fruit(1,'a',3.0)",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new Fruit(1, 'a', 3.0)",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "Fruit",
						Children: []ASTExpectation{
							{
								NodeType: "Identifier",
								Value:    "Fruit",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "IntLiteral",
						Value:    "1",
						Children: []ASTExpectation{},
					},
					{
						NodeType: "StringLiteral",
						Value:    "'a'",
						Children: []ASTExpectation{},
					},
					{
						NodeType: "RealLiteral",
						Value:    "3.0",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "带布尔参数的构造函数",
			expression: "new Boolean(true)",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new Boolean(true)",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "Boolean",
						Children: []ASTExpectation{
							{
								NodeType: "Identifier",
								Value:    "Boolean",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "BooleanLiteral",
						Value:    "true",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "带null参数的构造函数",
			expression: "new String(null)",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new String(null)",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "String",
						Children: []ASTExpectation{
							{
								NodeType: "Identifier",
								Value:    "String",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "NullLiteral",
						Value:    "null",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "构造函数链式调用",
			expression: "new StringBuilder().append('test')",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "new StringBuilder().append('test')",
				Children: []ASTExpectation{
					{
						NodeType: "ConstructorReference",
						Value:    "new StringBuilder()",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "StringBuilder",
								Children: []ASTExpectation{
									{
										NodeType: "Identifier",
										Value:    "StringBuilder",
										Children: []ASTExpectation{},
									},
								},
							},
						},
					},
					{
						NodeType: "MethodReference",
						Value:    "append('test')",
						Children: []ASTExpectation{
							{
								NodeType: "StringLiteral",
								Value:    "'test'",
								Children: []ASTExpectation{},
							},
						},
					},
				},
			},
		},
	}

	// 执行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- 测试: %s ---\n", tc.name)
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

// TestConstructorWithExpressions 测试构造函数与表达式的组合
func TestConstructorWithExpressions(t *testing.T) {
	fmt.Println("\n=== 构造函数与表达式组合测试 ===")

	testCases := []struct {
		name       string
		expression string
		expected   ASTExpectation
	}{
		{
			name:       "构造函数作为表达式的一部分",
			expression: "new String('hello') + ' world'",
			expected: ASTExpectation{
				NodeType: "OpPlus",
				Value:    "(new String('hello') + ' world')",
				Children: []ASTExpectation{
					{
						NodeType: "ConstructorReference",
						Value:    "new String('hello')",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "String",
								Children: []ASTExpectation{
									{
										NodeType: "Identifier",
										Value:    "String",
										Children: []ASTExpectation{},
									},
								},
							},
							{
								NodeType: "StringLiteral",
								Value:    "'hello'",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "StringLiteral",
						Value:    "' world'",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "构造函数比较",
			expression: "new Integer(5) > 3",
			expected: ASTExpectation{
				NodeType: "OpGT",
				Value:    "(new Integer(5) > 3)",
				Children: []ASTExpectation{
					{
						NodeType: "ConstructorReference",
						Value:    "new Integer(5)",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "Integer",
								Children: []ASTExpectation{
									{
										NodeType: "Identifier",
										Value:    "Integer",
										Children: []ASTExpectation{},
									},
								},
							},
							{
								NodeType: "IntLiteral",
								Value:    "5",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "IntLiteral",
						Value:    "3",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "构造函数在三元表达式中",
			expression: "true ? new String('yes') : new String('no')",
			expected: ASTExpectation{
				NodeType: "Ternary",
				Value:    "(true ? new String('yes') : new String('no'))",
				Children: []ASTExpectation{
					{
						NodeType: "BooleanLiteral",
						Value:    "true",
						Children: []ASTExpectation{},
					},
					{
						NodeType: "ConstructorReference",
						Value:    "new String('yes')",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "String",
								Children: []ASTExpectation{
									{
										NodeType: "Identifier",
										Value:    "String",
										Children: []ASTExpectation{},
									},
								},
							},
							{
								NodeType: "StringLiteral",
								Value:    "'yes'",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "ConstructorReference",
						Value:    "new String('no')",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "String",
								Children: []ASTExpectation{
									{
										NodeType: "Identifier",
										Value:    "String",
										Children: []ASTExpectation{},
									},
								},
							},
							{
								NodeType: "StringLiteral",
								Value:    "'no'",
								Children: []ASTExpectation{},
							},
						},
					},
				},
			},
		},
	}

	// 执行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- 测试: %s ---\n", tc.name)
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

// TestConstructorErrors 测试构造函数错误处理
func TestConstructorErrors(t *testing.T) {
	fmt.Println("\n=== 构造函数错误处理测试 ===")

	// 错误情况测试
	errorCases := []struct {
		name       string
		expression string
		shouldFail bool
	}{
		{
			name:       "缺少构造函数名",
			expression: "new ()",
			shouldFail: true,
		},
		{
			name:       "缺少括号",
			expression: "new String",
			shouldFail: true,
		},
		{
			name:       "不匹配的括号",
			expression: "new String(",
			shouldFail: true,
		},
		{
			name:       "空的new关键字",
			expression: "new",
			shouldFail: true,
		},
		{
			name:       "非法的类型名",
			expression: "new 123()",
			shouldFail: true,
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- 错误测试: %s ---\n", tc.name)
			fmt.Printf("表达式: %s\n", tc.expression)

			// 解析表达式
			parser := ast.NewSpelExpressionParser()
			_, err := parser.ParseExpression(tc.expression)

			if tc.shouldFail && err == nil {
				t.Errorf("期望解析失败，但实际成功了")
			} else if !tc.shouldFail && err != nil {
				t.Errorf("期望解析成功，但实际失败了: %v", err)
			} else if tc.shouldFail && err != nil {
				fmt.Printf("✓ 正确捕获到错误: %v\n", err)
			} else {
				fmt.Printf("✓ 解析成功\n")
			}
		})
	}
}

// TestConstructorInMethodChain 测试构造函数在方法链中的使用
func TestConstructorInMethodChain(t *testing.T) {
	fmt.Println("\n=== 构造函数方法链测试 ===")

	testCases := []struct {
		name       string
		expression string
		expected   ASTExpectation
	}{
		{
			name:       "构造函数后调用方法",
			expression: "new StringBuilder('hello').toString()",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "new StringBuilder('hello').toString()",
				Children: []ASTExpectation{
					{
						NodeType: "ConstructorReference",
						Value:    "new StringBuilder('hello')",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "StringBuilder",
								Children: []ASTExpectation{
									{
										NodeType: "Identifier",
										Value:    "StringBuilder",
										Children: []ASTExpectation{},
									},
								},
							},
							{
								NodeType: "StringLiteral",
								Value:    "'hello'",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "MethodReference",
						Value:    "toString()",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "构造函数方法链",
			expression: "new StringBuilder().append('a').append('b')",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "new StringBuilder().append('a').append('b')",
				Children: []ASTExpectation{
					{
						NodeType: "ConstructorReference",
						Value:    "new StringBuilder()",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "StringBuilder",
								Children: []ASTExpectation{
									{
										NodeType: "Identifier",
										Value:    "StringBuilder",
										Children: []ASTExpectation{},
									},
								},
							},
						},
					},
					{
						NodeType: "MethodReference",
						Value:    "append('a')",
						Children: []ASTExpectation{
							{
								NodeType: "StringLiteral",
								Value:    "'a'",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "MethodReference",
						Value:    "append('b')",
						Children: []ASTExpectation{
							{
								NodeType: "StringLiteral",
								Value:    "'b'",
								Children: []ASTExpectation{},
							},
						},
					},
				},
			},
		},
	}

	// 执行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- 测试: %s ---\n", tc.name)
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
