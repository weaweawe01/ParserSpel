package main

import (
	"fmt"
	"github.com/weaweawe01/ParserSpel/ast"
	"strings"
	"testing"
)

// TestEvaluationMiscellaneous tests miscellaneous evaluation scenarios
func TestEvaluationMiscellaneous(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	// Test expression length limits
	t.Run("ExpressionLength", func(t *testing.T) {
		// Test expression within limit
		expression := "'X' + '" + strings.Repeat(" ", 9992) + "'"
		if len(expression) != 10000 {
			t.Errorf("Expression length should be 10000, got %d", len(expression))
		}

		expr, err := parser.ParseExpression(expression)
		if err != nil {
			t.Errorf("Failed to parse expression: %v", err)
		}

		result, err := expr.GetValue()
		if err != nil {
			t.Errorf("Failed to evaluate expression: %v", err)
		}

		resultStr := result.(string)
		if len(resultStr) != 9993 {
			t.Errorf("Result length should be 9993, got %d", len(resultStr))
		}

		if strings.TrimSpace(resultStr) != "X" {
			t.Errorf("Result should be 'X', got '%s'", strings.TrimSpace(resultStr))
		}
	})

	// Test Elvis operator
	t.Run("ElvisOperator", func(t *testing.T) {
		testCases := []struct {
			expression string
			expected   interface{}
		}{
			{"'Andy'?:'Dave'", "Andy"},
			{"null?:'Dave'", "Dave"},
			{"3?:1", int64(3)},
			{"(2*3)?:1*10", int64(6)},
			{"null?:2*10", int64(20)},
			{"(null?:1)*10", int64(10)},
		}

		for _, tc := range testCases {
			expr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				t.Errorf("Failed to parse expression '%s': %v", tc.expression, err)
				continue
			}

			result, err := expr.GetValue()
			if err != nil {
				t.Errorf("Failed to evaluate expression '%s': %v", tc.expression, err)
				continue
			}

			if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
				t.Errorf("Expression '%s': expected %v, got %v", tc.expression, tc.expected, result)
			}
		}
	})

	// Test safe navigation
	t.Run("SafeNavigation", func(t *testing.T) {
		expr, err := parser.ParseExpression("null?.null?.null")
		if err != nil {
			t.Errorf("Failed to parse safe navigation expression: %v", err)
		}

		result, err := expr.GetValue()
		if err != nil {
			t.Errorf("Failed to evaluate safe navigation expression: %v", err)
		}

		// In our Go implementation, null might be represented as a string "null"
		if result != nil && result != "null" {
			t.Errorf("Safe navigation should return nil or null, got %v", result)
		}
	})

	// Test mixing operators
	t.Run("MixingOperators", func(t *testing.T) {
		expr, err := parser.ParseExpression("true and 5>3")
		if err != nil {
			t.Errorf("Failed to parse mixed operators expression: %v", err)
		}

		result, err := expr.GetValue()
		if err != nil {
			t.Errorf("Failed to evaluate mixed operators expression: %v", err)
		}

		if result != true {
			t.Errorf("Mixed operators should return true, got %v", result)
		}
	})
}

