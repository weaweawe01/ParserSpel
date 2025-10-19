package ast

import (
	"fmt"
	"strings"
)

// InternalSpelExpressionParser handles the parsing of SpEL expressions
type InternalSpelExpressionParser struct {
	Configuration      *SpelParserConfiguration
	ExpressionString   string
	TokenStream        []*Token
	TokenStreamLength  int
	TokenStreamPointer int
	ConstructedNodes   []SpelNode
}

// NewInternalSpelExpressionParser creates a new parser instance
func NewInternalSpelExpressionParser(config *SpelParserConfiguration) *InternalSpelExpressionParser {
	return &InternalSpelExpressionParser{
		Configuration:    config,
		ConstructedNodes: make([]SpelNode, 0),
	}
}

// DoParseExpression parses a SpEL expression string into an AST with debug output
func (p *InternalSpelExpressionParser) DoParseExpression(expressionString string) (*SpelExpression, error) {
	p.ExpressionString = expressionString

	// Check expression length
	if len(expressionString) > p.Configuration.MaximumExpressionLength {
		return nil, fmt.Errorf("expression exceeds maximum length of %d", p.Configuration.MaximumExpressionLength)
	}

	// Tokenize the expression
	tokenizer := NewTokenizer(expressionString)
	tokens, err := tokenizer.Process()
	if err != nil {
		return nil, fmt.Errorf("tokenization failed: %v", err)
	}

	p.TokenStream = tokens
	p.TokenStreamLength = len(tokens)
	p.TokenStreamPointer = 0
	p.ConstructedNodes = make([]SpelNode, 0)

	// Print token stream (matching Java debug output)
	for count, token := range p.TokenStream {
		fmt.Printf("[%d] %s\n", count, token)
	}

	// Parse the tokens into an AST
	ast, err := p.eatExpression()
	fmt.Println("AST->", ast)
	if err != nil {
		return nil, fmt.Errorf("parsing failed: %v", err)
	}

	if ast == nil {
		return nil, fmt.Errorf("empty expression - OOD")
	}

	// Check if there's an assignment operator after the first expression
	if p.peekToken(ASSIGN) {
		p.takeToken() // consume '='

		// Parse the right-hand side expression
		rightExpr, err := p.eatExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse assignment right-hand side: %v", err)
		}

		if rightExpr == nil {
			return nil, fmt.Errorf("expected expression after '='")
		}

		// Create assignment node
		assignNode := NewAssign(ast, rightExpr, ast.GetStartPosition(), rightExpr.GetEndPosition())
		ast = assignNode
	}

	// Print the AST (matching Java debug output)
	fmt.Println(ast.ToStringAST())

	// Check if all tokens were consumed
	if p.TokenStreamPointer < p.TokenStreamLength {
		nextToken := p.peekTokenRaw()
		return nil, fmt.Errorf("unexpected tokens after expression: %v", nextToken)
	}

	return NewSpelExpression(expressionString, ast, p.Configuration), nil
}

// ParseExpression parses a SpEL expression string into an AST (without debug output)
func (p *InternalSpelExpressionParser) ParseExpression(expressionString string) (*SpelExpression, error) {
	p.ExpressionString = expressionString

	// Check expression length
	if len(expressionString) > p.Configuration.MaximumExpressionLength {
		return nil, fmt.Errorf("expression exceeds maximum length of %d", p.Configuration.MaximumExpressionLength)
	}

	// Tokenize the expression
	tokenizer := NewTokenizer(expressionString)
	tokens, err := tokenizer.Process()
	if err != nil {
		return nil, fmt.Errorf("tokenization failed: %v", err)
	}
	p.TokenStream = tokens
	p.TokenStreamLength = len(tokens)
	p.TokenStreamPointer = 0
	p.ConstructedNodes = make([]SpelNode, 0)

	// Parse the tokens into an AST
	ast, err := p.eatExpression()
	if err != nil {
		return nil, fmt.Errorf("parsing failed: %v", err)
	}

	if ast == nil {
		return nil, fmt.Errorf("empty expression")
	}

	// Check if there's an assignment operator after the first expression
	if p.peekToken(ASSIGN) {
		p.takeToken() // consume '='

		// Parse the right-hand side expression
		rightExpr, err := p.eatExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse assignment right-hand side: %v", err)
		}

		if rightExpr == nil {
			return nil, fmt.Errorf("expected expression after '='")
		}

		// Create assignment node
		assignNode := NewAssign(ast, rightExpr, ast.GetStartPosition(), rightExpr.GetEndPosition())
		ast = assignNode
	}

	// Check if all tokens were consumed
	if p.TokenStreamPointer < p.TokenStreamLength {
		return nil, fmt.Errorf("unexpected tokens after expression")
	}

	return NewSpelExpression(expressionString, ast, p.Configuration), nil
}

// eatExpression parses the top-level expression
func (p *InternalSpelExpressionParser) eatExpression() (SpelNode, error) {
	return p.eatTernaryExpression()
}

// eatTernaryExpression parses ternary conditional expressions (condition ? trueValue : falseValue)
func (p *InternalSpelExpressionParser) eatTernaryExpression() (SpelNode, error) {
	expr, err := p.eatLogicalOrExpression()
	if err != nil {
		return nil, err
	}

	// Check for Elvis operator (?:)
	if p.peekToken(ELVIS) {
		p.takeToken() // consume '?:'

		// Parse default value
		defaultValue, err := p.eatLogicalOrExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse default value in Elvis expression: %v", err)
		}

		if defaultValue == nil {
			return nil, fmt.Errorf("expected expression after '?:' in Elvis operator")
		}

		// Create Elvis node
		elvis := NewElvis(expr, defaultValue, expr.GetStartPosition(), defaultValue.GetEndPosition())
		return elvis, nil
	}

	// Check for ternary operator
	if p.peekToken(QMARK) {
		p.takeToken() // consume '?'

		// Parse true value
		trueValue, err := p.eatLogicalOrExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse true value in ternary expression: %v", err)
		}

		if trueValue == nil {
			return nil, fmt.Errorf("expected expression after '?' in ternary")
		}

		// Expect colon
		if !p.peekToken(COLON) {
			return nil, fmt.Errorf("expected ':' in ternary expression")
		}
		p.takeToken() // consume ':'

		// Parse false value
		falseValue, err := p.eatLogicalOrExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse false value in ternary expression: %v", err)
		}

		if falseValue == nil {
			return nil, fmt.Errorf("expected expression after ':' in ternary")
		}

		// Create ternary node
		ternary := NewTernary(expr, trueValue, falseValue, expr.GetStartPosition(), falseValue.GetEndPosition())
		return ternary, nil
	}

	return expr, nil
}

