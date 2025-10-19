package main

import (
	"fmt"
	"github.com/weaweawe01/ParserSpel/ast"
	"math/big"
	"testing"
)

// TypeComparator 接口定义类型比较器
type TypeComparator interface {
	CanCompare(firstObject, secondObject interface{}) bool
	Compare(firstObject, secondObject interface{}) (int, error)
}

// StandardTypeComparator 标准类型比较器实现
type StandardTypeComparator struct{}

// NewStandardTypeComparator 创建标准类型比较器
func NewStandardTypeComparator() *StandardTypeComparator {
	return &StandardTypeComparator{}
}

// CanCompare 检查是否可以比较两个对象
func (c *StandardTypeComparator) CanCompare(firstObject, secondObject interface{}) bool {
	// 如果任一对象为 nil，可以比较
	if firstObject == nil || secondObject == nil {
		return true
	}

	// 获取类型信息
	type1 := fmt.Sprintf("%T", firstObject)
	type2 := fmt.Sprintf("%T", secondObject)

	// 数字类型之间可以比较
	numericTypes := map[string]bool{
		"int":     true,
		"int8":    true,
		"int16":   true,
		"int32":   true,
		"int64":   true,
		"uint":    true,
		"uint8":   true,
		"uint16":  true,
		"uint32":  true,
		"uint64":  true,
		"float32": true,
		"float64": true,
	}

	isFirstNumeric := numericTypes[type1]
	isSecondNumeric := numericTypes[type2]

	if isFirstNumeric && isSecondNumeric {
		return true
	}

	// BigDecimal 与数字类型可以比较
	if (type1 == "*big.Float" || type1 == "*big.Int") && isSecondNumeric {
		return true
	}
	if isFirstNumeric && (type2 == "*big.Float" || type2 == "*big.Int") {
		return true
	}
	if (type1 == "*big.Float" || type1 == "*big.Int") && (type2 == "*big.Float" || type2 == "*big.Int") {
		return true
	}

	// 相同类型可以比较
	if type1 == type2 {
		return true
	}

	// 字符串和数字不能比较（除非自定义比较器）
	if type1 == "string" && isSecondNumeric {
		return false
	}
	if isFirstNumeric && type2 == "string" {
		return false
	}

	return false
}

// Compare 比较两个对象
func (c *StandardTypeComparator) Compare(firstObject, secondObject interface{}) (int, error) {
	// 处理 nil 值
	if firstObject == nil && secondObject == nil {
		return 0, nil
	}
	if firstObject == nil {
		return -1, nil
	}
	if secondObject == nil {
		return 1, nil
	}

	// 转换为浮点数进行比较（数字类型）
	first, err1 := toFloat64(firstObject)
	second, err2 := toFloat64(secondObject)

	if err1 == nil && err2 == nil {
		if first < second {
			return -1, nil
		} else if first > second {
			return 1, nil
		} else {
			return 0, nil
		}
	}

	// 字符串比较
	if str1, ok1 := firstObject.(string); ok1 {
		if str2, ok2 := secondObject.(string); ok2 {
			if str1 < str2 {
				return -1, nil
			} else if str1 > str2 {
				return 1, nil
			} else {
				return 0, nil
			}
		}
	}

	return 0, fmt.Errorf("无法比较类型 %T 和 %T", firstObject, secondObject)
}

// toFloat64 将各种数字类型转换为 float64
func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case *big.Float:
		result, _ := v.Float64()
		return result, nil
	case *big.Int:
		return float64(v.Int64()), nil
	default:
		return 0, fmt.Errorf("无法转换类型 %T 为 float64", value)
	}
}

// CustomTypeComparator 自定义比较器，声明所有东西都相等
type CustomTypeComparator struct{}

func (c *CustomTypeComparator) CanCompare(firstObject, secondObject interface{}) bool {
	return true
}

func (c *CustomTypeComparator) Compare(firstObject, secondObject interface{}) (int, error) {
	return 0, nil // 所有东西都相等
}