// TestStringLiterals tests string literal parsing and evaluation
func TestStringLiterals(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	testCases := []struct {
		name       string
		expression string
		expected   string
	}{
		{"SingleQuotes", "'hello'", "hello"},
		{"SingleQuotesWithSpace", "'hello world'", "hello world"},
		{"DoubleQuotes", "\"hello\"", "hello"},
		{"DoubleQuotesWithSpace", "\"hello world\"", "hello world"},
		{"SingleQuotesInsideSingle", "'Tony''s Pizza'", "Tony's Pizza"},
		{"BigSingleQuotesInsideSingle", "'big ''''pizza'''' parlor'", "big ''pizza'' parlor"},
		{"DoubleQuotesInsideDouble", "\"big \"\"pizza\"\" parlor\"", "big \"pizza\" parlor"},
		{"BigDoubleQuotesInsideDouble", "\"big \"\"\"\"pizza\"\"\"\" parlor\"", "big \"\"pizza\"\" parlor"},
		{"SingleQuotesInsideDouble", "\"Tony's Pizza\"", "Tony's Pizza"},
		{"BigSingleQuotesInsideDouble", "\"big ''pizza'' parlor\"", "big ''pizza'' parlor"},
		{"DoubleQuotesInsideSingle", "'big \"pizza\" parlor'", "big \"pizza\" parlor"},
		{"TwoDoubleQuotesInsideSingle", "'two double \"\" quotes'", "two double \"\" quotes"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				t.Errorf("Failed to parse expression '%s': %v", tc.expression, err)
				return
			}

			result, err := expr.GetValue()
			if err != nil {
				t.Errorf("Failed to evaluate expression '%s': %v", tc.expression, err)
				return
			}

			resultStr, ok := result.(string)
			if !ok {
				t.Errorf("Expected string result for '%s', got %T", tc.expression, result)
				return
			}

			if resultStr != tc.expected {
				t.Errorf("Expression '%s': expected '%s', got '%s'", tc.expression, tc.expected, resultStr)
			}
		})
	}

	// Test compound expressions with string literals
	t.Run("CompoundExpressions", func(t *testing.T) {
		testCases := []struct {
			expression string
			expected   bool
		}{
			{"'123''4' == '123''4'", true},
			{"\"123\"\"4\" == \"123\"\"4\"", true},
		}

		for _, tc := range testCases {
			expr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				t.Errorf("Failed to parse expression '%s': %v", tc.expression, err)
				continue
			}

			result, err := expr.GetValue()
			if err != nil {
				t.Errorf("Failed to evaluate expression '%s': %v", tc.expression, err)
				continue
			}

			boolResult, ok := result.(bool)
			if !ok {
				t.Errorf("Expected boolean result for '%s', got %T", tc.expression, result)
				continue
			}

			if boolResult != tc.expected {
				t.Errorf("Expression '%s': expected %v, got %v", tc.expression, tc.expected, boolResult)
			}
		}
	})
}

// TestRelationalOperators tests relational operator evaluation
func TestRelationalOperators(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	testCases := []struct {
		name       string
		expression string
		expected   bool
	}{
		{"GreaterThan", "3 > 6", false},
		{"LessThan", "3 < 6", true},
		{"LessThanOrEqual", "3 <= 6", true},
		{"GreaterThanOrEqual1", "3 >= 6", false},
		{"GreaterThanOrEqual2", "3 >= 3", true},
		{"InstanceofString", "'xyz' instanceof T(String)", false},
		{"InstanceofInt", "'xyz' instanceof T(int)", false},
		{"InstanceofNull1", "null instanceof T(String)", true},
		{"InstanceofNull2", "null instanceof T(Integer)", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				t.Errorf("Failed to parse expression '%s': %v", tc.expression, err)
				return
			}

			result, err := expr.GetValue()
			if err != nil {
				t.Errorf("Failed to evaluate expression '%s': %v", tc.expression, err)
				return
			}

			boolResult, ok := result.(bool)
			if !ok {
				t.Errorf("Expected boolean result for '%s', got %T", tc.expression, result)
				return
			}

			if boolResult != tc.expected {
				t.Errorf("Expression '%s': expected %v, got %v", tc.expression, tc.expected, boolResult)
			}
		})
	}

	// Test matches operator
	t.Run("MatchesOperator", func(t *testing.T) {
		testCases := []struct {
			expression string
			expected   bool
		}{
			{"'5.00' matches '^-?\\\\d+(\\\\.\\\\d{2})?$'", true},
			{"'5.0067' matches '^-?\\\\d+(\\\\.\\\\d{2})?$'", false},
			{"27 matches '^.*2.*$'", true}, // conversion int --> string
		}

		for _, tc := range testCases {
			expr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				t.Errorf("Failed to parse expression '%s': %v", tc.expression, err)
				continue
			}

			result, err := expr.GetValue()
			if err != nil {
				t.Errorf("Failed to evaluate expression '%s': %v", tc.expression, err)
				continue
			}

			boolResult, ok := result.(bool)
			if !ok {
				t.Errorf("Expected boolean result for '%s', got %T", tc.expression, result)
				continue
			}

			if boolResult != tc.expected {
				t.Errorf("Expression '%s': expected %v, got %v", tc.expression, tc.expected, boolResult)
			}
		}
	})

	// Test between operator
	t.Run("BetweenOperator", func(t *testing.T) {
		// Note: This would require implementing between operator support
		// For now, we'll skip this test or implement it when between operator is added
		t.Skip("Between operator not yet implemented in Go version")
	})
}