// eatLogicalOrExpression parses logical OR expressions
func (p *InternalSpelExpressionParser) eatLogicalOrExpression() (SpelNode, error) {
	expr, err := p.eatLogicalAndExpression()
	if err != nil {
		return nil, err
	}

	for p.peekToken(SYMBOLIC_OR) || p.peekIdentifierToken("or") {
		token := p.takeToken()
		right, err := p.eatLogicalAndExpression()
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, fmt.Errorf("missing right operand for OR at position %d", token.StartPos)
		}
		expr = NewOpOr(expr, right, expr.GetStartPosition(), right.GetEndPosition())
	}

	return expr, nil
}

// eatLogicalAndExpression parses logical AND expressions
func (p *InternalSpelExpressionParser) eatLogicalAndExpression() (SpelNode, error) {
	expr, err := p.eatRelationalExpression()
	if err != nil {
		return nil, err
	}

	for p.peekToken(SYMBOLIC_AND) || p.peekIdentifierToken("and") {
		token := p.takeToken()
		right, err := p.eatRelationalExpression()
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, fmt.Errorf("missing right operand for AND at position %d", token.StartPos)
		}
		expr = NewOpAnd(expr, right, expr.GetStartPosition(), right.GetEndPosition())
	}

	return expr, nil
}

// eatRelationalExpression parses relational expressions
func (p *InternalSpelExpressionParser) eatRelationalExpression() (SpelNode, error) {
	expr, err := p.eatSumExpression()
	if err != nil {
		return nil, err
	}

	relationalOperatorToken := p.maybeEatRelationalOperator()
	if relationalOperatorToken != nil {
		right, err := p.eatSumExpression()
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, fmt.Errorf("missing right operand for %s at position %d",
				relationalOperatorToken.Kind.String(), relationalOperatorToken.StartPos)
		}

		startPos := expr.GetStartPosition()
		endPos := right.GetEndPosition()

		switch relationalOperatorToken.Kind {
		case EQ:
			expr = NewOpEQ(expr, right, startPos, endPos)
		case NE:
			expr = NewOpNE(expr, right, startPos, endPos)
		case GT:
			expr = NewOpGT(expr, right, startPos, endPos)
		case GE:
			expr = NewOpGE(expr, right, startPos, endPos)
		case LT:
			expr = NewOpLT(expr, right, startPos, endPos)
		case LE:
			expr = NewOpLE(expr, right, startPos, endPos)
		case MATCHES:
			expr = NewOperatorMatches(expr, right, startPos, endPos)
		case BETWEEN:
			expr = NewOperatorBetween(expr, right, startPos, endPos)
		default:
			return nil, fmt.Errorf("unsupported relational operator: %s", relationalOperatorToken.Kind.String())
		}
	}

	return expr, nil
}

// eatSumExpression parses addition and subtraction expressions
func (p *InternalSpelExpressionParser) eatSumExpression() (SpelNode, error) {
	expr, err := p.eatProductExpression()
	if err != nil {
		return nil, err
	}

	for p.peekToken(PLUS) || p.peekToken(MINUS) {
		token := p.takeToken()
		right, err := p.eatProductExpression()
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, fmt.Errorf("missing right operand for %s at position %d",
				token.Kind.String(), token.StartPos)
		}

		startPos := expr.GetStartPosition()
		endPos := right.GetEndPosition()

		switch token.Kind {
		case PLUS:
			expr = NewOpPlus(expr, right, startPos, endPos)
		case MINUS:
			expr = NewOpMinus(expr, right, startPos, endPos)
		}
	}

	return expr, nil
}

// eatProductExpression parses multiplication, division, and modulo expressions
func (p *InternalSpelExpressionParser) eatProductExpression() (SpelNode, error) {
	expr, err := p.eatPowerIncDecExpression()
	if err != nil {
		return nil, err
	}

	for p.peekToken(STAR) || p.peekToken(DIV) || p.peekToken(MOD) {
		token := p.takeToken()
		right, err := p.eatPowerIncDecExpression()
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, fmt.Errorf("missing right operand for %s at position %d",
				token.Kind.String(), token.StartPos)
		}

		startPos := expr.GetStartPosition()
		endPos := right.GetEndPosition()

		switch token.Kind {
		case STAR:
			expr = NewOpMultiply(expr, right, startPos, endPos)
		case DIV:
			expr = NewOpDivide(expr, right, startPos, endPos)
		case MOD:
			expr = NewOpModulus(expr, right, startPos, endPos)
		}
	}

	return expr, nil
}

// eatPowerIncDecExpression parses power expressions
func (p *InternalSpelExpressionParser) eatPowerIncDecExpression() (SpelNode, error) {
	expr, err := p.eatUnaryExpression()
	if err != nil {
		return nil, err
	}

	if p.peekToken(POWER) {
		token := p.takeToken() // consume POWER
		right, err := p.eatUnaryExpression()
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, fmt.Errorf("missing right operand for power operator at position %d", token.StartPos)
		}

		startPos := expr.GetStartPosition()
		endPos := right.GetEndPosition()
		return NewOperatorPower(expr, right, startPos, endPos), nil
	}

	// TODO: Handle INC/DEC operators if needed
	return expr, nil
}

// eatUnaryExpression parses unary expressions
func (p *InternalSpelExpressionParser) eatUnaryExpression() (SpelNode, error) {
	if p.peekToken(NOT) || p.peekToken(PLUS) || p.peekToken(MINUS) {
		token := p.takeToken()
		child, err := p.eatUnaryExpression()
		if err != nil {
			return nil, err
		}
		if child == nil {
			return nil, fmt.Errorf("missing operand for unary %s at position %d",
				token.Kind.String(), token.StartPos)
		}

		switch token.Kind {
		case NOT:
			return NewOperatorNot(child, token.StartPos, child.GetEndPosition()), nil
		case PLUS:
			// Unary plus - just return the child
			return child, nil
		case MINUS:
			// Unary minus
			return NewUnaryOpMinus(child, token.StartPos, child.GetEndPosition()), nil
		}
	}

	return p.eatPrimaryExpression()
}

// eatPrimaryExpression parses primary expressions
func (p *InternalSpelExpressionParser) eatPrimaryExpression() (SpelNode, error) {
	start, err := p.eatStartNode()
	if err != nil {
		return nil, err
	}
	if start == nil {
		return nil, fmt.Errorf("expected primary expression")
	}

	// Handle compound expressions (property access, method calls, etc.)
	nodes := []SpelNode{start}

	for {
		node, err := p.eatNode()
		if err != nil {
			return nil, err
		}
		if node == nil {
			break
		}
		nodes = append(nodes, node)
	}

	if len(nodes) == 1 {
		return nodes[0], nil
	}

	// Create compound expression
	return NewCompoundExpression(
		nodes[0].GetStartPosition(),
		nodes[len(nodes)-1].GetEndPosition(),
		nodes...), nil
}