// TestComparatorPrimitives 测试基本类型比较
func TestComparatorPrimitives(t *testing.T) {
	fmt.Println("=== 基本类型比较测试 ===")

	comparator := NewStandardTypeComparator()

	testCases := []struct {
		name     string
		first    interface{}
		second   interface{}
		expected int
		desc     string
	}{
		// 整数比较
		{"int 1 < 2", 1, 2, -1, "整数1小于2"},
		{"int 1 == 1", 1, 1, 0, "整数1等于1"},
		{"int 2 > 1", 2, 1, 1, "整数2大于1"},

		// 浮点数与整数比较
		{"double 1.0 < 2", 1.0, 2, -1, "浮点数1.0小于整数2"},
		{"double 1.0 == 1", 1.0, 1, 0, "浮点数1.0等于整数1"},
		{"double 2.0 > 1", 2.0, 1, 1, "浮点数2.0大于整数1"},

		// float32 与整数比较
		{"float32 1.0 < 2", float32(1.0), 2, -1, "float32 1.0小于整数2"},
		{"float32 1.0 == 1", float32(1.0), 1, 0, "float32 1.0等于整数1"},
		{"float32 2.0 > 1", float32(2.0), 1, 1, "float32 2.0大于整数1"},

		// int64 与整数比较
		{"int64 1 < 2", int64(1), 2, -1, "int64 1小于整数2"},
		{"int64 1 == 1", int64(1), 1, 0, "int64 1等于整数1"},
		{"int64 2 > 1", int64(2), 1, 1, "int64 2大于整数1"},

		// 整数与 int64 比较
		{"int 1 < int64 2", 1, int64(2), -1, "整数1小于int64 2"},
		{"int 1 == int64 1", 1, int64(1), 0, "整数1等于int64 1"},
		{"int 2 > int64 1", 2, int64(1), 1, "整数2大于int64 1"},

		// int64 之间比较
		{"int64 1 < int64 2", int64(1), int64(2), -1, "int64 1小于int64 2"},
		{"int64 1 == int64 1", int64(1), int64(1), 0, "int64 1等于int64 1"},
		{"int64 2 > int64 1", int64(2), int64(1), 1, "int64 2大于int64 1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- %s ---\n", tc.name)
			fmt.Printf("描述: %s\n", tc.desc)
			fmt.Printf("比较: %v 与 %v\n", tc.first, tc.second)

			result, err := comparator.Compare(tc.first, tc.second)
			if err != nil {
				t.Fatalf("比较失败: %v", err)
			}

			fmt.Printf("期望结果: %d\n", tc.expected)
			fmt.Printf("实际结果: %d\n", result)

			// 检查符号是否正确
			if (tc.expected < 0 && result >= 0) ||
				(tc.expected == 0 && result != 0) ||
				(tc.expected > 0 && result <= 0) {
				t.Errorf("比较结果不正确: 期望 %d, 实际 %d", tc.expected, result)
			} else {
				fmt.Println("✅ 比较结果正确")
			}
		})
	}
}

// TestComparatorNonPrimitiveNumbers 测试非基本数字类型
func TestComparatorNonPrimitiveNumbers(t *testing.T) {
	fmt.Println("\n=== 非基本数字类型比较测试 ===")

	comparator := NewStandardTypeComparator()

	// 创建 BigDecimal 等价的 big.Float
	bdOne, _ := new(big.Float).SetString("1")
	bdTwo, _ := new(big.Float).SetString("2")
	bdOneAnother, _ := new(big.Float).SetString("1")

	testCases := []struct {
		name     string
		first    interface{}
		second   interface{}
		expected int
		desc     string
	}{
		// BigDecimal 之间比较
		{"BigFloat 1 < 2", bdOne, bdTwo, -1, "BigFloat 1小于2"},
		{"BigFloat 1 == 1", bdOne, bdOneAnother, 0, "BigFloat 1等于1"},
		{"BigFloat 2 > 1", bdTwo, bdOne, 1, "BigFloat 2大于1"},

		// 整数与 BigDecimal 比较
		{"int 1 < BigFloat 2", 1, bdTwo, -1, "整数1小于BigFloat 2"},
		{"int 1 == BigFloat 1", 1, bdOne, 0, "整数1等于BigFloat 1"},
		{"int 2 > BigFloat 1", 2, bdOne, 1, "整数2大于BigFloat 1"},

		// 浮点数与 BigDecimal 比较
		{"double 1.0 < BigFloat 2", 1.0, bdTwo, -1, "浮点数1.0小于BigFloat 2"},
		{"double 1.0 == BigFloat 1", 1.0, bdOne, 0, "浮点数1.0等于BigFloat 1"},
		{"double 2.0 > BigFloat 1", 2.0, bdOne, 1, "浮点数2.0大于BigFloat 1"},

		// float32 与 BigDecimal 比较
		{"float32 1.0 < BigFloat 2", float32(1.0), bdTwo, -1, "float32 1.0小于BigFloat 2"},
		{"float32 1.0 == BigFloat 1", float32(1.0), bdOne, 0, "float32 1.0等于BigFloat 1"},
		{"float32 2.0 > BigFloat 1", float32(2.0), bdOne, 1, "float32 2.0大于BigFloat 1"},

		// int64 与 BigDecimal 比较
		{"int64 1 < BigFloat 2", int64(1), bdTwo, -1, "int64 1小于BigFloat 2"},
		{"int64 1 == BigFloat 1", int64(1), bdOne, 0, "int64 1等于BigFloat 1"},
		{"int64 2 > BigFloat 1", int64(2), bdOne, 1, "int64 2大于BigFloat 1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- %s ---\n", tc.name)
			fmt.Printf("描述: %s\n", tc.desc)
			fmt.Printf("比较: %v 与 %v\n", tc.first, tc.second)

			result, err := comparator.Compare(tc.first, tc.second)
			if err != nil {
				t.Fatalf("比较失败: %v", err)
			}

			fmt.Printf("期望结果: %d\n", tc.expected)
			fmt.Printf("实际结果: %d\n", result)

			// 检查符号是否正确
			if (tc.expected < 0 && result >= 0) ||
				(tc.expected == 0 && result != 0) ||
				(tc.expected > 0 && result <= 0) {
				t.Errorf("比较结果不正确: 期望 %d, 实际 %d", tc.expected, result)
			} else {
				fmt.Println("✅ 比较结果正确")
			}
		})
	}
}

