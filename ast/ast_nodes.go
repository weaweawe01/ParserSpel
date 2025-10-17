package ast

import (
	"fmt"
	"strconv"
	"strings"
)

// Literal represents a literal value in the expression
type Literal struct {
	*SpelNodeImpl
	Value interface{}
}

func NewLiteral(value interface{}, startPos, endPos int) *Literal {
	return &Literal{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos),
		Value:        value,
	}
}

func (l *Literal) GetValue(state *ExpressionState) (interface{}, error) {
	return l.Value, nil
}

func (l *Literal) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	return NewTypedValue(l.Value), nil
}

func (l *Literal) ToStringAST() string {
	if str, ok := l.Value.(string); ok {
		return fmt.Sprintf("'%s'", str)
	}
	return fmt.Sprintf("%v", l.Value)
}

// StringLiteral represents a string literal
type StringLiteral struct {
	*Literal
}

func NewStringLiteral(value string, startPos, endPos int) *StringLiteral {
	return &StringLiteral{
		Literal: NewLiteral(value, startPos, endPos),
	}
}

// BooleanLiteral represents a boolean literal
type BooleanLiteral struct {
	*Literal
}

func NewBooleanLiteral(value bool, startPos, endPos int) *BooleanLiteral {
	return &BooleanLiteral{
		Literal: NewLiteral(value, startPos, endPos),
	}
}

// NullLiteral represents a null literal
type NullLiteral struct {
	*SpelNodeImpl
}

func NewNullLiteral(startPos, endPos int) *NullLiteral {
	return &NullLiteral{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos),
	}
}

func (n *NullLiteral) GetValue(state *ExpressionState) (interface{}, error) {
	return nil, nil
}

func (n *NullLiteral) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	return NewTypedValue(nil), nil
}

func (n *NullLiteral) ToStringAST() string {
	return "null"
}

// Identifier represents an identifier (variable, property name, etc.)
type Identifier struct {
	*SpelNodeImpl
	Name string
}

func NewIdentifier(name string, startPos, endPos int) *Identifier {
	return &Identifier{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos),
		Name:         name,
	}
}

func (i *Identifier) GetValue(state *ExpressionState) (interface{}, error) {
	// For this simple implementation, return default values for common identifiers
	// In a real implementation, this would look up the identifier value from context
	switch strings.ToLower(i.Name) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "null":
		return nil, nil
	default:
		// For unknown identifiers, return the name as a string
		// This prevents parse errors in arithmetic operations
		return i.Name, nil
	}
}

func (i *Identifier) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := i.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (i *Identifier) ToStringAST() string {
	return i.Name
}

// PropertyOrFieldReference represents property or field access
type PropertyOrFieldReference struct {
	*SpelNodeImpl
	Name               string
	NullSafeNavigation bool
	IsDirectReference  bool // true for direct references like "userName", false for chained like ".name"
}

func NewPropertyOrFieldReference(nullSafeNavigation bool, name string, startPos, endPos int) *PropertyOrFieldReference {
	return &PropertyOrFieldReference{
		SpelNodeImpl:       NewSpelNodeImpl(startPos, endPos),
		Name:               name,
		NullSafeNavigation: nullSafeNavigation,
		IsDirectReference:  false, // Default to chained reference
	}
}

// NewDirectPropertyOrFieldReference creates a direct property reference (like "userName")
func NewDirectPropertyOrFieldReference(name string, startPos, endPos int) *PropertyOrFieldReference {
	return &PropertyOrFieldReference{
		SpelNodeImpl:       NewSpelNodeImpl(startPos, endPos),
		Name:               name,
		NullSafeNavigation: false,
		IsDirectReference:  true,
	}
}

func (p *PropertyOrFieldReference) GetValue(state *ExpressionState) (interface{}, error) {
	// Placeholder implementation
	return p.Name, nil
}