// eatStartNode parses start nodes (literals, identifiers, parenthesized expressions)
func (p *InternalSpelExpressionParser) eatStartNode() (SpelNode, error) {
	if p.maybeEatLiteral() {
		return p.pop(), nil
	}

	if p.maybeEatParenExpression() {
		return p.pop(), nil
	}

	if p.maybeEatBeanReference() {
		return p.pop(), nil
	}

	if p.maybeEatVariableReference() {
		return p.pop(), nil
	}

	if p.maybeEatNullReference() {
		return p.pop(), nil
	}

	if p.maybeEatTypeReference() {
		return p.pop(), nil
	}

	if p.maybeEatConstructorExpression() {
		return p.pop(), nil
	}

	if p.maybeEatInlineCollection() {
		return p.pop(), nil
	}

	if p.maybeEatMethodCall() {
		return p.pop(), nil
	}

	if p.maybeEatIdentifier() {
		return p.pop(), nil
	}

	return nil, fmt.Errorf("unexpected token: %v", p.peekTokenRaw())
}

// eatNode parses node expressions (property access, indexing, method calls, etc.)
func (p *InternalSpelExpressionParser) eatNode() (SpelNode, error) {
	if p.peekToken(DOT) || p.peekToken(SAFE_NAVI) {
		return p.eatDottedNode()
	}

	if p.peekToken(LSQUARE) {
		// Handle indexing properly
		startToken := p.takeToken() // consume '['

		// Parse the index expression
		indexExpr, err := p.eatExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse index expression: %v", err)
		}

		// Expect closing bracket
		if !p.peekToken(RSQUARE) {
			return nil, fmt.Errorf("expected ']' after index expression")
		}
		endToken := p.takeToken() // consume ']'

		// Create indexer node
		indexer := NewIndexer(indexExpr, startToken.StartPos, endToken.EndPos)
		return indexer, nil
	}

	// Handle method calls
	if p.peekToken(LPAREN) {
		return p.eatMethodCall()
	}

	return nil, nil
}

// eatDottedNode parses dotted expressions (property access and method calls)
func (p *InternalSpelExpressionParser) eatDottedNode() (SpelNode, error) {
	token := p.takeToken() // consume '.' or '?.'
	nullSafeNavigation := (token.Kind == SAFE_NAVI)

	// Check for selection expressions like .?[criteria] or .![criteria]
	if p.peekToken(SELECT) {
		p.takeToken() // consume '?['

		// Parse the selection criteria
		criteria, err := p.eatExpression()
		if err != nil {
			return nil, fmt.Errorf("error parsing selection criteria: %v", err)
		}

		if criteria == nil {
			return nil, fmt.Errorf("expected criteria in selection expression")
		}

		// Expect closing bracket
		if !p.peekToken(RSQUARE) {
			return nil, fmt.Errorf("expected ']' to close selection expression")
		}

		endToken := p.takeToken() // consume ']'
		endPos := endToken.EndPos

		// Create selection node (.?[ is always null-safe selection)
		selection := NewSelection(true, SelectionAll, criteria, token.StartPos, endPos)
		p.push(selection)
		return selection, nil
	}

	// Check for first selection expressions like .^[criteria] or ?^[criteria]
	if p.peekToken(SELECT_FIRST) {
		p.takeToken() // consume '^['

		// Parse the selection criteria
		criteria, err := p.eatExpression()
		if err != nil {
			return nil, fmt.Errorf("error parsing first selection criteria: %v", err)
		}

		if criteria == nil {
			return nil, fmt.Errorf("expected criteria in first selection expression")
		}

		// Expect closing bracket
		if !p.peekToken(RSQUARE) {
			return nil, fmt.Errorf("expected ']' to close first selection expression")
		}

		endToken := p.takeToken() // consume ']'
		endPos := endToken.EndPos

		// Create first selection node
		selection := NewSelection(nullSafeNavigation, SelectionFirst, criteria, token.StartPos, endPos)
		p.push(selection)
		return selection, nil
	}

	// Check for last selection expressions like .$[criteria] or ?$[criteria]
	if p.peekToken(SELECT_LAST) {
		p.takeToken() // consume '$['

		// Parse the selection criteria
		criteria, err := p.eatExpression()
		if err != nil {
			return nil, fmt.Errorf("error parsing last selection criteria: %v", err)
		}

		if criteria == nil {
			return nil, fmt.Errorf("expected criteria in last selection expression")
		}

		// Expect closing bracket
		if !p.peekToken(RSQUARE) {
			return nil, fmt.Errorf("expected ']' to close last selection expression")
		}

		endToken := p.takeToken() // consume ']'
		endPos := endToken.EndPos

		// Create last selection node
		selection := NewSelection(nullSafeNavigation, SelectionLast, criteria, token.StartPos, endPos)
		p.push(selection)
		return selection, nil
	}

	// Check for projection expressions like .![expression]
	if p.peekToken(PROJECT) {
		p.takeToken() // consume '!['

		// Parse the projection expression
		projectionExpr, err := p.eatExpression()
		if err != nil {
			return nil, fmt.Errorf("error parsing projection expression: %v", err)
		}

		if projectionExpr == nil {
			return nil, fmt.Errorf("expected expression in projection")
		}

		// Expect closing bracket
		if !p.peekToken(RSQUARE) {
			return nil, fmt.Errorf("expected ']' to close projection expression")
		}

		endToken := p.takeToken() // consume ']'
		endPos := endToken.EndPos

		// Create projection node
		projection := NewProjection(projectionExpr, token.StartPos, endPos)
		p.push(projection)
		return projection, nil
	}

	// Check for safe navigation indexing like ?.[index]
	if p.peekToken(LSQUARE) {
		p.takeToken() // consume '['

		// Parse the index expression
		indexExpr, err := p.eatExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse index expression: %v", err)
		}

		if indexExpr == nil {
			return nil, fmt.Errorf("expected index expression")
		}

		// Expect closing bracket
		if !p.peekToken(RSQUARE) {
			return nil, fmt.Errorf("expected ']' after index expression")
		}
		endBracket := p.takeToken() // consume ']'

		// Create indexer node (null-safe if preceded by ?.)
		var indexer *Indexer
		if nullSafeNavigation {
			indexer = NewNullSafeIndexer(indexExpr, token.StartPos, endBracket.EndPos)
		} else {
			indexer = NewIndexer(indexExpr, token.StartPos, endBracket.EndPos)
		}
		p.push(indexer)
		return indexer, nil
	}

	if p.peekToken(IDENTIFIER) {
		identToken := p.takeToken()
		propertyName := identToken.StringValue()

		// Check if this is followed by parentheses (method call)
		if p.peekToken(LPAREN) {
			// This is a method call
			p.takeToken() // consume '('

			// Parse method arguments
			var arguments []SpelNode

			// Handle empty argument list
			if p.peekToken(RPAREN) {
				p.takeToken() // consume ')'
				endPos := p.TokenStream[p.TokenStreamPointer-1].EndPos
				methodRef := NewMethodReference(nullSafeNavigation, propertyName, arguments, token.StartPos, endPos)
				p.push(methodRef)
				return methodRef, nil
			}

			// Parse arguments separated by commas
			for {
				arg, err := p.eatExpression()
				if err != nil {
					return nil, fmt.Errorf("error parsing method argument: %v", err)
				}
				if arg == nil {
					return nil, fmt.Errorf("expected argument in method call")
				}
				arguments = append(arguments, arg)

				if p.peekToken(COMMA) {
					p.takeToken() // consume ','
					continue
				}
				break
			}

			if !p.peekToken(RPAREN) {
				return nil, fmt.Errorf("expected ')' to close method call")
			}

			endToken := p.takeToken() // consume ')'
			endPos := endToken.EndPos

			methodRef := NewMethodReference(nullSafeNavigation, propertyName, arguments, token.StartPos, endPos)
			p.push(methodRef)
			return methodRef, nil
		} else {
			// This is a property/field reference
			node := NewPropertyOrFieldReference(nullSafeNavigation, propertyName, token.StartPos, identToken.EndPos)
			p.push(node)
			return node, nil
		}
	}

	return nil, fmt.Errorf("expected identifier after %s at position %d", token.Kind.String(), token.StartPos)
}