// TestComparatorNulls 测试 nil 值比较
func TestComparatorNulls(t *testing.T) {
	fmt.Println("\n=== nil 值比较测试 ===")

	comparator := NewStandardTypeComparator()

	testCases := []struct {
		name     string
		first    interface{}
		second   interface{}
		expected int
		desc     string
	}{
		{"nil < string", nil, "abc", -1, "nil小于字符串"},
		{"nil == nil", nil, nil, 0, "nil等于nil"},
		{"string > nil", "abc", nil, 1, "字符串大于nil"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- %s ---\n", tc.name)
			fmt.Printf("描述: %s\n", tc.desc)
			fmt.Printf("比较: %v 与 %v\n", tc.first, tc.second)

			result, err := comparator.Compare(tc.first, tc.second)
			if err != nil {
				t.Fatalf("比较失败: %v", err)
			}

			fmt.Printf("期望结果: %d\n", tc.expected)
			fmt.Printf("实际结果: %d\n", result)

			// 检查符号是否正确
			if (tc.expected < 0 && result >= 0) ||
				(tc.expected == 0 && result != 0) ||
				(tc.expected > 0 && result <= 0) {
				t.Errorf("比较结果不正确: 期望 %d, 实际 %d", tc.expected, result)
			} else {
				fmt.Println("✅ 比较结果正确")
			}
		})
	}
}

// TestComparatorObjects 测试对象比较
func TestComparatorObjects(t *testing.T) {
	fmt.Println("\n=== 对象比较测试 ===")

	comparator := NewStandardTypeComparator()

	testCases := []struct {
		name     string
		first    interface{}
		second   interface{}
		expected int
		desc     string
	}{
		{"string a == a", "a", "a", 0, "字符串a等于a"},
		{"string a < b", "a", "b", -1, "字符串a小于b"},
		{"string b > a", "b", "a", 1, "字符串b大于a"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- %s ---\n", tc.name)
			fmt.Printf("描述: %s\n", tc.desc)
			fmt.Printf("比较: %v 与 %v\n", tc.first, tc.second)

			result, err := comparator.Compare(tc.first, tc.second)
			if err != nil {
				t.Fatalf("比较失败: %v", err)
			}

			fmt.Printf("期望结果: %d\n", tc.expected)
			fmt.Printf("实际结果: %d\n", result)

			// 检查符号是否正确
			if (tc.expected < 0 && result >= 0) ||
				(tc.expected == 0 && result != 0) ||
				(tc.expected > 0 && result <= 0) {
				t.Errorf("比较结果不正确: 期望 %d, 实际 %d", tc.expected, result)
			} else {
				fmt.Println("✅ 比较结果正确")
			}
		})
	}
}

// TestCanCompare 测试能否比较
func TestCanCompare(t *testing.T) {
	fmt.Println("\n=== 能否比较测试 ===")

	comparator := NewStandardTypeComparator()

	testCases := []struct {
		name     string
		first    interface{}
		second   interface{}
		expected bool
		desc     string
	}{
		{"nil 和 int", nil, 1, true, "nil和整数可以比较"},
		{"int 和 nil", 1, nil, true, "整数和nil可以比较"},
		{"int 和 int", 2, 1, true, "整数和整数可以比较"},
		{"string 和 string", "abc", "def", true, "字符串和字符串可以比较"},
		{"string 和 int", "abc", 3, false, "字符串和整数不能比较"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- %s ---\n", tc.name)
			fmt.Printf("描述: %s\n", tc.desc)
			fmt.Printf("检查: %v 与 %v\n", tc.first, tc.second)

			result := comparator.CanCompare(tc.first, tc.second)

			fmt.Printf("期望结果: %t\n", tc.expected)
			fmt.Printf("实际结果: %t\n", result)

			if result != tc.expected {
				t.Errorf("CanCompare结果不正确: 期望 %t, 实际 %t", tc.expected, result)
			} else {
				fmt.Println("✅ CanCompare结果正确")
			}
		})
	}
}