func (p *PropertyOrFieldReference) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := p.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (p *PropertyOrFieldReference) ToStringAST() string {
	if p.IsDirectReference {
		// Direct reference like "userName" - no dot prefix
		return p.Name
	}
	if p.NullSafeNavigation {
		return "?." + p.Name
	}
	return "." + p.Name
}

// CompoundExpression represents a compound expression with multiple parts
type CompoundExpression struct {
	*SpelNodeImpl
}

func NewCompoundExpression(startPos, endPos int, children ...SpelNode) *CompoundExpression {
	return &CompoundExpression{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos, children...),
	}
}

func (c *CompoundExpression) GetValue(state *ExpressionState) (interface{}, error) {
	if len(c.Children) == 0 {
		return nil, nil
	}

	var result interface{}
	var err error

	for _, child := range c.Children {
		result, err = child.GetValue(state)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (c *CompoundExpression) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := c.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (c *CompoundExpression) ToStringAST() string {
	if len(c.Children) == 0 {
		return ""
	}

	var result strings.Builder
	for i, child := range c.Children {
		childStr := child.ToStringAST()

		if i == 0 {
			// First child - use as is
			result.WriteString(childStr)
		} else {
			// For subsequent children, add appropriate connector
			if _, isPropertyRef := child.(*PropertyOrFieldReference); isPropertyRef {
				// PropertyOrFieldReference handles its own dot prefix
				result.WriteString(childStr)
			} else if _, isMethodRef := child.(*MethodReference); isMethodRef {
				// MethodReference needs a dot prefix in compound expressions
				if !strings.HasPrefix(childStr, ".") && !strings.HasPrefix(childStr, "?.") {
					result.WriteString(".")
				}
				result.WriteString(childStr)
			} else {
				// Other node types
				result.WriteString(childStr)
			}
		}
	}
	return result.String()
} // VariableReference represents a variable reference (#var)
type VariableReference struct {
	*SpelNodeImpl
	Name string
}

func NewVariableReference(name string, startPos, endPos int) *VariableReference {
	return &VariableReference{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos),
		Name:         name,
	}
}

func (v *VariableReference) GetValue(state *ExpressionState) (interface{}, error) {
	// Placeholder implementation
	return "#" + v.Name, nil
}

func (v *VariableReference) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := v.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (v *VariableReference) ToStringAST() string {
	return "#" + v.Name
}

// BeanReference represents a bean reference (@bean)
type BeanReference struct {
	*SpelNodeImpl
	Name string
}

func NewBeanReference(name string, startPos, endPos int) *BeanReference {
	return &BeanReference{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos),
		Name:         name,
	}
}

func (b *BeanReference) GetValue(state *ExpressionState) (interface{}, error) {
	// Placeholder implementation
	return "@" + b.Name, nil
}

func (b *BeanReference) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := b.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (b *BeanReference) ToStringAST() string {
	return "@" + b.Name
}

// parseNumber parses a string token into appropriate numeric type
func parseNumber(tokenData string, tokenKind TokenKind) (interface{}, error) {
	switch tokenKind {
	case LITERAL_INT:
		return strconv.Atoi(tokenData)
	case LITERAL_LONG:
		// Remove 'L' suffix
		data := strings.TrimSuffix(tokenData, "L")
		data = strings.TrimSuffix(data, "l")
		return strconv.ParseInt(data, 10, 64)
	case LITERAL_HEXINT:
		return strconv.ParseInt(tokenData, 16, 32)
	case LITERAL_HEXLONG:
		// Remove 'L' suffix
		data := strings.TrimSuffix(tokenData, "L")
		data = strings.TrimSuffix(data, "l")
		return strconv.ParseInt(data, 16, 64)
	case LITERAL_REAL:
		return strconv.ParseFloat(tokenData, 64)
	case LITERAL_REAL_FLOAT:
		// Remove 'F' suffix
		data := strings.TrimSuffix(tokenData, "F")
		data = strings.TrimSuffix(data, "f")
		return strconv.ParseFloat(data, 32)
	default:
		return nil, fmt.Errorf("unsupported numeric token kind: %v", tokenKind)
	}
}

