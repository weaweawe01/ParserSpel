package main

import (
	"github.com/weaweawe01/ParserSpel/ast"
	"strings"
	"testing"
)

// TestInlineListCreation tests basic inline list creation
func TestInlineListCreation(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	tests := []struct {
		expression string
		desc       string
		shouldPass bool
	}{
		{"{1, 2, 3, 4, 5}", "simple integer list", false}, // Not implemented yet
		{"{'abc', 'xyz'}", "string list", false},          // Not implemented yet
		{"{}", "empty list", false},                       // Not implemented yet
		{"{'abc'=='xyz'}", "list with expression", false}, // Not implemented yet
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			expr, err := parser.ParseExpression(test.expression)
			if test.shouldPass {
				if err != nil {
					t.Errorf("Expected expression '%s' to parse successfully, but got error: %v", test.expression, err)
				} else {
					// Check AST contains inline list structure
					ast := expr.AST.ToStringAST()
					if !strings.Contains(ast, "InlineList") && !strings.Contains(ast, "List") {
						t.Logf("Expression '%s' parsed to AST: %s", test.expression, ast)
					}
				}
			} else {
				// For now, expect parsing to fail since inline lists aren't implemented
				if err == nil {
					t.Logf("Expression '%s' unexpectedly parsed successfully: %s", test.expression, expr.AST.ToStringAST())
				} else {
					t.Logf("Expression '%s' failed as expected: %v", test.expression, err)
				}
			}
		})
	}
}

// TestInlineListAndNesting tests nested list structures
func TestInlineListAndNesting(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	tests := []struct {
		expression string
		desc       string
		shouldPass bool
	}{
		{"{{1,2,3},{4,5,6}}", "nested integer lists", false},
		{"{{1,'2',3},{4,{'a','b'},5,6}}", "complex nested lists", false},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			expr, err := parser.ParseExpression(test.expression)
			if test.shouldPass {
				if err != nil {
					t.Errorf("Expected expression '%s' to parse successfully, but got error: %v", test.expression, err)
				}
			} else {
				// For now, expect parsing to fail since inline lists aren't implemented
				if err == nil {
					t.Logf("Expression '%s' unexpectedly parsed successfully: %s", test.expression, expr.AST.ToStringAST())
				} else {
					t.Logf("Expression '%s' failed as expected: %v", test.expression, err)
				}
			}
		})
	}
}

// TestInlineListError tests error cases for inline lists
func TestInlineListError(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	// Test incomplete list expression
	expr, err := parser.ParseExpression("{'abc'")
	if err == nil {
		t.Errorf("Expected parsing of incomplete list to fail, but got AST: %s", expr.AST.ToStringAST())
	} else {
		t.Logf("Incomplete list correctly failed with error: %v", err)
	}
}

// TestInlineListAndInstanceofOperator tests instanceof with lists
func TestInlineListAndInstanceofOperator(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	// This test will require both inline lists and instanceof operator
	expr, err := parser.ParseExpression("{1, 2, 3, 4, 5} instanceof T(java.util.List)")
	if err != nil {
		t.Logf("Expression with instanceof not yet supported: %v", err)
	} else {
		t.Logf("Expression parsed: %s", expr.AST.ToStringAST())
	}
}

// TestInlineListAndBetweenOperator tests between operator with lists
func TestInlineListAndBetweenOperator(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	tests := []struct {
		expression string
		desc       string
		expected   string
	}{
		{"1 between {1,5}", "integer in range", "true"},
		{"3 between {1,5}", "integer in middle", "true"},
		{"5 between {1,5}", "integer at end", "true"},
		{"0 between {1,5}", "integer below range", "false"},
		{"8 between {1,5}", "integer above range", "false"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			expr, err := parser.ParseExpression(test.expression)
			if err != nil {
				t.Logf("Between operator not yet supported for '%s': %v", test.expression, err)
			} else {
				t.Logf("Expression '%s' parsed: %s", test.expression, expr.AST.ToStringAST())
				// TODO: Add evaluation when between operator is implemented
			}
		})
	}
}

