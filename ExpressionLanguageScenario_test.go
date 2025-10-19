package main

import (
	"github.com/weaweawe01/ParserSpel/ast"
	"strings"
	"testing"
)

// Lightweight parser/AST tests inspired by ExpressionLanguageScenarioTests.java.
// These tests focus on parsing and AST string outputs using existing Go APIs.
func TestBasicParsingAndAST(t *testing.T) {
	parser := ast.NewSpelExpressionParser()

	// Simple literal parsing
	expr, err := parser.ParseExpressionWithContext("'hello world'", nil)
	if err != nil {
		t.Fatalf("failed to parse literal expression: %v", err)
	}
	// Ensure AST contains the literal content
	if !strings.Contains(expr.ToStringAST(), "hello world") {
		t.Fatalf("AST should contain literal 'hello world', got: %s", expr.ToStringAST())
	}

	// Variable reference parsing (we only test parsing, not evaluation wiring)
	expr, err = parser.ParseExpressionWithContext("#favouriteColour", nil)
	if err != nil {
		t.Fatalf("failed to parse variable reference: %v", err)
	}
	if !strings.Contains(expr.ToStringAST(), "#favouriteColour") {
		t.Fatalf("AST should contain variable reference '#favouriteColour', got: %s", expr.ToStringAST())
	}

	// List access parsing - just ensure parse succeeds
	expr, err = parser.ParseExpressionWithContext("{2,3,5,7}.get(1)", nil)
	if err != nil {
		t.Fatalf("failed to parse list access: %v", err)
	}
	if !strings.Contains(expr.ToStringAST(), "get(1)") {
		t.Fatalf("AST should contain get(1), got: %s", expr.ToStringAST())
	}
}
