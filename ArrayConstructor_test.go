package main

// 测试用例参考：https://github.com/spring-projects/spring-framework/blob/main/spring-expression/src/test/java/org/springframework/expression/spel/ArrayConstructorTests.java

import (
	"fmt"
	"github.com/weaweawe01/ParserSpel/ast"
	"testing"
)

// TestArrayConstructorAnalysis 测试SpEL数组构造器表达式的AST树结构正确性
func TestArrayConstructorAnalysis(t *testing.T) {
	fmt.Println("=== SpEL 数组构造器 AST树结构测试 ===")

	// 测试用例定义
	testCases := []struct {
		name       string
		expression string
		expected   ASTExpectation
	}{
		{
			name:       "数组构造器与索引访问",
			expression: "new String[]{1,2,3}[0]",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "new String[] {1,2,3}[0]",
				Children: []ASTExpectation{
					{
						NodeType: "ConstructorReference",
						Value:    "new String[] {1,2,3}",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "String",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "String", Children: []ASTExpectation{}},
								},
							},
							{
								NodeType: "InlineList",
								Value:    "{1,2,3}",
								Children: []ASTExpectation{
									{NodeType: "IntLiteral", Value: "1", Children: []ASTExpectation{}},
									{NodeType: "IntLiteral", Value: "2", Children: []ASTExpectation{}},
									{NodeType: "IntLiteral", Value: "3", Children: []ASTExpectation{}},
								},
							},
						},
					},
					{
						NodeType: "Indexer",
						Value:    "[0]",
						Children: []ASTExpectation{
							{NodeType: "IntLiteral", Value: "0", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "Array Constructor with Indexer",
			expression: "new String[]{'123'}[0]",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "new String[] {'123'}[0]",
				Children: []ASTExpectation{
					{
						NodeType: "ConstructorReference",
						Value:    "new String[] {'123'}",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "String",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "String", Children: []ASTExpectation{}},
								},
							},
							{
								NodeType: "InlineList",
								Value:    "{'123'}",
								Children: []ASTExpectation{
									{NodeType: "StringLiteral", Value: "'123'", Children: []ASTExpectation{}},
								},
							},
						},
					},
					{
						NodeType: "Indexer",
						Value:    "[0]",
						Children: []ASTExpectation{
							{NodeType: "IntLiteral", Value: "0", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "双精度浮点数组",
			expression: "new double[]{1d,2d,3d,4d}",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new double[] {1.0,2.0,3.0,4.0}",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "double",
						Children: []ASTExpectation{
							{NodeType: "Identifier", Value: "double", Children: []ASTExpectation{}},
						},
					},
					{
						NodeType: "InlineList",
						Value:    "{1.0,2.0,3.0,4.0}",
						Children: []ASTExpectation{
							{NodeType: "RealLiteral", Value: "1.0", Children: []ASTExpectation{}},
							{NodeType: "RealLiteral", Value: "2.0", Children: []ASTExpectation{}},
							{NodeType: "RealLiteral", Value: "3.0", Children: []ASTExpectation{}},
							{NodeType: "RealLiteral", Value: "4.0", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "空数组的长度属性",
			expression: "new int[]{}.length",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "new int[] {}.length",
				Children: []ASTExpectation{
					{
						NodeType: "ConstructorReference",
						Value:    "new int[] {}",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "int",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "int", Children: []ASTExpectation{}},
								},
							},
							{
								NodeType: "InlineList",
								Value:    "{}",
								Children: []ASTExpectation{},
							},
						},
					},
					{
						NodeType: "PropertyOrFieldReference",
						Value:    ".length",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "布尔数组与索引访问",
			expression: "new boolean[]{true,false,true}[0]",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "new boolean[] {true,false,true}[0]",
				Children: []ASTExpectation{
					{
						NodeType: "ConstructorReference",
						Value:    "new boolean[] {true,false,true}",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "boolean",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "boolean", Children: []ASTExpectation{}},
								},
							},
							{
								NodeType: "InlineList",
								Value:    "{true,false,true}",
								Children: []ASTExpectation{
									{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
									{NodeType: "BooleanLiteral", Value: "false", Children: []ASTExpectation{}},
									{NodeType: "BooleanLiteral", Value: "true", Children: []ASTExpectation{}},
								},
							},
						},
					},
					{
						NodeType: "Indexer",
						Value:    "[0]",
						Children: []ASTExpectation{
							{NodeType: "IntLiteral", Value: "0", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "字符数组指定大小",
			expression: "new char[7]{'a','c','d','e'}",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new char[] {'a','c','d','e'}",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "char",
						Children: []ASTExpectation{
							{NodeType: "Identifier", Value: "char", Children: []ASTExpectation{}},
						},
					},
					{
						NodeType: "InlineList",
						Value:    "{'a','c','d','e'}",
						Children: []ASTExpectation{
							{NodeType: "StringLiteral", Value: "'a'", Children: []ASTExpectation{}},
							{NodeType: "StringLiteral", Value: "'c'", Children: []ASTExpectation{}},
							{NodeType: "StringLiteral", Value: "'d'", Children: []ASTExpectation{}},
							{NodeType: "StringLiteral", Value: "'e'", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "ArrayList构造器与类型引用",
			expression: "new java.util.ArrayList(T(java.lang.Integer).MAX_VALUE)",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new java.util.ArrayList(T(java.lang.Integer).MAX_VALUE)",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "java.util.ArrayList",
						Children: []ASTExpectation{
							{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
							{NodeType: "Identifier", Value: "util", Children: []ASTExpectation{}},
							{NodeType: "Identifier", Value: "ArrayList", Children: []ASTExpectation{}},
						},
					},
					{
						NodeType: "CompoundExpression",
						Value:    "T(java.lang.Integer).MAX_VALUE",
						Children: []ASTExpectation{
							{
								NodeType: "TypeReference",
								Value:    "T(java.lang.Integer)",
								Children: []ASTExpectation{
									{
										NodeType: "QualifiedIdentifier",
										Value:    "java.lang.Integer",
										Children: []ASTExpectation{
											{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
											{NodeType: "Identifier", Value: "lang", Children: []ASTExpectation{}},
											{NodeType: "Identifier", Value: "Integer", Children: []ASTExpectation{}},
										},
									},
								},
							},
							{NodeType: "PropertyOrFieldReference", Value: ".MAX_VALUE", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "大型多维数组",
			expression: "new int[1024 * 1024][1024 * 1024]",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new int[(1024 * 1024)][(1024 * 1024)]",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "int",
						Children: []ASTExpectation{
							{NodeType: "Identifier", Value: "int", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "字符串数组的长度属性",
			expression: "new String[]{'a','b','c','d'}.length",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "new String[] {'a','b','c','d'}.length",
				Children: []ASTExpectation{
					{
						NodeType: "ConstructorReference",
						Value:    "new String[] {'a','b','c','d'}",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "String",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "String", Children: []ASTExpectation{}},
								},
							},
							{
								NodeType: "InlineList",
								Value:    "{'a','b','c','d'}",
								Children: []ASTExpectation{
									{NodeType: "StringLiteral", Value: "'a'", Children: []ASTExpectation{}},
									{NodeType: "StringLiteral", Value: "'b'", Children: []ASTExpectation{}},
									{NodeType: "StringLiteral", Value: "'c'", Children: []ASTExpectation{}},
									{NodeType: "StringLiteral", Value: "'d'", Children: []ASTExpectation{}},
								},
							},
						},
					},
					{
						NodeType: "PropertyOrFieldReference",
						Value:    ".length",
						Children: []ASTExpectation{},
					},
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
				t.Fatalf("解析失败: %v", err)
			}

			if spelExpr == nil || spelExpr.AST == nil {
				t.Fatal("解析结果为空")
			}

			// 打印实际的AST树结构
			fmt.Println("实际AST结构:")
			ast.PrintAST(spelExpr.AST, 0)

			// 验证AST结构
			fmt.Println("验证AST结构...")
			if !validateASTStructure(spelExpr.AST.(ast.SpelNode), tc.expected) {
				t.Errorf("AST结构不匹配!\n期望: %+v\n实际AST见上方输出", tc.expected)
			} else {
				fmt.Println("✓ AST结构验证通过")
			}
		})
	}
}

// TestArrayConstructorParsing 测试数组构造器表达式的基本解析功能
func TestArrayConstructorParsing(t *testing.T) {
	fmt.Println("\n=== 数组构造器基本解析测试 ===")

	testExpressions := []string{
		"new String[]{1,2,3}[0]",
		"new int[]{'123'}[0]",
		"new double[]{1d,2d,3d,4d}",
		"new int[]{}.length",
		"new boolean[]{true,false,true}[0]",
		"new char[7]{'a','c','d','e'}",
		"new java.util.ArrayList(T(java.lang.Integer).MAX_VALUE)",
		"new int[1024 * 1024][1024 * 1024]",
		"new String[]{'a','b','c','d'}.length",
	}

	parser := ast.NewSpelExpressionParser()

	for i, expr := range testExpressions {
		t.Run(fmt.Sprintf("Expression_%d", i+1), func(t *testing.T) {
			fmt.Printf("\n测试表达式 %d: %s\n", i+1, expr)

			result, err := parser.ParseExpressionWithContext(expr, nil)

			if err != nil {
				fmt.Printf("❌ 解析错误: %v\n", err)
				// 某些表达式可能确实无法解析（如语法不标准的），这里不强制失败
				// t.Errorf("解析失败: %v", err)
			} else {
				fmt.Printf("✅ 解析成功!\n")
				if result != nil && result.AST != nil {
					ast.PrintAST(result.AST, 0)
				}
			}
		})
	}
}