// MethodReference represents a method call in the SpEL AST
type MethodReference struct {
	*SpelNodeImpl
	Name      string
	NullSafe  bool
	Arguments []SpelNode
}

func NewMethodReference(nullSafe bool, name string, arguments []SpelNode, startPos, endPos int) *MethodReference {
	return &MethodReference{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos, arguments...),
		Name:         name,
		NullSafe:     nullSafe,
		Arguments:    arguments,
	}
}

func (m *MethodReference) GetValue(state *ExpressionState) (interface{}, error) {
	// Placeholder implementation
	var argValues []interface{}
	for _, arg := range m.Arguments {
		val, err := arg.GetValue(state)
		if err != nil {
			return nil, err
		}
		argValues = append(argValues, val)
	}
	return fmt.Sprintf("methodCall(%s, %v)", m.Name, argValues), nil
}

func (m *MethodReference) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := m.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (m *MethodReference) ToStringAST() string {
	if len(m.Arguments) == 0 {
		return fmt.Sprintf("%s()", m.Name)
	}

	var argStrings []string
	for _, arg := range m.Arguments {
		argStrings = append(argStrings, arg.ToStringAST())
	}
	return fmt.Sprintf("%s(%s)", m.Name, strings.Join(argStrings, ", "))
}

// ConstructorReference represents a constructor call in the SpEL AST
type ConstructorReference struct {
	*SpelNodeImpl
	TypeName      string
	QualifierNode SpelNode // The qualified identifier for the type
	Arguments     []SpelNode
}

func NewConstructorReference(typeName string, qualifierNode SpelNode, arguments []SpelNode, startPos, endPos int) *ConstructorReference {
	// Create children array including the qualifier and arguments
	var children []SpelNode
	if qualifierNode != nil {
		children = append(children, qualifierNode)
	}
	children = append(children, arguments...)

	return &ConstructorReference{
		SpelNodeImpl:  NewSpelNodeImpl(startPos, endPos, children...),
		TypeName:      typeName,
		QualifierNode: qualifierNode,
		Arguments:     arguments,
	}
}

func (c *ConstructorReference) GetValue(state *ExpressionState) (interface{}, error) {
	// Placeholder implementation
	var argValues []interface{}
	for _, arg := range c.Arguments {
		val, err := arg.GetValue(state)
		if err != nil {
			return nil, err
		}
		argValues = append(argValues, val)
	}
	return fmt.Sprintf("new %s(%v)", c.TypeName, argValues), nil
}

func (c *ConstructorReference) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := c.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (c *ConstructorReference) ToStringAST() string {
	if len(c.Arguments) == 0 {
		return fmt.Sprintf("new %s()", c.TypeName)
	}

	var argStrings []string
	for _, arg := range c.Arguments {
		argStrings = append(argStrings, arg.ToStringAST())
	}
	return fmt.Sprintf("new %s(%s)", c.TypeName, strings.Join(argStrings, ", "))
}

// ArrayConstructor represents an array constructor like new byte[]{1,2,3}
type ArrayConstructor struct {
	*SpelNodeImpl
	TypeName string
	Elements []SpelNode
}

func NewArrayConstructor(typeName string, elements []SpelNode, startPos, endPos int) *ArrayConstructor {
	return &ArrayConstructor{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos, elements...),
		TypeName:     typeName,
		Elements:     elements,
	}
}

func (a *ArrayConstructor) GetValue(state *ExpressionState) (interface{}, error) {
	// Placeholder implementation
	var elementValues []interface{}
	for _, element := range a.Elements {
		val, err := element.GetValue(state)
		if err != nil {
			return nil, err
		}
		elementValues = append(elementValues, val)
	}
	return fmt.Sprintf("new %s[]{%v}", a.TypeName, elementValues), nil
}