// TestCustomComparatorWorksWithEquality 测试自定义比较器与相等性
func TestCustomComparatorWorksWithEquality(t *testing.T) {
	fmt.Println("\n=== 自定义比较器相等性测试 ===")

	// 模拟 SpEL 表达式解析和求值
	testCases := []struct {
		name       string
		expression string
		expected   ASTExpectation
	}{
		{
			name:       "字符串和数字相等比较",
			expression: "'1' == 1",
			expected: ASTExpectation{
				NodeType: "OpEQ",
				Value:    "('1' == 1)",
				Children: []ASTExpectation{
					{NodeType: "StringLiteral", Value: "'1'", Children: []ASTExpectation{}},
					{NodeType: "IntLiteral", Value: "1", Children: []ASTExpectation{}},
				},
			},
		},
	}

	parser := ast.NewSpelExpressionParser()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- %s ---\n", tc.name)
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

			// 注意：在有完整上下文支持后，可以测试自定义比较器
			// customComparator := &CustomTypeComparator{}
			// context := NewEvaluationContextWithComparator(customComparator)
			// result, err := expr.GetValueWithContext(context)
			// if err != nil {
			//     t.Fatalf("求值失败: %v", err)
			// }
			// if result != true {
			//     t.Errorf("自定义比较器测试失败: 期望 true, 实际 %v", result)
			// }

			fmt.Println("✅ 表达式解析成功")
		})
	}
}

// TestComparisonOperators 测试比较运算符
func TestComparisonOperators(t *testing.T) {
	fmt.Println("\n=== 比较运算符测试 ===")

	testCases := []struct {
		name       string
		expression string
		expected   ASTExpectation
	}{
		{
			name:       "等于运算符",
			expression: "1 == 1",
			expected: ASTExpectation{
				NodeType: "OpEQ",
				Value:    "(1 == 1)",
				Children: []ASTExpectation{
					{NodeType: "IntLiteral", Value: "1", Children: []ASTExpectation{}},
					{NodeType: "IntLiteral", Value: "1", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "不等于运算符",
			expression: "1 != 2",
			expected: ASTExpectation{
				NodeType: "OpNE",
				Value:    "(1 != 2)",
				Children: []ASTExpectation{
					{NodeType: "IntLiteral", Value: "1", Children: []ASTExpectation{}},
					{NodeType: "IntLiteral", Value: "2", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "小于运算符",
			expression: "1 < 2",
			expected: ASTExpectation{
				NodeType: "OpLT",
				Value:    "(1 < 2)",
				Children: []ASTExpectation{
					{NodeType: "IntLiteral", Value: "1", Children: []ASTExpectation{}},
					{NodeType: "IntLiteral", Value: "2", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "小于等于运算符",
			expression: "1 <= 2",
			expected: ASTExpectation{
				NodeType: "OpLE",
				Value:    "(1 <= 2)",
				Children: []ASTExpectation{
					{NodeType: "IntLiteral", Value: "1", Children: []ASTExpectation{}},
					{NodeType: "IntLiteral", Value: "2", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "大于运算符",
			expression: "2 > 1",
			expected: ASTExpectation{
				NodeType: "OpGT",
				Value:    "(2 > 1)",
				Children: []ASTExpectation{
					{NodeType: "IntLiteral", Value: "2", Children: []ASTExpectation{}},
					{NodeType: "IntLiteral", Value: "1", Children: []ASTExpectation{}},
				},
			},
		},
		{
			name:       "大于等于运算符",
			expression: "2 >= 1",
			expected: ASTExpectation{
				NodeType: "OpGE",
				Value:    "(2 >= 1)",
				Children: []ASTExpectation{
					{NodeType: "IntLiteral", Value: "2", Children: []ASTExpectation{}},
					{NodeType: "IntLiteral", Value: "1", Children: []ASTExpectation{}},
				},
			},
		},
	}

	parser := ast.NewSpelExpressionParser()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n--- %s ---\n", tc.name)
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