// Literal parsing methods
func (p *InternalSpelExpressionParser) maybeEatLiteral() bool {
	token := p.peekTokenRaw()
	if token == nil {
		return false
	}

	switch token.Kind {
	case LITERAL_INT, LITERAL_LONG, LITERAL_HEXINT, LITERAL_HEXLONG:
		p.takeToken()
		value, err := parseNumber(token.StringValue(), token.Kind)
		if err != nil {
			return false
		}
		literal := NewIntLiteral(value, token.StartPos, token.EndPos)
		p.push(literal)
		return true

	case LITERAL_REAL, LITERAL_REAL_FLOAT:
		p.takeToken()
		value, err := parseNumber(token.StringValue(), token.Kind)
		if err != nil {
			return false
		}
		// Create RealLiteral for floating point numbers
		if floatVal, ok := value.(float64); ok {
			realLiteral := NewRealLiteral(floatVal, token.StartPos, token.EndPos)
			p.push(realLiteral)
		} else if floatVal, ok := value.(float32); ok {
			realLiteral := NewRealLiteral(float64(floatVal), token.StartPos, token.EndPos)
			p.push(realLiteral)
		} else {
			// Fallback to regular literal
			literal := NewIntLiteral(value, token.StartPos, token.EndPos)
			p.push(literal)
		}
		return true

	case LITERAL_STRING:
		p.takeToken()
		// Remove quotes from string literal
		strValue := token.StringValue()
		if len(strValue) >= 2 {
			strValue = strValue[1 : len(strValue)-1] // Remove surrounding quotes
		}
		literal := NewStringLiteral(strValue, token.StartPos, token.EndPos)
		p.push(literal)
		return true

	case IDENTIFIER:
		// Check for boolean literals
		tokenValue := strings.ToLower(token.StringValue())
		if tokenValue == "true" || tokenValue == "false" {
			p.takeToken()
			boolValue := tokenValue == "true"
			literal := NewBooleanLiteral(boolValue, token.StartPos, token.EndPos)
			p.push(literal)
			return true
		}
	}

	return false
}

func (p *InternalSpelExpressionParser) maybeEatParenExpression() bool {
	if !p.peekToken(LPAREN) {
		return false
	}

	p.takeToken() // consume '('
	expr, err := p.eatExpression()
	if err != nil {
		return false
	}

	if !p.peekToken(RPAREN) {
		return false
	}

	p.takeToken() // consume ')'
	p.push(expr)
	return true
}

func (p *InternalSpelExpressionParser) maybeEatBeanReference() bool {
	if !p.peekToken(BEAN_REF) && !p.peekToken(FACTORY_BEAN_REF) {
		return false
	}

	token := p.takeToken()

	if !p.peekToken(IDENTIFIER) && !p.peekToken(LITERAL_STRING) {
		return false
	}

	nameToken := p.takeToken()
	beanName := nameToken.StringValue()

	// Remove quotes if string literal
	if nameToken.Kind == LITERAL_STRING && len(beanName) >= 2 {
		beanName = beanName[1 : len(beanName)-1]
	}

	beanRef := NewBeanReference(beanName, token.StartPos, nameToken.EndPos)
	p.push(beanRef)
	return true
}

func (p *InternalSpelExpressionParser) maybeEatVariableReference() bool {
	if !p.peekToken(HASH) {
		return false
	}

	hashToken := p.takeToken()

	if !p.peekToken(IDENTIFIER) {
		return false
	}

	identToken := p.takeToken()
	varName := identToken.StringValue()

	// Check if this is a function call (#identifier(...))
	if p.peekToken(LPAREN) {
		// This is a function call, not a variable reference
		p.takeToken() // consume '('

		var arguments []SpelNode

		// Handle empty argument list
		if p.peekToken(RPAREN) {
			endToken := p.takeToken() // consume ')'
			endPos := endToken.EndPos
			funcRef := NewFunctionReference(varName, arguments, hashToken.StartPos, endPos)
			p.push(funcRef)
			return true
		}

		// Parse arguments separated by commas
		for {
			arg, err := p.eatExpression()
			if err != nil {
				return false
			}
			if arg == nil {
				return false
			}
			arguments = append(arguments, arg)

			if p.peekToken(COMMA) {
				p.takeToken() // consume ','
				continue
			}
			break
		}

		if !p.peekToken(RPAREN) {
			return false
		}

		endToken := p.takeToken() // consume ')'
		endPos := endToken.EndPos

		funcRef := NewFunctionReference(varName, arguments, hashToken.StartPos, endPos)
		p.push(funcRef)
		return true
	} else {
		// This is a variable reference
		varRef := NewVariableReference(varName, hashToken.StartPos, identToken.EndPos)
		p.push(varRef)
		return true
	}
}

func (p *InternalSpelExpressionParser) maybeEatNullReference() bool {
	if !p.peekToken(IDENTIFIER) {
		return false
	}

	token := p.peekTokenRaw()
	if token != nil && strings.ToLower(token.StringValue()) == "null" {
		p.takeToken()
		nullLiteral := NewNullLiteral(token.StartPos, token.EndPos)
		p.push(nullLiteral)
		return true
	}

	return false
}