func (a *ArrayConstructor) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := a.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (a *ArrayConstructor) ToStringAST() string {
	if len(a.Elements) == 0 {
		return fmt.Sprintf("new %s[]{}", a.TypeName)
	}

	var elementStrings []string
	for _, element := range a.Elements {
		elementStrings = append(elementStrings, element.ToStringAST())
	}
	return fmt.Sprintf("new %s[]{%s}", a.TypeName, strings.Join(elementStrings, ", "))
}

// TemplateExpression represents a template expression with mixed literal and SpEL parts
type TemplateExpression struct {
	*SpelNodeImpl
	Parts []SpelNode
}

func NewTemplateExpression(parts []SpelNode, startPos, endPos int) *TemplateExpression {
	return &TemplateExpression{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos, parts...),
		Parts:        parts,
	}
}

func (t *TemplateExpression) GetValue(state *ExpressionState) (interface{}, error) {
	var result strings.Builder

	for _, part := range t.Parts {
		value, err := part.GetValue(state)
		if err != nil {
			return nil, err
		}
		result.WriteString(fmt.Sprintf("%v", value))
	}

	return result.String(), nil
}

func (t *TemplateExpression) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := t.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (t *TemplateExpression) ToStringAST() string {
	var parts []string
	for _, part := range t.Parts {
		parts = append(parts, part.ToStringAST())
	}
	return fmt.Sprintf("template[%s]", strings.Join(parts, " + "))
}

// TypeReference represents a type reference like T(java.lang.String)
type TypeReference struct {
	*SpelNodeImpl
	TypeName            string
	QualifiedIdentifier *QualifiedIdentifier
}

func NewTypeReference(typeName string, startPos, endPos int) *TypeReference {
	// Split the type name into qualifiers
	qualifiers := strings.Split(typeName, ".")

	// Create a QualifiedIdentifier for the type name
	// Position calculation: T( = 2 chars, then the qualified identifier, then )
	qualifierStartPos := startPos + 2
	qualifierEndPos := endPos - 1
	qualifiedId := NewQualifiedIdentifier(qualifiers, qualifierStartPos, qualifierEndPos)

	return &TypeReference{
		SpelNodeImpl:        NewSpelNodeImpl(startPos, endPos, qualifiedId),
		TypeName:            typeName,
		QualifiedIdentifier: qualifiedId,
	}
}

func (t *TypeReference) GetValue(state *ExpressionState) (interface{}, error) {
	// Placeholder implementation - type references typically return Class objects
	return fmt.Sprintf("Class<%s>", t.TypeName), nil
}

func (t *TypeReference) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := t.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (t *TypeReference) ToStringAST() string {
	return fmt.Sprintf("T(%s)", t.TypeName)
}

// QualifiedIdentifier represents a qualified identifier like java.lang.String
type QualifiedIdentifier struct {
	*SpelNodeImpl
	Qualifiers []string
}

func NewQualifiedIdentifier(qualifiers []string, startPos, endPos int) *QualifiedIdentifier {
	// Create child identifier nodes for each qualifier
	var children []SpelNode
	currentPos := startPos
	for i, qualifier := range qualifiers {
		// Calculate positions for each identifier
		identifierEndPos := currentPos + len(qualifier)
		identifier := NewIdentifier(qualifier, currentPos, identifierEndPos)
		children = append(children, identifier)

		// Move position past the identifier and dot (if not last)
		if i < len(qualifiers)-1 {
			currentPos = identifierEndPos + 1 // +1 for the dot
		}
	}

	return &QualifiedIdentifier{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos, children...),
		Qualifiers:   qualifiers,
	}
}

func (q *QualifiedIdentifier) GetValue(state *ExpressionState) (interface{}, error) {
	// Return the full qualified name
	return strings.Join(q.Qualifiers, "."), nil
}

func (q *QualifiedIdentifier) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := q.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (q *QualifiedIdentifier) ToStringAST() string {
	return strings.Join(q.Qualifiers, ".")
}
