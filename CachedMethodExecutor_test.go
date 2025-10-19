package main

// 参考 https://github.com/spring-projects/spring-framework/blob/main/spring-expression/src/test/java/org/springframework/expression/spel/CachedMethodExecutorTests.java

import (
	"fmt"
	"github.com/weaweawe01/ParserSpel/ast"
	"testing"
)

// BaseObject represents the base class with echo method for strings
type BaseObject struct{}

func (b *BaseObject) Echo(value string) string {
	return "String: " + value
}

// RootObject extends BaseObject and adds echo method for integers
type RootObject struct {
	*BaseObject
}

func (r *RootObject) Echo(value int) string {
	return fmt.Sprintf("int: %d", value)
}

// TestCachedMethodExecutor 测试方法执行器的缓存机制
func TestCachedMethodExecutor(t *testing.T) {
	fmt.Println("=== 缓存方法执行器测试 ===")

	// 测试用例定义
	testCases := []struct {
		name       string
		expression string
		expected   ASTExpectation
	}{
		{
			name:       "方法调用 - echo(#var)",
			expression: "echo(#var)",
			expected: ASTExpectation{
				NodeType: "MethodReference",
				Value:    "echo(#var)",
				Children: []ASTExpectation{
					{
						NodeType: "VariableReference",
						Value:    "#var",
						Children: []ASTExpectation{},
					},
				},
			},
		},
		{
			name:       "目标方法调用 - #var.echo(42)",
			expression: "#var.echo(42)",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "#var.echo(42)",
				Children: []ASTExpectation{
					{
						NodeType: "VariableReference",
						Value:    "#var",
						Children: []ASTExpectation{},
					},
					{
						NodeType: "MethodReference",
						Value:    "echo(42)",
						Children: []ASTExpectation{
							{
								NodeType: "IntLiteral",
								Value:    "42",
								Children: []ASTExpectation{},
							},
						},
					},
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

// TestCachedExecutionForParameters 测试参数缓存执行
func TestCachedExecutionForParameters(t *testing.T) {
	fmt.Println("\n=== 参数缓存执行测试 ===")

	// 模拟测试用例
	testCases := []struct {
		name        string
		expression  string
		variable    interface{}
		expected    string
		description string
	}{
		{
			name:        "整数参数 - 第一次",
			expression:  "echo(#var)",
			variable:    42,
			expected:    "int: 42",
			description: "使用整数参数调用echo方法",
		},
		{
			name:        "整数参数 - 第二次（缓存）",
			expression:  "echo(#var)",
			variable:    42,
			expected:    "int: 42",
			description: "重复使用相同整数参数，应该使用缓存",
		},
		{
			name:        "字符串参数",
			expression:  "echo(#var)",
			variable:    "Deep Thought",
			expected:    "String: Deep Thought",
			description: "使用字符串参数调用echo方法",
		},
		{
			name:        "整数参数 - 第三次（缓存）",
			expression:  "echo(#var)",
			variable:    42,
			expected:    "int: 42",
			description: "再次使用整数参数，应该使用缓存",
		},
	}

	parser := ast.NewSpelExpressionParser()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- %s ---\n", tc.name)
			fmt.Printf("表达式: %s\n", tc.expression)
			fmt.Printf("变量值: %v\n", tc.variable)
			fmt.Printf("描述: %s\n", tc.description)

			expr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				t.Fatalf("解析表达式失败: %v", err)
			}

			// 注意：在实际实现中，这里需要设置变量和上下文
			// 由于当前的 Go 实现可能还没有完整的变量上下文支持，
			// 这里主要验证表达式解析的正确性

			fmt.Printf("解析结果: %s\n", expr.GetExpressionString())
			fmt.Printf("期望结果: %s\n", tc.expected)

			// 在有完整上下文支持后，可以添加实际的求值测试
			// result, err := expr.GetValueWithContext(context)
			// if err != nil {
			//     t.Fatalf("求值失败: %v", err)
			// }
			// if result != tc.expected {
			//     t.Errorf("结果不匹配: 期望 %v, 实际 %v", tc.expected, result)
			// }

			fmt.Println("✅ 表达式解析成功")
		})
	}
}