func (p *InternalSpelExpressionParser) maybeEatIdentifier() bool {
	if !p.peekToken(IDENTIFIER) {
		return false
	}

	token := p.peekTokenRaw()
	// Don't treat reserved keywords as regular identifiers
	if token != nil && strings.ToLower(token.StringValue()) == "new" {
		return false
	}

	token = p.takeToken()
	// In SpEL, standalone identifiers are treated as direct property/field references
	// This matches the Java implementation behavior
	propertyRef := NewDirectPropertyOrFieldReference(token.StringValue(), token.StartPos, token.EndPos)
	p.push(propertyRef)
	return true
}

// maybeEatInlineCollection parses inline collections like {1,2,3} or {key:value}
func (p *InternalSpelExpressionParser) maybeEatInlineCollection() bool {
	if !p.peekToken(LCURLY) {
		return false
	}

	startToken := p.takeToken() // consume '{'

	// Handle empty collection {}
	if p.peekToken(RCURLY) {
		endToken := p.takeToken() // consume '}'
		// Default to empty list for {}
		inlineList := NewInlineList([]SpelNode{}, startToken.StartPos, endToken.EndPos)
		p.push(inlineList)
		return true
	}

	// Parse first element to determine if it's a list or map
	firstExpr, err := p.eatExpression()
	if err != nil {
		return false
	}

	// Check if this is a map (key:value format)
	if p.peekToken(COLON) {
		// This is a map literal
		p.takeToken() // consume ':'

		// Parse first value
		firstValue, err := p.eatExpression()
		if err != nil {
			return false
		}

		pairs := []KeyValuePair{
			{Key: firstExpr, Value: firstValue},
		}

		// Parse remaining key-value pairs
		for p.peekToken(COMMA) {
			p.takeToken() // consume ','

			// Parse key
			key, err := p.eatExpression()
			if err != nil {
				return false
			}

			// Expect colon
			if !p.peekToken(COLON) {
				return false
			}
			p.takeToken() // consume ':'

			// Parse value
			value, err := p.eatExpression()
			if err != nil {
				return false
			}

			pairs = append(pairs, KeyValuePair{Key: key, Value: value})
		}

		// Expect closing brace
		if !p.peekToken(RCURLY) {
			return false
		}
		endToken := p.takeToken() // consume '}'

		inlineMap := NewInlineMap(pairs, startToken.StartPos, endToken.EndPos)
		p.push(inlineMap)
		return true
	} else {
		// This is a list literal
		elements := []SpelNode{firstExpr}

		// Parse remaining elements
		for p.peekToken(COMMA) {
			p.takeToken() // consume ','
			expr, err := p.eatExpression()
			if err != nil {
				return false
			}
			elements = append(elements, expr)
		}

		// Expect closing brace
		if !p.peekToken(RCURLY) {
			return false
		}
		endToken := p.takeToken() // consume '}'

		inlineList := NewInlineList(elements, startToken.StartPos, endToken.EndPos)
		p.push(inlineList)
		return true
	}
}

// maybeEatMethodCall parses direct method calls like methodName(args...)
func (p *InternalSpelExpressionParser) maybeEatMethodCall() bool {
	// Look ahead to see if we have IDENTIFIER followed by LPAREN
	if !p.peekToken(IDENTIFIER) {
		return false
	}

	// Don't treat reserved keywords as method names
	token := p.peekTokenRaw()
	if token != nil && strings.ToLower(token.StringValue()) == "new" {
		return false
	}

	// Save current position to backtrack if needed
	savedPos := p.TokenStreamPointer

	// Check if next token after identifier is LPAREN
	p.TokenStreamPointer++
	if p.TokenStreamPointer >= p.TokenStreamLength || !p.peekToken(LPAREN) {
		// Backtrack and return false
		p.TokenStreamPointer = savedPos
		return false
	}

	// Restore position and parse as method call
	p.TokenStreamPointer = savedPos

	// Parse method name
	nameToken := p.takeToken()
	methodName := nameToken.StringValue()

	// Parse method arguments
	if !p.peekToken(LPAREN) {
		return false
	}

	startPos := nameToken.StartPos
	p.takeToken() // consume '('

	var arguments []SpelNode

	// Handle empty argument list
	if p.peekToken(RPAREN) {
		endToken := p.takeToken() // consume ')'
		endPos := endToken.EndPos
		methodRef := NewMethodReference(false, methodName, arguments, startPos, endPos)
		p.push(methodRef)
		return true
	}

	// Parse arguments separated by commas
	for {
		arg, err := p.eatExpression()
		if err != nil {
			return false
		}
		if arg == nil {
			return false
		}
		arguments = append(arguments, arg)

		if p.peekToken(COMMA) {
			p.takeToken() // consume ','
			continue
		}
		break
	}

	if !p.peekToken(RPAREN) {
		return false
	}

	endToken := p.takeToken() // consume ')'
	endPos := endToken.EndPos

	methodRef := NewMethodReference(false, methodName, arguments, startPos, endPos)
	p.push(methodRef)
	return true
} // Token manipulation methods
func (p *InternalSpelExpressionParser) peekToken(desiredKind TokenKind) bool {
	token := p.peekTokenRaw()
	return token != nil && token.Kind == desiredKind
}

func (p *InternalSpelExpressionParser) peekIdentifierToken(identifier string) bool {
	token := p.peekTokenRaw()
	return token != nil && token.Kind == IDENTIFIER &&
		strings.ToLower(token.StringValue()) == strings.ToLower(identifier)
}

func (p *InternalSpelExpressionParser) peekTokenRaw() *Token {
	if p.TokenStreamPointer >= p.TokenStreamLength {
		return nil
	}
	return p.TokenStream[p.TokenStreamPointer]
}

func (p *InternalSpelExpressionParser) takeToken() *Token {
	if p.TokenStreamPointer >= p.TokenStreamLength {
		return nil
	}
	token := p.TokenStream[p.TokenStreamPointer]
	p.TokenStreamPointer++
	return token
}

func (p *InternalSpelExpressionParser) maybeEatRelationalOperator() *Token {
	token := p.peekTokenRaw()
	if token == nil {
		return nil
	}

	switch token.Kind {
	case EQ, NE, GT, GE, LT, LE, BETWEEN, MATCHES, INSTANCEOF:
		return p.takeToken()
	case IDENTIFIER:
		tokenValue := strings.ToLower(token.StringValue())
		if tokenValue == "matches" || tokenValue == "instanceof" || tokenValue == "between" {
			token := p.takeToken()
			// Convert identifier token to appropriate operator token
			switch tokenValue {
			case "matches":
				token.Kind = MATCHES
			case "instanceof":
				token.Kind = INSTANCEOF
			case "between":
				token.Kind = BETWEEN
			}
			return token
		}
	}

	return nil
}

// Node stack operations
func (p *InternalSpelExpressionParser) push(node SpelNode) {
	p.ConstructedNodes = append(p.ConstructedNodes, node)
}