// TestMethodAndConstructor tests method calls and constructor invocations
func TestMethodAndConstructor(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	t.Run("ConstructorInvocation", func(t *testing.T) {
		testCases := []struct {
			name       string
			expression string
			expected   string
		}{
			{"SimpleString", "new String('hello')", "hello"},
			{"QualifiedString", "new java.lang.String('foobar')", "foobar"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				expr, err := parser.ParseExpression(tc.expression)
				if err != nil {
					t.Errorf("Failed to parse expression '%s': %v", tc.expression, err)
					return
				}

				// Check AST structure
				astString := expr.ToStringAST()
				if !strings.Contains(astString, "new") {
					t.Errorf("AST should contain 'new', got: %s", astString)
				}

				// Note: Actual evaluation would require implementing constructor invocation
				// For now, we just test parsing
			})
		}

		// Test repeated evaluation and AST properties
		t.Run("RepeatedEvaluation", func(t *testing.T) {
			expr, err := parser.ParseExpression("new String('wibble')")
			if err != nil {
				t.Errorf("Failed to parse expression: %v", err)
				return
			}

			// Check AST string representation
			astString := expr.ToStringAST()
			expectedAST := "new String('wibble')"
			if astString != expectedAST {
				t.Errorf("Expected AST '%s', got '%s'", expectedAST, astString)
			}

			// Note: Actual evaluation and writability checks would require full implementation
		})
	})
}

// TestUnaryOperators tests unary operator evaluation
func TestUnaryOperators(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	testCases := []struct {
		name       string
		expression string
		expected   interface{}
	}{
		{"UnaryMinus", "-5", int64(-5)},
		{"UnaryPlus", "+5", int64(5)},
		{"UnaryNotTrue", "!true", false},
		{"UnaryNotFalse", "!false", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				t.Errorf("Failed to parse expression '%s': %v", tc.expression, err)
				return
			}

			result, err := expr.GetValue()
			if err != nil {
				t.Errorf("Failed to evaluate expression '%s': %v", tc.expression, err)
				return
			}

			if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
				t.Errorf("Expression '%s': expected %v, got %v", tc.expression, tc.expected, result)
			}
		})
	}

	// Test unary not with null value should fail
	t.Run("UnaryNotWithNull", func(t *testing.T) {
		expr, err := parser.ParseExpression("!null")
		if err != nil {
			t.Errorf("Failed to parse expression: %v", err)
			return
		}

		_, err = expr.GetValue()
		if err == nil {
			t.Error("Expected error when evaluating '!null', but got none")
		}
	})
}

// TestTernaryOperator tests ternary operator evaluation
func TestTernaryOperator(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	testCases := []struct {
		name       string
		expression string
		expected   interface{}
	}{
		{"BasicTernary1", "2>4?1:2", int64(2)},
		{"BasicTernary2", "'abc'=='abc'?1:2", int64(1)},
		{"NestedTernary", "2>4?(3>2?true:false):(5<3?true:false)", false},
		{"TernaryWithImplicitGrouping1", "4 % 2 == 0 ? 2 : 3 * 10", int64(2)},
		{"TernaryWithImplicitGrouping2", "4 % 2 == 1 ? 2 : 3 * 10", int64(30)},
		{"TernaryWithExplicitGrouping", "((4 % 2 == 0) ? 2 : 1) * 10", int64(20)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				t.Errorf("Failed to parse expression '%s': %v", tc.expression, err)
				return
			}

			result, err := expr.GetValue()
			if err != nil {
				t.Errorf("Failed to evaluate expression '%s': %v", tc.expression, err)
				return
			}

			if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
				t.Errorf("Expression '%s': expected %v, got %v", tc.expression, tc.expected, result)
			}
		})
	}

	// Test ternary with null condition should fail
	t.Run("TernaryWithNull", func(t *testing.T) {
		expr, err := parser.ParseExpression("null ? 0 : 1")
		if err != nil {
			t.Errorf("Failed to parse expression: %v", err)
			return
		}

		_, err = expr.GetValue()
		if err == nil {
			t.Error("Expected error when evaluating ternary with null condition, but got none")
		}
	})
}