// TestCachedExecutionForTarget 测试目标对象缓存执行
func TestCachedExecutionForTarget(t *testing.T) {
	fmt.Println("\n=== 目标对象缓存执行测试 ===")

	// 模拟测试用例
	testCases := []struct {
		name        string
		expression  string
		target      interface{}
		expected    string
		description string
	}{
		{
			name:        "RootObject目标 - 第一次",
			expression:  "#var.echo(42)",
			target:      &RootObject{BaseObject: &BaseObject{}},
			expected:    "int: 42",
			description: "使用RootObject调用echo方法",
		},
		{
			name:        "RootObject目标 - 第二次（缓存）",
			expression:  "#var.echo(42)",
			target:      &RootObject{BaseObject: &BaseObject{}},
			expected:    "int: 42",
			description: "重复使用RootObject，应该使用缓存",
		},
		{
			name:        "BaseObject目标",
			expression:  "#var.echo(42)",
			target:      &BaseObject{},
			expected:    "String: 42",
			description: "使用BaseObject调用echo方法",
		},
		{
			name:        "RootObject目标 - 第三次（缓存）",
			expression:  "#var.echo(42)",
			target:      &RootObject{BaseObject: &BaseObject{}},
			expected:    "int: 42",
			description: "再次使用RootObject，应该使用缓存",
		},
	}

	parser := ast.NewSpelExpressionParser()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- %s ---\n", tc.name)
			fmt.Printf("表达式: %s\n", tc.expression)
			fmt.Printf("目标类型: %T\n", tc.target)
			fmt.Printf("描述: %s\n", tc.description)

			expr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				t.Fatalf("解析表达式失败: %v", err)
			}

			// 注意：在实际实现中，这里需要设置目标对象和上下文
			// 由于当前的 Go 实现可能还没有完整的对象上下文支持，
			// 这里主要验证表达式解析的正确性

			fmt.Printf("解析结果: %s\n", expr.GetExpressionString())
			fmt.Printf("期望结果: %s\n", tc.expected)

			// 在有完整上下文支持后，可以添加实际的求值测试
			// context.SetVariable("var", tc.target)
			// result, err := expr.GetValueWithContext(context)
			// if err != nil {
			//     t.Fatalf("求值失败: %v", err)
			// }
			// if result != tc.expected {
			//     t.Errorf("结果不匹配: 期望 %v, 实际 %v", tc.expected, result)
			// }

			fmt.Println("✅ 表达式解析成功")
		})
	}
}

// TestMethodCachingBehavior 测试方法缓存行为
func TestMethodCachingBehavior(t *testing.T) {
	fmt.Println("\n=== 方法缓存行为测试 ===")

	testCases := []struct {
		name       string
		expression string
		expected   ASTExpectation
	}{
		{
			name:       "简单方法调用",
			expression: "toString()",
			expected: ASTExpectation{
				NodeType: "MethodReference",
				Value:    "toString()",
				Children: []ASTExpectation{},
			},
		},
		{
			name:       "带参数方法调用",
			expression: "substring(1, 5)",
			expected: ASTExpectation{
				NodeType: "MethodReference",
				Value:    "substring(1, 5)",
				Children: []ASTExpectation{
					{NodeType: "IntLiteral", Value: "1", Children: []ASTExpectation{}},
					{NodeType: "IntLiteral", Value: "5", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "链式方法调用",
			expression: "getValue().toString()",
			expected: ASTExpectation{
				NodeType: "CompoundExpression",
				Value:    "getValue().toString()",
				Children: []ASTExpectation{
					{
						NodeType: "MethodReference",
						Value:    "getValue()",
						Children: []ASTExpectation{},
					},
					{
						NodeType: "MethodReference",
						Value:    "toString()",
						Children: []ASTExpectation{},
					},
				},
			},
		},
	}

	parser := ast.NewSpelExpressionParser()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- 测试用例: %s ---\n", tc.name)
			fmt.Printf("表达式: %s\n", tc.expression)

			expr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				t.Fatalf("解析表达式失败: %v", err)
			}

			fmt.Printf("AST-> %v\n", expr.AST)
			fmt.Printf("%s\n", expr.GetExpressionString())

			fmt.Println("实际AST结构:")
			ast.PrintAST(expr.AST, 0)

			fmt.Println("验证AST结构...")
			if !validateASTStructure(expr.AST, tc.expected) {
				t.Errorf("AST结构验证失败")
			} else {
				fmt.Println("✅ AST结构验证通过")
			}
		})
	}
}

// TestMethodOverloading 测试方法重载
func TestMethodOverloading(t *testing.T) {
	fmt.Println("\n=== 方法重载测试 ===")

	// 模拟方法重载测试
	overloadTests := []struct {
		name        string
		expression  string
		description string
	}{
		{
			name:        "重载方法 - 整数参数",
			expression:  "echo(42)",
			description: "调用接受整数参数的echo方法",
		},
		{
			name:        "重载方法 - 字符串参数",
			expression:  "echo('hello')",
			description: "调用接受字符串参数的echo方法",
		},
		{
			name:        "重载方法 - 多个参数",
			expression:  "echo('test', 123)",
			description: "调用接受多个参数的echo方法",
		},
	}

	parser := ast.NewSpelExpressionParser()

	for _, tc := range overloadTests {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- %s ---\n", tc.name)
			fmt.Printf("表达式: %s\n", tc.expression)
			fmt.Printf("描述: %s\n", tc.description)

			expr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				t.Fatalf("解析表达式失败: %v", err)
			}

			fmt.Printf("解析结果: %s\n", expr.GetExpressionString())
			fmt.Println("实际AST结构:")
			ast.PrintAST(expr.AST, 0)

			// 验证方法调用是否正确解析
			if expr.AST == nil {
				t.Error("AST不应该为空")
			} else {
				fmt.Println("✅ 方法重载表达式解析成功")
			}
		})
	}
}