func (p *InternalSpelExpressionParser) pop() SpelNode {
	if len(p.ConstructedNodes) == 0 {
		return nil
	}

	node := p.ConstructedNodes[len(p.ConstructedNodes)-1]
	p.ConstructedNodes = p.ConstructedNodes[:len(p.ConstructedNodes)-1]
	return node
}

// eatMethodCall parses method calls with arguments
func (p *InternalSpelExpressionParser) eatMethodCall() (SpelNode, error) {
	if !p.peekToken(LPAREN) {
		return nil, fmt.Errorf("expected '(' for method call")
	}

	startPos := p.peekTokenRaw().StartPos
	p.takeToken() // consume '('

	// Parse arguments
	var arguments []SpelNode

	// Handle empty argument list
	if p.peekToken(RPAREN) {
		p.takeToken() // consume ')'
		endPos := p.TokenStream[p.TokenStreamPointer-1].EndPos
		methodCall := NewMethodReference(false, "", arguments, startPos, endPos)
		p.push(methodCall)
		return methodCall, nil
	}

	// Parse arguments separated by commas
	for {
		arg, err := p.eatExpression()
		if err != nil {
			return nil, fmt.Errorf("error parsing method argument: %v", err)
		}
		if arg == nil {
			return nil, fmt.Errorf("expected argument in method call")
		}
		arguments = append(arguments, arg)

		if p.peekToken(COMMA) {
			p.takeToken() // consume ','
			continue
		}
		break
	}

	if !p.peekToken(RPAREN) {
		return nil, fmt.Errorf("expected ')' to close method call")
	}

	endToken := p.takeToken() // consume ')'
	endPos := endToken.EndPos

	methodCall := NewMethodReference(false, "", arguments, startPos, endPos)
	p.push(methodCall)
	return methodCall, nil
}

// maybeEatConstructorExpression parses constructor expressions (new ClassName(...))
func (p *InternalSpelExpressionParser) maybeEatConstructorExpression() bool {
	if !p.peekIdentifierToken("new") {
		return false
	}

	newToken := p.takeToken() // consume 'new'

	// Parse the type/class name (may include dots like java.lang.String)
	var typeParts []string

	// Java Spring SpEL allows numbers and other tokens as type names in constructor expressions
	// Check for IDENTIFIER or numeric literals
	if !p.peekToken(IDENTIFIER) && !p.peekToken(LITERAL_INT) && !p.peekToken(LITERAL_LONG) &&
		!p.peekToken(LITERAL_HEXINT) && !p.peekToken(LITERAL_HEXLONG) {
		// 'new' keyword found but no valid type name follows
		// Put the token back and let the normal error handling in eatStartNode handle it
		p.TokenStreamPointer--
		return false
	}

	// Collect type name parts - accept identifiers and numeric literals
	for {
		token := p.peekTokenRaw()
		if token == nil {
			break
		}

		// Accept IDENTIFIER or numeric literals as type name parts
		if p.peekToken(IDENTIFIER) || p.peekToken(LITERAL_INT) || p.peekToken(LITERAL_LONG) ||
			p.peekToken(LITERAL_HEXINT) || p.peekToken(LITERAL_HEXLONG) {
			identToken := p.takeToken()
			typeParts = append(typeParts, identToken.StringValue())
		} else {
			break
		}

		if p.peekToken(DOT) {
			p.takeToken() // consume '.'
			continue
		}
		break
	}

	if len(typeParts) == 0 {
		// Not a valid constructor, put back tokens
		p.TokenStreamPointer = newToken.StartPos
		return false
	}

	typeName := strings.Join(typeParts, ".")

	// Create QualifiedIdentifier for the type name - always use QualifiedIdentifier to match Java behavior
	var qualifierNode SpelNode
	qualifierNode = NewQualifiedIdentifier(typeParts, 0, 0) // Position will be updated

	// Check for array notation []{...}
	if p.peekToken(LSQUARE) {
		return p.maybeEatArrayConstructor(newToken, typeName)
	}

	// Regular constructor call with parentheses
	if !p.peekToken(LPAREN) {
		// Not a constructor call, put back tokens
		p.TokenStreamPointer = newToken.StartPos
		return false
	}

	p.takeToken() // consume '('

	// Parse constructor arguments
	var arguments []SpelNode

	// Handle empty argument list
	if p.peekToken(RPAREN) {
		p.takeToken() // consume ')'
		endPos := p.TokenStream[p.TokenStreamPointer-1].EndPos
		constructor := NewConstructorReference(typeName, qualifierNode, arguments, newToken.StartPos, endPos)
		p.push(constructor)
		return true
	}

	// Parse arguments separated by commas
	for {
		arg, err := p.eatExpression()
		if err != nil {
			return false
		}
		if arg == nil {
			return false
		}
		arguments = append(arguments, arg)

		if p.peekToken(COMMA) {
			p.takeToken() // consume ','
			continue
		}
		break
	}

	if !p.peekToken(RPAREN) {
		return false
	}

	endToken := p.takeToken() // consume ')'
	endPos := endToken.EndPos

	constructor := NewConstructorReference(typeName, qualifierNode, arguments, newToken.StartPos, endPos)
	p.push(constructor)
	return true
}

// maybeEatArrayConstructor handles array constructor syntax like new String[]{1,2,3} or new int[][]{{1,2},{3,4}}
// Also handles array constructors with size expressions like new int[1024 * 1024][1024 * 1024]
// This should be parsed as ConstructorReference with InlineList, not ArrayConstructor
func (p *InternalSpelExpressionParser) maybeEatArrayConstructor(newToken *Token, typeName string) bool {
	if !p.peekToken(LSQUARE) {
		return false
	}

	// Try to parse array constructor with size expressions first
	if p.tryParseArrayConstructorWithSizes(newToken, typeName) {
		return true
	}

	// Count and consume array dimensions [][]...
	dimensionCount := 0
	for p.peekToken(LSQUARE) {
		p.takeToken() // consume '['

		if !p.peekToken(RSQUARE) {
			// Not an empty bracket pair, not an array constructor
			return false
		}

		p.takeToken() // consume ']'
		dimensionCount++
	}

	if dimensionCount == 0 || !p.peekToken(LCURLY) {
		return false
	}

	// Build array type name with dimensions (e.g., "int" -> "int[][]")
	arrayTypeName := typeName
	for i := 0; i < dimensionCount; i++ {
		arrayTypeName += "[]"
	}

	// Parse the inline list {1,2,3}
	if !p.peekToken(LCURLY) {
		return false
	}

	startToken := p.takeToken() // consume '{'

	var elements []SpelNode

	// Handle empty array
	if p.peekToken(RCURLY) {
		p.takeToken() // consume '}'
	} else {
		// Parse elements separated by commas
		for {
			element, err := p.eatExpression()
			if err != nil {
				return false
			}
			if element == nil {
				return false
			}
			elements = append(elements, element)

			if p.peekToken(COMMA) {
				p.takeToken() // consume ','
				continue
			}
			break
		}

		if !p.peekToken(RCURLY) {
			return false
		}
		p.takeToken() // consume '}'
	}

	endToken := p.TokenStream[p.TokenStreamPointer-1]
	endPos := endToken.EndPos

	// Create InlineList for the array elements
	inlineList := NewInlineList(elements, startToken.StartPos, endToken.EndPos)

	// Create type identifier - always use QualifiedIdentifier to match Java behavior
	var typeIdentifier SpelNode
	if strings.Contains(typeName, ".") {
		// Qualified identifier like "java.lang.String"
		parts := strings.Split(typeName, ".")
		typeIdentifier = NewQualifiedIdentifier(parts, newToken.StartPos, newToken.EndPos)
	} else {
		// Simple identifier like "String" - still wrap in QualifiedIdentifier
		parts := []string{typeName}
		typeIdentifier = NewQualifiedIdentifier(parts, newToken.StartPos, newToken.EndPos)
	}

	// Create ConstructorReference with array type and InlineList as arguments
	arguments := []SpelNode{inlineList}
	constructorRef := NewConstructorReference(arrayTypeName, typeIdentifier, arguments, newToken.StartPos, endPos)
	p.push(constructorRef)
	return true
}

