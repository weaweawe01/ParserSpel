package main

import (
	"fmt"
	"github.com/weaweawe01/ParserSpel/ast"
	"testing"
)

// TestSpelParserAnalysis 测试SpEL AST树结构的正确性
func TestSpelParserMapAnalysis(t *testing.T) {
	fmt.Println("=== SpEL AST树结构测试 ===")

	// 测试用例定义
	testCases := []struct {
		name       string
		expression string
		expected   ASTExpectation
	}{
		{
			name:       "数组索引访问",
			expression: "inventions[3]",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "inventions[3]",
				Children: []ASTExpectation{
					{NodeType: "PropertyOrFieldReference", Value: "inventions", Children: []ASTExpectation{}},
					{
						NodeType: "Indexer",
						Value:    "[3]",
						Children: []ASTExpectation{
							{NodeType: "IntLiteral", Value: "3", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "成员索引和属性访问",
			expression: "Members[0].Name",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "Members[0].Name",
				Children: []ASTExpectation{
					{NodeType: "PropertyOrFieldReference", Value: "Members", Children: []ASTExpectation{}},
					{
						NodeType: "Indexer",
						Value:    "[0]",
						Children: []ASTExpectation{
							{NodeType: "IntLiteral", Value: "0", Children: []ASTExpectation{}},
						},
					},
					{NodeType: "PropertyOrFieldReference", Value: ".Name", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "映射键访问",
			expression: "Officers['president']",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "Officers['president']",
				Children: []ASTExpectation{
					{NodeType: "PropertyOrFieldReference", Value: "Officers", Children: []ASTExpectation{}},
					{
						NodeType: "Indexer",
						Value:    "['president']",
						Children: []ASTExpectation{
							{NodeType: "StringLiteral", Value: "'president'", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "嵌套映射和索引访问",
			expression: "Officers['advisors'][0].PlaceOfBirth.Country",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "Officers['advisors'][0].PlaceOfBirth.Country",
				Children: []ASTExpectation{
					{NodeType: "PropertyOrFieldReference", Value: "Officers", Children: []ASTExpectation{}},
					{
						NodeType: "Indexer",
						Value:    "['advisors']",
						Children: []ASTExpectation{
							{NodeType: "StringLiteral", Value: "'advisors'", Children: []ASTExpectation{}},
						},
					},
					{
						NodeType: "Indexer",
						Value:    "[0]",
						Children: []ASTExpectation{
							{NodeType: "IntLiteral", Value: "0", Children: []ASTExpectation{}},
						},
					},
					{NodeType: "PropertyOrFieldReference", Value: ".PlaceOfBirth", Children: []ASTExpectation{}},
					{NodeType: "PropertyOrFieldReference", Value: ".Country", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "内联列表",
			expression: "{1,2,3,4}",
			expected: ASTExpectation{
				NodeType: "InlineList",
				Value:    "{1,2,3,4}",
				Children: []ASTExpectation{
					{NodeType: "IntLiteral", Value: "1", Children: []ASTExpectation{}},
					{NodeType: "IntLiteral", Value: "2", Children: []ASTExpectation{}},
					{NodeType: "IntLiteral", Value: "3", Children: []ASTExpectation{}},
					{NodeType: "IntLiteral", Value: "4", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "嵌套列表",
			expression: "{{'a','b'},{'x','y'}}",
			expected: ASTExpectation{
				NodeType: "InlineList",
				Value:    "{{'a','b'},{'x','y'}}",
				Children: []ASTExpectation{
					{
						NodeType: "InlineList",
						Value:    "{'a','b'}",
						Children: []ASTExpectation{
							{NodeType: "StringLiteral", Value: "'a'", Children: []ASTExpectation{}},
							{NodeType: "StringLiteral", Value: "'b'", Children: []ASTExpectation{}},
						},
					},
					{
						NodeType: "InlineList",
						Value:    "{'x','y'}",
						Children: []ASTExpectation{
							{NodeType: "StringLiteral", Value: "'x'", Children: []ASTExpectation{}},
							{NodeType: "StringLiteral", Value: "'y'", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "映射字面量",
			expression: "{name:{first:'Nikola',last:'Tesla'},dob:{day:10,month:'July',year:1856}}",
			expected: ASTExpectation{
				NodeType: "InlineMap",
				Value:    "{name:{first:'Nikola',last:'Tesla'},dob:{day:10,month:'July',year:1856}}",
				Children: []ASTExpectation{
					{NodeType: "PropertyOrFieldReference", Value: "name", Children: []ASTExpectation{}},
					{
						NodeType: "InlineMap",
						Value:    "{first:'Nikola',last:'Tesla'}",
						Children: []ASTExpectation{
							{NodeType: "PropertyOrFieldReference", Value: "first", Children: []ASTExpectation{}},
							{NodeType: "StringLiteral", Value: "'Nikola'", Children: []ASTExpectation{}},
							{NodeType: "PropertyOrFieldReference", Value: "last", Children: []ASTExpectation{}},
							{NodeType: "StringLiteral", Value: "'Tesla'", Children: []ASTExpectation{}},
						},
					},
					{NodeType: "PropertyOrFieldReference", Value: "dob", Children: []ASTExpectation{}},
					{
						NodeType: "InlineMap",
						Value:    "{day:10,month:'July',year:1856}",
						Children: []ASTExpectation{
							{NodeType: "PropertyOrFieldReference", Value: "day", Children: []ASTExpectation{}},
							{NodeType: "IntLiteral", Value: "10", Children: []ASTExpectation{}},
							{NodeType: "PropertyOrFieldReference", Value: "month", Children: []ASTExpectation{}},
							{NodeType: "StringLiteral", Value: "'July'", Children: []ASTExpectation{}},
							{NodeType: "PropertyOrFieldReference", Value: "year", Children: []ASTExpectation{}},
							{NodeType: "IntLiteral", Value: "1856", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "数组构造器",
			expression: "new int[]{1,2,3}",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new int[] {1,2,3}", // 注意实际输出包含空格
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
						Value:    "{1,2,3}",
						Children: []ASTExpectation{
							{NodeType: "IntLiteral", Value: "1", Children: []ASTExpectation{}},
							{NodeType: "IntLiteral", Value: "2", Children: []ASTExpectation{}},
							{NodeType: "IntLiteral", Value: "3", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "方法引用",
			expression: "isMember('Mihajlo Pupin')",
			expected: ASTExpectation{
				NodeType: "MethodReference",
				Value:    "isMember('Mihajlo Pupin')",
				Children: []ASTExpectation{
					{NodeType: "StringLiteral", Value: "'Mihajlo Pupin'", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "字符串方法调用",
			expression: "'abc'.substring(1, 3)",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "'abc'.substring(1, 3)",
				Children: []ASTExpectation{
					{NodeType: "StringLiteral", Value: "'abc'", Children: []ASTExpectation{}},
					{
						NodeType: "MethodReference",
						Value:    "substring(1, 3)",
						Children: []ASTExpectation{
							{NodeType: "IntLiteral", Value: "1", Children: []ASTExpectation{}},
							{NodeType: "IntLiteral", Value: "3", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "赋值表达式字符串",
			expression: "Name = 'Aleksandar Seovic2'",
			expected: ASTExpectation{
				NodeType: "Assign",
				Value:    "Name = 'Aleksandar Seovic2'",
				Children: []ASTExpectation{
					{NodeType: "PropertyOrFieldReference", Value: "Name", Children: []ASTExpectation{}},
					{NodeType: "StringLiteral", Value: "'Aleksandar Seovic2'", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "枚举比较",
			expression: "T(java.math.RoundingMode).CEILING < T(java.math.RoundingMode).FLOOR",
			expected: ASTExpectation{
				NodeType: "OpLT",
				Value:    "(T(java.math.RoundingMode).CEILING < T(java.math.RoundingMode).FLOOR)",
				Children: []ASTExpectation{
					{
						NodeType: "CompoundExpression",
						Value:    "T(java.math.RoundingMode).CEILING",
						Children: []ASTExpectation{
							{
								NodeType: "TypeReference",
								Value:    "T(java.math.RoundingMode)",
								Children: []ASTExpectation{
									{
										NodeType: "QualifiedIdentifier",
										Value:    "java.math.RoundingMode",
										Children: []ASTExpectation{
											{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
											{NodeType: "Identifier", Value: "math", Children: []ASTExpectation{}},
											{NodeType: "Identifier", Value: "RoundingMode", Children: []ASTExpectation{}},
										},
									},
								},
							},
							{NodeType: "PropertyOrFieldReference", Value: ".CEILING", Children: []ASTExpectation{}},
						},
					},
					{
						NodeType: "CompoundExpression",
						Value:    "T(java.math.RoundingMode).FLOOR",
						Children: []ASTExpectation{
							{
								NodeType: "TypeReference",
								Value:    "T(java.math.RoundingMode)",
								Children: []ASTExpectation{
									{
										NodeType: "QualifiedIdentifier",
										Value:    "java.math.RoundingMode",
										Children: []ASTExpectation{
											{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
											{NodeType: "Identifier", Value: "math", Children: []ASTExpectation{}},
											{NodeType: "Identifier", Value: "RoundingMode", Children: []ASTExpectation{}},
										},
									},
								},
							},
							{NodeType: "PropertyOrFieldReference", Value: ".FLOOR", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "构造器表达式",
			expression: "new org.spring.samples.spel.inventor.Inventor('Albert Einstein', 'German')",
			expected: ASTExpectation{
				NodeType: "ConstructorReference", // 注意实际是ConstructorReference不是ConstructorExpression
				Value:    "new org.spring.samples.spel.inventor.Inventor('Albert Einstein', 'German')",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "org.spring.samples.spel.inventor.Inventor",
						Children: []ASTExpectation{
							{NodeType: "Identifier", Value: "org", Children: []ASTExpectation{}},
							{NodeType: "Identifier", Value: "spring", Children: []ASTExpectation{}},
							{NodeType: "Identifier", Value: "samples", Children: []ASTExpectation{}},
							{NodeType: "Identifier", Value: "spel", Children: []ASTExpectation{}},
							{NodeType: "Identifier", Value: "inventor", Children: []ASTExpectation{}},
							{NodeType: "Identifier", Value: "Inventor", Children: []ASTExpectation{}},
						},
					},
					{NodeType: "StringLiteral", Value: "'Albert Einstein'", Children: []ASTExpectation{}},
					{NodeType: "StringLiteral", Value: "'German'", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "赋值变量引用",
			expression: "Name = #newName",
			expected: ASTExpectation{
				NodeType: "Assign",
				Value:    "Name = #newName",
				Children: []ASTExpectation{
					{NodeType: "PropertyOrFieldReference", Value: "Name", Children: []ASTExpectation{}},
					{NodeType: "VariableReference", Value: "#newName", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "选择表达式",
			expression: "#primes.?[#this>10]",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "#primes.?[(#this > 10)]", // 注意格式化包含括号和空格
				Children: []ASTExpectation{
					{NodeType: "VariableReference", Value: "#primes", Children: []ASTExpectation{}},
					{
						NodeType: "Selection",
						Value:    ".?[(#this > 10)]",
						Children: []ASTExpectation{
							{
								NodeType: "OpGT",
								Value:    "(#this > 10)",
								Children: []ASTExpectation{
									{NodeType: "VariableReference", Value: "#this", Children: []ASTExpectation{}},
									{NodeType: "IntLiteral", Value: "10", Children: []ASTExpectation{}},
								},
							},
						},
					},
				},
			},
		},
		{
			name:       "函数引用",
			expression: "#reverseString('hello')",
			expected: ASTExpectation{
				NodeType: "FunctionReference",
				Value:    "#reverseString('hello')",
				Children: []ASTExpectation{
					{NodeType: "StringLiteral", Value: "'hello'", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "Bean引用",
			expression: "@something",
			expected: ASTExpectation{
				NodeType: "BeanReference",
				Value:    "@something",
				Children: []ASTExpectation{},
			},
		},
		{
			name:       "工厂Bean引用",
			expression: "&foo",
			expected: ASTExpectation{
				NodeType: "BeanReference", // 注意&foo实际被解析为BeanReference
				Value:    "@foo",          // 注意实际显示为@foo
				Children: []ASTExpectation{},
			},
		},
		{
			name:       "安全导航属性",
			expression: "PlaceOfBirth?.City",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "PlaceOfBirth?.City",
				Children: []ASTExpectation{
					{NodeType: "PropertyOrFieldReference", Value: "PlaceOfBirth", Children: []ASTExpectation{}},
					{NodeType: "PropertyOrFieldReference", Value: "?.City", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "映射选择",
			expression: "map.?[value<27]",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "map.?[(value < 27)]", // 注意格式化包含括号和空格
				Children: []ASTExpectation{
					{NodeType: "PropertyOrFieldReference", Value: "map", Children: []ASTExpectation{}},
					{
						NodeType: "Selection",
						Value:    ".?[(value < 27)]",
						Children: []ASTExpectation{
							{
								NodeType: "OpLT",
								Value:    "(value < 27)",
								Children: []ASTExpectation{
									{NodeType: "PropertyOrFieldReference", Value: "value", Children: []ASTExpectation{}},
									{NodeType: "IntLiteral", Value: "27", Children: []ASTExpectation{}},
								},
							},
						},
					},
				},
			},
		},
		{
			name:       "投影表达式",
			expression: "Members.![placeOfBirth.city]",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "Members.![placeOfBirth.city]",
				Children: []ASTExpectation{
					{NodeType: "PropertyOrFieldReference", Value: "Members", Children: []ASTExpectation{}},
					{
						NodeType: "Projection",
						Value:    ".![placeOfBirth.city]",
						Children: []ASTExpectation{
							{
								NodeType: "CompoundExpression",
								Value:    "placeOfBirth.city",
								Children: []ASTExpectation{
									{NodeType: "PropertyOrFieldReference", Value: "placeOfBirth", Children: []ASTExpectation{}},
									{NodeType: "PropertyOrFieldReference", Value: ".city", Children: []ASTExpectation{}},
								},
							},
						},
					},
				},
			},
		},
		{
			name:       "三元运算符",
			expression: "falseValue ? 'trueExp' : 'falseExp'",
			expected: ASTExpectation{
				NodeType: "Ternary",
				Value:    "falseValue ? 'trueExp' : 'falseExp'", // 实际输出包含空格
				Children: []ASTExpectation{
					{NodeType: "PropertyOrFieldReference", Value: "falseValue", Children: []ASTExpectation{}},
					{NodeType: "StringLiteral", Value: "'trueExp'", Children: []ASTExpectation{}},
					{NodeType: "StringLiteral", Value: "'falseExp'", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "Elvis运算符",
			expression: "name?:'Unknown'",
			expected: ASTExpectation{
				NodeType: "Elvis",
				Value:    "name ?: 'Unknown'", // 实际输出包含空格
				Children: []ASTExpectation{
					{NodeType: "PropertyOrFieldReference", Value: "name", Children: []ASTExpectation{}},
					{NodeType: "StringLiteral", Value: "'Unknown'", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "算术表达式加法",
			expression: "1 + 1",
			expected: ASTExpectation{
				NodeType: "OpPlus",
				Value:    "(1 + 1)",
				Children: []ASTExpectation{
					{NodeType: "IntLiteral", Value: "1", Children: []ASTExpectation{}},
					{NodeType: "IntLiteral", Value: "1", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "字符串连接",
			expression: "'Hello' + ' World'",
			expected: ASTExpectation{
				NodeType: "OpPlus",
				Value:    "('Hello' + ' World')",
				Children: []ASTExpectation{
					{NodeType: "StringLiteral", Value: "'Hello'", Children: []ASTExpectation{}},
					{NodeType: "StringLiteral", Value: "' World'", Children: []ASTExpectation{}},
				},
			},
		},
	}

	// 执行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- 测试用例: %s ---\n", tc.name)
			fmt.Printf("表达式: %s\n", tc.expression)

			// 解析表达式
			config := ast.NewSpelParserConfiguration()
			parser := ast.NewInternalSpelExpressionParser(config)

			spelExpr, err := parser.DoParseExpression(tc.expression)
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