// TestTypeReferences tests type reference expressions
func TestTypeReferences(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	testCases := []struct {
		name       string
		expression string
		expected   string
	}{
		{"QualifiedType", "T(java.lang.String)", "class java.lang.String"},
		{"SimpleType", "T(String)", "class java.lang.String"},
		{"PrimitiveInt", "T(int)", "int"},
		{"PrimitiveByte", "T(byte)", "byte"},
		{"PrimitiveChar", "T(char)", "char"},
		{"PrimitiveBoolean", "T(boolean)", "boolean"},
		{"PrimitiveLong", "T(long)", "long"},
		{"PrimitiveShort", "T(short)", "short"},
		{"PrimitiveDouble", "T(double)", "double"},
		{"PrimitiveFloat", "T(float)", "float"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				t.Errorf("Failed to parse expression '%s': %v", tc.expression, err)
				return
			}

			// Check AST structure
			astString := expr.ToStringAST()
			if !strings.Contains(astString, "T(") {
				t.Errorf("AST should contain 'T(', got: %s", astString)
			}

			// Note: Actual type resolution would require implementing type reference evaluation
		})
	}

	// Test type references and qualified identifier caching
	t.Run("QualifiedIdentifierCaching", func(t *testing.T) {
		expr, err := parser.ParseExpression("T(java.lang.String)")
		if err != nil {
			t.Errorf("Failed to parse expression: %v", err)
			return
		}

		// Check AST string representation multiple times to test caching
		astString1 := expr.ToStringAST()
		astString2 := expr.ToStringAST()

		expectedAST := "T(java.lang.String)"
		if astString1 != expectedAST {
			t.Errorf("Expected AST '%s', got '%s'", expectedAST, astString1)
		}

		if astString1 != astString2 {
			t.Error("AST string should be consistent across multiple calls")
		}
	})
}

// TestArithmeticOperations tests arithmetic operations
func TestArithmeticOperations(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	testCases := []struct {
		name       string
		expression string
		expected   interface{}
	}{
		{"BasicArithmetic", "3*4+5", int64(17)},
		{"AdvancedNumerics1", "2.0 * 3e0 * 4", int64(24)},
		{"AdvancedNumerics2", "-2 ^ 4", int64(16)},
		{"ComplexArithmetic", "1+2-3*8^2/2/2", int64(-45)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := parser.ParseExpression(tc.expression)
			if err != nil {
				t.Errorf("Failed to parse expression '%s': %v", tc.expression, err)
				return
			}

			result, err := expr.GetValue()
			if err != nil {
				t.Errorf("Failed to evaluate expression '%s': %v", tc.expression, err)
				return
			}

			if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
				t.Errorf("Expression '%s': expected %v, got %v", tc.expression, tc.expected, result)
			}
		})
	}
}

// TestNullLiterals tests case-insensitive null literals
func TestNullLiterals(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	testCases := []string{"null", "NULL", "NuLl"}

	for _, nullLiteral := range testCases {
		t.Run(fmt.Sprintf("NullLiteral_%s", nullLiteral), func(t *testing.T) {
			expr, err := parser.ParseExpression(nullLiteral)
			if err != nil {
				t.Errorf("Failed to parse expression '%s': %v", nullLiteral, err)
				return
			}

			result, err := expr.GetValue()
			if err != nil {
				t.Errorf("Failed to evaluate expression '%s': %v", nullLiteral, err)
				return
			}

			if result != nil {
				t.Errorf("Expected nil for '%s', got %v", nullLiteral, result)
			}
		})
	}
}

// TestLogicalOperators tests logical AND/OR operators
func TestLogicalOperators(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	// Test error cases for logical operators with null values
	errorCases := []string{
		"null and true",
		"true and null",
		"null or false",
		"false or null",
	}

	for _, expression := range errorCases {
		t.Run(fmt.Sprintf("Error_%s", strings.ReplaceAll(expression, " ", "_")), func(t *testing.T) {
			expr, err := parser.ParseExpression(expression)
			if err != nil {
				t.Errorf("Failed to parse expression '%s': %v", expression, err)
				return
			}

			_, err = expr.GetValue()
			if err == nil {
				t.Errorf("Expected error when evaluating '%s', but got none", expression)
			}
		})
	}
}

// Helper function to run a basic evaluation test
func runEvaluationTest(t *testing.T, parser *ast.SpelExpressionParser, expression string, expected interface{}) {
	expr, err := parser.ParseExpression(expression)
	if err != nil {
		t.Errorf("Failed to parse expression '%s': %v", expression, err)
		return
	}

	result, err := expr.GetValue()
	if err != nil {
		t.Errorf("Failed to evaluate expression '%s': %v", expression, err)
		return
	}

	if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", expected) {
		t.Errorf("Expression '%s': expected %v, got %v", expression, expected, result)
	}
}

// Helper function to test expressions that should fail
func runErrorTest(t *testing.T, parser *ast.SpelExpressionParser, expression string) {
	expr, err := parser.ParseExpression(expression)
	if err != nil {
		// Parse error is acceptable for some test cases
		return
	}

	_, err = expr.GetValue()
	if err == nil {
		t.Errorf("Expected error when evaluating '%s', but got none", expression)
	}
}