// tryParseArrayConstructorWithSizes handles array constructor syntax with size expressions
// like new int[1024 * 1024][1024 * 1024] or new char[7]{'a','c','d','e'}
func (p *InternalSpelExpressionParser) tryParseArrayConstructorWithSizes(newToken *Token, typeName string) bool {
	// Save current position in case we need to backtrack
	originalPosition := p.TokenStreamPointer

	var sizeExpressions []SpelNode

	// Parse array dimensions with size expressions [expr][expr]...
	for p.peekToken(LSQUARE) {
		p.takeToken() // consume '['

		if p.peekToken(RSQUARE) {
			// Empty bracket [], revert and let the original function handle it
			p.TokenStreamPointer = originalPosition
			return false
		}

		// Parse the size expression inside brackets
		sizeExpr, err := p.eatExpression()
		if err != nil || sizeExpr == nil {
			// Failed to parse size expression, revert
			p.TokenStreamPointer = originalPosition
			return false
		}

		if !p.peekToken(RSQUARE) {
			// Missing closing bracket, revert
			p.TokenStreamPointer = originalPosition
			return false
		}

		p.takeToken() // consume ']'
		sizeExpressions = append(sizeExpressions, sizeExpr)
	}

	if len(sizeExpressions) == 0 {
		// No size expressions found, revert
		p.TokenStreamPointer = originalPosition
		return false
	}

	// Create type identifier
	var typeIdentifier SpelNode
	if strings.Contains(typeName, ".") {
		// Qualified identifier like "java.lang.String"
		parts := strings.Split(typeName, ".")
		typeIdentifier = NewQualifiedIdentifier(parts, newToken.StartPos, newToken.EndPos)
	} else {
		// Simple identifier like "String" - still wrap in QualifiedIdentifier
		parts := []string{typeName}
		typeIdentifier = NewQualifiedIdentifier(parts, newToken.StartPos, newToken.EndPos)
	}

	// In Java Spring SpEL, size expressions are parsed but not retained in the AST
	// Only the type identifier and optional initializer list are kept as children
	var arguments []SpelNode

	// Check if there's also an initializer list like {'a','c','d','e'}
	if p.peekToken(LCURLY) {
		startToken := p.takeToken() // consume '{'

		var elements []SpelNode

		// Handle empty initializer
		if p.peekToken(RCURLY) {
			p.takeToken() // consume '}'
		} else {
			// Parse elements separated by commas
			for {
				element, err := p.eatExpression()
				if err != nil || element == nil {
					// Failed to parse element, revert
					p.TokenStreamPointer = originalPosition
					return false
				}
				elements = append(elements, element)

				if p.peekToken(COMMA) {
					p.takeToken() // consume ','
					continue
				}
				break
			}

			if !p.peekToken(RCURLY) {
				// Missing closing brace, revert
				p.TokenStreamPointer = originalPosition
				return false
			}
			p.takeToken() // consume '}'
		}

		endToken := p.TokenStream[p.TokenStreamPointer-1]

		// Create InlineList for the array elements and add to arguments
		inlineList := NewInlineList(elements, startToken.StartPos, endToken.EndPos)
		arguments = append(arguments, inlineList)
	}

	endPos := p.TokenStream[p.TokenStreamPointer-1].EndPos

	// Generate display format for array constructor with size expressions
	var dimensionStrings []string
	for _, sizeExpr := range sizeExpressions {
		dimensionStrings = append(dimensionStrings, fmt.Sprintf("[%s]", sizeExpr.ToStringAST()))
	}
	displayFormat := fmt.Sprintf("new %s%s", typeName, strings.Join(dimensionStrings, ""))

	// Add initializer list to display format if present
	if len(arguments) > 0 {
		if inlineList, ok := arguments[len(arguments)-1].(*InlineList); ok {
			displayFormat = fmt.Sprintf("new %s[] %s", typeName, inlineList.ToStringAST())
		}
	}

	// Create ConstructorReference with custom display format
	// Size expressions are used for parsing but not retained in AST (matching Java Spring SpEL behavior)
	constructorRef := NewConstructorReferenceWithDisplay(typeName, typeIdentifier, arguments, displayFormat, newToken.StartPos, endPos)
	p.push(constructorRef)
	return true
}

// maybeEatTypeReference handles type reference syntax like T(java.lang.String)
func (p *InternalSpelExpressionParser) maybeEatTypeReference() bool {
	if !p.peekIdentifierToken("T") {
		return false
	}

	tToken := p.takeToken() // consume 'T'

	if !p.peekToken(LPAREN) {
		// Not a type reference, put back the token
		p.TokenStreamPointer--
		return false
	}

	p.takeToken() // consume '('

	// Parse the qualified type name (may include dots like java.lang.String)
	var typeParts []string

	if !p.peekToken(IDENTIFIER) {
		// Not a valid type reference, put back tokens
		p.TokenStreamPointer = tToken.StartPos
		return false
	}

	// Collect type name parts
	for {
		if !p.peekToken(IDENTIFIER) {
			break
		}
		identToken := p.takeToken()
		typeParts = append(typeParts, identToken.StringValue())

		if p.peekToken(DOT) {
			p.takeToken() // consume '.'
			continue
		}
		break
	}

	if len(typeParts) == 0 {
		// Not a valid type reference, put back tokens
		p.TokenStreamPointer = tToken.StartPos
		return false
	}

	if !p.peekToken(RPAREN) {
		// Not a valid type reference, put back tokens
		p.TokenStreamPointer = tToken.StartPos
		return false
	}

	endToken := p.takeToken() // consume ')'
	typeName := strings.Join(typeParts, ".")

	typeRef := NewTypeReference(typeName, tToken.StartPos, endToken.EndPos)
	p.push(typeRef)
	return true
}