// TestInlineListAndBetweenOperatorForStrings tests between operator with string lists
func TestInlineListAndBetweenOperatorForStrings(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	tests := []struct {
		expression string
		desc       string
	}{
		{"'a' between {'a', 'c'}", "string at start"},
		{"'b' between {'a', 'c'}", "string in middle"},
		{"'c' between {'a', 'c'}", "string at end"},
		{"'z' between {'a', 'c'}", "string outside range"},
		{"'efg' between {'abc', 'xyz'}", "string lexical range"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			expr, err := parser.ParseExpression(test.expression)
			if err != nil {
				t.Logf("String between not yet supported for '%s': %v", test.expression, err)
			} else {
				t.Logf("Expression '%s' parsed: %s", test.expression, expr.AST.ToStringAST())
			}
		})
	}
}

// TestProjectionOnInlineList tests projection operator on inline lists
func TestProjectionOnInlineList(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	tests := []struct {
		expression string
		desc       string
	}{
		{"{1,2,3,4,5,6}.![#this>3]", "boolean projection"},
		{"{1,2,3,4,5,6,7,8,9,10}.![#this<5?'y':'n']", "conditional projection"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			expr, err := parser.ParseExpression(test.expression)
			if err != nil {
				t.Logf("Projection operator not yet supported for '%s': %v", test.expression, err)
			} else {
				t.Logf("Expression '%s' parsed: %s", test.expression, expr.AST.ToStringAST())
			}
		})
	}
}

// TestSelectionOnInlineList tests selection operator on inline lists
func TestSelectionOnInlineList(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	tests := []struct {
		expression string
		desc       string
	}{
		{"{1,2,3,4,5,6}.?[#this>3]", "filter selection"},
		{"{1,2,3,4,5,6,7,8,9,10}.?[#isEven(#this) == 'y']", "function-based selection"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			expr, err := parser.ParseExpression(test.expression)
			if err != nil {
				t.Logf("Selection operator not yet supported for '%s': %v", test.expression, err)
			} else {
				t.Logf("Expression '%s' parsed: %s", test.expression, expr.AST.ToStringAST())
			}
		})
	}
}

// TestSelectFirstAndLastOnList tests first/last selection operators
func TestSelectFirstAndLastOnList(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	tests := []struct {
		expression string
		desc       string
	}{
		{"listOfNumbersUpToTen.^[#isEven(#this) == 'y']", "select first"},
		{"listOfNumbersUpToTen.$[#isEven(#this) == 'y']", "select last"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			expr, err := parser.ParseExpression(test.expression)
			if err != nil {
				t.Logf("First/last selection not yet supported for '%s': %v", test.expression, err)
			} else {
				t.Logf("Expression '%s' parsed: %s", test.expression, expr.AST.ToStringAST())
			}
		})
	}
}

// TestSetConstructionWithInlineList tests using inline lists with constructors
func TestSetConstructionWithInlineList(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	expr, err := parser.ParseExpression("new java.util.HashSet().addAll({'a','b','c'})")
	if err != nil {
		t.Logf("Set construction with inline list not yet supported: %v", err)
	} else {
		t.Logf("Set construction parsed: %s", expr.AST.ToStringAST())
	}
}

// TestConstantRepresentation tests which lists are considered constant
func TestConstantRepresentation(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	tests := []struct {
		expression  string
		desc        string
		isConstant  bool
		shouldParse bool
	}{
		{"{1,2,3,4,5}", "literal numbers", true, false},
		{"{'abc'}", "literal string", true, false},
		{"{}", "empty list", true, false},
		{"{#a,2,3}", "variable reference", false, false},
		{"{1,2,Integer.valueOf(4)}", "method call", false, false},
		{"{1,2,{#a}}", "nested variable", false, false},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			expr, err := parser.ParseExpression(test.expression)
			if test.shouldParse {
				if err != nil {
					t.Errorf("Expected expression '%s' to parse successfully, but got error: %v", test.expression, err)
				} else {
					// TODO: Check if the AST node has constant representation when implemented
					t.Logf("Expression '%s' parsed successfully", test.expression)
				}
			} else {
				if err == nil {
					t.Logf("Expression '%s' unexpectedly parsed: %s", test.expression, expr.AST.ToStringAST())
				} else {
					t.Logf("Expression '%s' failed as expected: %v", test.expression, err)
				}
			}
		})
	}
}

// TestInlineListWriting tests list modification
func TestInlineListWriting(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	// Test that lists should be unmodifiable
	expr, err := parser.ParseExpression("{1, 2, 3, 4, 5}[0]=6")
	if err != nil {
		t.Logf("List indexing assignment not yet supported: %v", err)
	} else {
		t.Logf("List assignment parsed: %s", expr.AST.ToStringAST())
		// TODO: Add evaluation test to ensure UnsupportedOperationException equivalent
	}
}