// SpelExpression represents a parsed SpEL expression ready for evaluation
type SpelExpression struct {
	ExpressionString string
	AST              SpelNode
	Configuration    *SpelParserConfiguration
}

func NewSpelExpression(expression string, ast SpelNode, config *SpelParserConfiguration) *SpelExpression {
	return &SpelExpression{
		ExpressionString: expression,
		AST:              ast,
		Configuration:    config,
	}
}

func (expr *SpelExpression) GetValue() (interface{}, error) {
	state := NewExpressionState(expr.Configuration)
	return expr.AST.GetValue(state)
}

func (expr *SpelExpression) GetValueWithRoot(rootObject interface{}) (interface{}, error) {
	state := NewExpressionStateWithRoot(expr.Configuration, NewTypedValue(rootObject))
	return expr.AST.GetValue(state)
}

func (expr *SpelExpression) ToStringAST() string {
	return expr.AST.ToStringAST()
}

func (expr *SpelExpression) GetExpressionString() string {
	return expr.ExpressionString
}

// SpelExpressionParser is the main entry point for parsing SpEL expressions
type SpelExpressionParser struct {
	Configuration *SpelParserConfiguration
}

func NewSpelExpressionParser() *SpelExpressionParser {
	return &SpelExpressionParser{
		Configuration: NewSpelParserConfiguration(),
	}
}

func NewSpelExpressionParserWithConfig(config *SpelParserConfiguration) *SpelExpressionParser {
	return &SpelExpressionParser{
		Configuration: config,
	}
}

func (parser *SpelExpressionParser) ParseExpression(expressionString string) (*SpelExpression, error) {
	return parser.ParseExpressionWithContext(expressionString, nil)
}

// ParseExpressionWithContext parses expression with optional ParserContext (matching Java version)
func (parser *SpelExpressionParser) ParseExpressionWithContext(expressionString string, context *ParserContext) (*SpelExpression, error) {
	if expressionString == "" {
		return nil, fmt.Errorf("'expressionString' must not be null or blank")
	}

	if context != nil && context.IsTemplate {
		return parser.parseTemplate(expressionString, context)
	} else {
		return parser.doParseExpression(expressionString, context)
	}
}

// DoParseExpression parses expression with debug output (matching Java version)
func (parser *SpelExpressionParser) DoParseExpression(expressionString string) (*SpelExpression, error) {
	internalParser := NewInternalSpelExpressionParser(parser.Configuration)
	return internalParser.DoParseExpression(expressionString)
}

// doParseExpression parses expression without debug output
func (parser *SpelExpressionParser) doParseExpression(expressionString string, context *ParserContext) (*SpelExpression, error) {
	internalParser := NewInternalSpelExpressionParser(parser.Configuration)
	return internalParser.ParseExpression(expressionString)
}

// parseTemplate parses template expressions with embedded SpEL expressions
func (parser *SpelExpressionParser) parseTemplate(expressionString string, context *ParserContext) (*SpelExpression, error) {
	if context == nil {
		context = NewTemplateParserContext()
	}

	templateParts, err := parser.parseTemplateExpression(expressionString, context)
	if err != nil {
		return nil, err
	}

	if len(templateParts) == 1 && templateParts[0].IsLiteral {
		// Pure literal template, no expressions
		literal := NewStringLiteral(templateParts[0].Content, 0, len(expressionString))
		return NewSpelExpression(expressionString, literal, parser.Configuration), nil
	}

	// Create composite expression for template
	var nodes []SpelNode
	for _, part := range templateParts {
		if part.IsLiteral {
			// Literal text part
			literal := NewStringLiteral(part.Content, part.StartPos, part.EndPos)
			nodes = append(nodes, literal)
		} else {
			// SpEL expression part
			internalParser := NewInternalSpelExpressionParser(parser.Configuration)
			exprNode, err := internalParser.ParseExpression(part.Content)
			if err != nil {
				return nil, fmt.Errorf("failed to parse expression '%s' in template: %v", part.Content, err)
			}
			nodes = append(nodes, exprNode.AST)
		}
	}

	// Create template expression
	templateExpr := NewTemplateExpression(nodes, 0, len(expressionString))
	return NewSpelExpression(expressionString, templateExpr, parser.Configuration), nil
}

func (parser *SpelExpressionParser) ParseAST() (SpelNode, error) {
	internalParser := NewInternalSpelExpressionParser(parser.Configuration)
	return internalParser.eatExpression()
}

// TemplatePart represents a part of a template (either literal text or SpEL expression)
type TemplatePart struct {
	Content   string
	IsLiteral bool
	StartPos  int
	EndPos    int
}

// parseTemplateExpression parses a template string into literal and expression parts
func (parser *SpelExpressionParser) parseTemplateExpression(template string, context *ParserContext) ([]TemplatePart, error) {
	var parts []TemplatePart
	var currentPos int

	prefixLen := len(context.ExpressionPrefix)
	suffixLen := len(context.ExpressionSuffix)

	for currentPos < len(template) {
		// Look for expression prefix
		exprStart := strings.Index(template[currentPos:], context.ExpressionPrefix)

		if exprStart == -1 {
			// No more expressions, rest is literal
			if currentPos < len(template) {
				literal := template[currentPos:]
				parts = append(parts, TemplatePart{
					Content:   literal,
					IsLiteral: true,
					StartPos:  currentPos,
					EndPos:    len(template),
				})
			}
			break
		}

		// Adjust position to absolute
		exprStart += currentPos

		// Add literal part before expression (if any)
		if exprStart > currentPos {
			literal := template[currentPos:exprStart]
			parts = append(parts, TemplatePart{
				Content:   literal,
				IsLiteral: true,
				StartPos:  currentPos,
				EndPos:    exprStart,
			})
		}

		// Look for expression suffix
		exprContentStart := exprStart + prefixLen
		exprEnd := strings.Index(template[exprContentStart:], context.ExpressionSuffix)

		if exprEnd == -1 {
			return nil, fmt.Errorf("no closing '%s' for expression starting at position %d", context.ExpressionSuffix, exprStart)
		}

		// Adjust position to absolute
		exprEnd += exprContentStart

		// Extract expression content
		exprContent := template[exprContentStart:exprEnd]
		parts = append(parts, TemplatePart{
			Content:   exprContent,
			IsLiteral: false,
			StartPos:  exprStart,
			EndPos:    exprEnd + suffixLen,
		})

		// Move position past the expression
		currentPos = exprEnd + suffixLen
	}

	return parts, nil
}
