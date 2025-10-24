package ast

import (
	"fmt"
	"strconv"
	"strings"
)

// IntLiteral represents an integer literal value in the expression
type IntLiteral struct {
	*SpelNodeImpl
	Value interface{}
}

func NewIntLiteral(value interface{}, startPos, endPos int) *IntLiteral {
	return &IntLiteral{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos),
		Value:        value,
	}
}

func (l *IntLiteral) GetValue(state *ExpressionState) (interface{}, error) {
	return l.Value, nil
}

func (l *IntLiteral) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	return NewTypedValue(l.Value), nil
}

func (l *IntLiteral) ToStringAST() string {
	if str, ok := l.Value.(string); ok {
		return fmt.Sprintf("'%s'", str)
	}
	return fmt.Sprintf("%v", l.Value)
}

// StringLiteral represents a string literal
type StringLiteral struct {
	*IntLiteral
}

func NewStringLiteral(value string, startPos, endPos int) *StringLiteral {
	return &StringLiteral{
		IntLiteral: NewIntLiteral(value, startPos, endPos),
	}
}

// BooleanLiteral represents a boolean literal
type BooleanLiteral struct {
	*IntLiteral
}

func NewBooleanLiteral(value bool, startPos, endPos int) *BooleanLiteral {
	return &BooleanLiteral{
		IntLiteral: NewIntLiteral(value, startPos, endPos),
	}
}

// RealLiteral represents a real (floating point) literal
type RealLiteral struct {
	*IntLiteral
}

func NewRealLiteral(value float64, startPos, endPos int) *RealLiteral {
	return &RealLiteral{
		IntLiteral: NewIntLiteral(value, startPos, endPos),
	}
}

func (r *RealLiteral) ToStringAST() string {
	// Format float64 values with at least one decimal place
	if val, ok := r.Value.(float64); ok {
		// If it's a whole number, add .0
		if val == float64(int64(val)) {
			return fmt.Sprintf("%.1f", val)
		}
		return fmt.Sprintf("%g", val)
	}
	return fmt.Sprintf("%v", r.Value)
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
		// Remove 'd' or 'D' suffix for double literals
		data := strings.TrimSuffix(tokenData, "D")
		data = strings.TrimSuffix(data, "d")
		return strconv.ParseFloat(data, 64)
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
	DisplayFormat string // Optional custom display format for array constructors
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
		DisplayFormat: "", // Default empty, will use auto-generated format
	}
}

func NewConstructorReferenceWithDisplay(typeName string, qualifierNode SpelNode, arguments []SpelNode, displayFormat string, startPos, endPos int) *ConstructorReference {
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
		DisplayFormat: displayFormat,
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
	// Use custom display format if provided
	if c.DisplayFormat != "" {
		return c.DisplayFormat
	}

	if len(c.Arguments) == 0 {
		return fmt.Sprintf("new %s()", c.TypeName)
	}

	// Check if this is an array constructor (single InlineList argument)
	if len(c.Arguments) == 1 {
		if inlineList, ok := c.Arguments[0].(*InlineList); ok {
			// Format as array constructor
			// If TypeName already contains [], use as-is; otherwise add []
			if strings.Contains(c.TypeName, "[]") {
				return fmt.Sprintf("new %s %s", c.TypeName, inlineList.ToStringAST())
			} else {
				return fmt.Sprintf("new %s[] %s", c.TypeName, inlineList.ToStringAST())
			}
		}
	}

	// Check if this is an array constructor with size and initializer
	// e.g., new char[7]{'a','c','d','e'} -> arguments: [Literal(7), InlineList]
	if len(c.Arguments) >= 2 {
		lastArg := c.Arguments[len(c.Arguments)-1]
		if inlineList, ok := lastArg.(*InlineList); ok {
			// This looks like an array constructor with size and initializer
			// Format as: new type[] {elements}
			return fmt.Sprintf("new %s[] %s", c.TypeName, inlineList.ToStringAST())
		}
	}

	// This logic should only be triggered by the custom DisplayFormat from parser
	// Normal constructor calls should fall through to the default formatting

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

// InlineList represents a list literal like {1,2,3}
type InlineList struct {
	*SpelNodeImpl
	Elements []SpelNode
}

func NewInlineList(elements []SpelNode, startPos, endPos int) *InlineList {
	return &InlineList{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos, elements...),
		Elements:     elements,
	}
}

func (i *InlineList) GetValue(state *ExpressionState) (interface{}, error) {
	// Return a slice of element values
	var elementValues []interface{}
	for _, element := range i.Elements {
		val, err := element.GetValue(state)
		if err != nil {
			return nil, err
		}
		elementValues = append(elementValues, val)
	}
	return elementValues, nil
}

func (i *InlineList) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	val, err := i.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(val), nil
}

func (i *InlineList) ToStringAST() string {
	var elements []string
	for _, element := range i.Elements {
		elements = append(elements, element.ToStringAST())
	}
	return "{" + strings.Join(elements, ",") + "}"
}

// Indexer represents indexing operations like array[index] or list[index]
type Indexer struct {
	*SpelNodeImpl
	IndexExpression SpelNode
	NullSafe        bool // true for ?.[...], false for [...]
}

func NewIndexer(indexExpression SpelNode, startPos, endPos int) *Indexer {
	children := []SpelNode{}
	if indexExpression != nil {
		children = append(children, indexExpression)
	}
	return &Indexer{
		SpelNodeImpl:    NewSpelNodeImpl(startPos, endPos, children...),
		IndexExpression: indexExpression,
		NullSafe:        false,
	}
}

func NewNullSafeIndexer(indexExpression SpelNode, startPos, endPos int) *Indexer {
	children := []SpelNode{}
	if indexExpression != nil {
		children = append(children, indexExpression)
	}
	return &Indexer{
		SpelNodeImpl:    NewSpelNodeImpl(startPos, endPos, children...),
		IndexExpression: indexExpression,
		NullSafe:        true,
	}
}

func (i *Indexer) GetValue(state *ExpressionState) (interface{}, error) {
	// Placeholder implementation - would need context about the indexed object
	return nil, fmt.Errorf("indexer evaluation not yet implemented")
}

func (i *Indexer) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	val, err := i.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(val), nil
}

func (i *Indexer) ToStringAST() string {
	prefix := ""
	if i.NullSafe {
		prefix = "?."
	}
	if i.IndexExpression != nil {
		return prefix + "[" + i.IndexExpression.ToStringAST() + "]"
	}
	return prefix + "[]"
}

// Assign represents assignment expressions like variable = value
type Assign struct {
	*SpelNodeImpl
	Left  SpelNode // Left-hand side (usually PropertyOrFieldReference)
	Right SpelNode // Right-hand side (value to assign)
}

func NewAssign(left, right SpelNode, startPos, endPos int) *Assign {
	children := []SpelNode{}
	if left != nil {
		children = append(children, left)
	}
	if right != nil {
		children = append(children, right)
	}
	return &Assign{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos, children...),
		Left:         left,
		Right:        right,
	}
}

func (a *Assign) GetValue(state *ExpressionState) (interface{}, error) {
	// Placeholder implementation - would set the value and return it
	rightValue, err := a.Right.GetValue(state)
	if err != nil {
		return nil, err
	}
	// In a real implementation, this would assign the value to the left-hand side
	return rightValue, nil
}

func (a *Assign) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	val, err := a.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(val), nil
}

func (a *Assign) ToStringAST() string {
	left := ""
	if a.Left != nil {
		left = a.Left.ToStringAST()
	}
	right := ""
	if a.Right != nil {
		right = a.Right.ToStringAST()
	}
	return left + " = " + right
}

// Selection represents selection expressions like collection.?[criteria]
type SelectionKind int

const (
	SelectionAll   SelectionKind = iota // ?[...] - select all matching
	SelectionFirst                      // ^[...] - select first matching
	SelectionLast                       // $[...] - select last matching
)

type Selection struct {
	*SpelNodeImpl
	NullSafe bool          // true for safe navigation (?.)
	Kind     SelectionKind // Type of selection (all, first, last)
	Criteria SpelNode      // The selection criteria expression
}

func NewSelection(nullSafe bool, kind SelectionKind, criteria SpelNode, startPos, endPos int) *Selection {
	children := []SpelNode{}
	if criteria != nil {
		children = append(children, criteria)
	}
	return &Selection{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos, children...),
		NullSafe:     nullSafe,
		Kind:         kind,
		Criteria:     criteria,
	}
}

func (s *Selection) GetValue(state *ExpressionState) (interface{}, error) {
	// Placeholder implementation - would filter the collection based on criteria
	return nil, fmt.Errorf("selection evaluation not yet implemented")
}

func (s *Selection) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	val, err := s.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(val), nil
}

func (s *Selection) ToStringAST() string {
	criteriaStr := ""
	if s.Criteria != nil {
		criteriaStr = s.Criteria.ToStringAST()
	}

	prefix := ""
	if s.NullSafe {
		prefix = "?."
	}

	switch s.Kind {
	case SelectionAll:
		return prefix + "?[" + criteriaStr + "]"
	case SelectionFirst:
		return prefix + "^[" + criteriaStr + "]"
	case SelectionLast:
		return prefix + "$[" + criteriaStr + "]"
	default:
		return prefix + "?[" + criteriaStr + "]"
	}
}

// FunctionReference represents function calls like #functionName(args...)
type FunctionReference struct {
	*SpelNodeImpl
	FunctionName string
	Arguments    []SpelNode
}

func NewFunctionReference(functionName string, arguments []SpelNode, startPos, endPos int) *FunctionReference {
	return &FunctionReference{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos, arguments...),
		FunctionName: functionName,
		Arguments:    arguments,
	}
}

func (f *FunctionReference) GetValue(state *ExpressionState) (interface{}, error) {
	// Placeholder implementation - would call the registered function
	return nil, fmt.Errorf("function evaluation not yet implemented")
}

func (f *FunctionReference) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	val, err := f.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(val), nil
}

func (f *FunctionReference) ToStringAST() string {
	var args []string
	for _, arg := range f.Arguments {
		args = append(args, arg.ToStringAST())
	}
	return "#" + f.FunctionName + "(" + strings.Join(args, ", ") + ")"
}

// Ternary represents ternary conditional expressions like condition ? trueValue : falseValue
type Ternary struct {
	*SpelNodeImpl
	Condition  SpelNode // The condition expression
	TrueValue  SpelNode // Value returned if condition is true
	FalseValue SpelNode // Value returned if condition is false
}

func NewTernary(condition, trueValue, falseValue SpelNode, startPos, endPos int) *Ternary {
	children := []SpelNode{}
	if condition != nil {
		children = append(children, condition)
	}
	if trueValue != nil {
		children = append(children, trueValue)
	}
	if falseValue != nil {
		children = append(children, falseValue)
	}
	return &Ternary{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos, children...),
		Condition:    condition,
		TrueValue:    trueValue,
		FalseValue:   falseValue,
	}
}

func (t *Ternary) GetValue(state *ExpressionState) (interface{}, error) {
	// Evaluate condition
	conditionValue, err := t.Condition.GetValue(state)
	if err != nil {
		return nil, err
	}

	// Convert to boolean
	var conditionBool bool
	switch v := conditionValue.(type) {
	case bool:
		conditionBool = v
	case string:
		conditionBool = v != ""
	case int:
		conditionBool = v != 0
	default:
		conditionBool = v != nil
	}

	// Return appropriate value
	if conditionBool {
		return t.TrueValue.GetValue(state)
	} else {
		return t.FalseValue.GetValue(state)
	}
}

func (t *Ternary) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	val, err := t.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(val), nil
}

func (t *Ternary) ToStringAST() string {
	condition := ""
	if t.Condition != nil {
		condition = t.Condition.ToStringAST()
	}
	trueVal := ""
	if t.TrueValue != nil {
		trueVal = t.TrueValue.ToStringAST()
	}
	falseVal := ""
	if t.FalseValue != nil {
		falseVal = t.FalseValue.ToStringAST()
	}
	return "(" + condition + " ? " + trueVal + " : " + falseVal + ")"
}

// Elvis represents Elvis operator expressions like value ?: defaultValue
type Elvis struct {
	*SpelNodeImpl
	Expression   SpelNode // The expression to evaluate
	DefaultValue SpelNode // Default value if expression is null/empty
}

func NewElvis(expression, defaultValue SpelNode, startPos, endPos int) *Elvis {
	children := []SpelNode{}
	if expression != nil {
		children = append(children, expression)
	}
	if defaultValue != nil {
		children = append(children, defaultValue)
	}
	return &Elvis{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos, children...),
		Expression:   expression,
		DefaultValue: defaultValue,
	}
}

func (e *Elvis) GetValue(state *ExpressionState) (interface{}, error) {
	// Evaluate the main expression
	expressionValue, err := e.Expression.GetValue(state)
	if err != nil {
		return nil, err
	}

	// Check if the value is "empty" (null, empty string, zero, false)
	isEmpty := false
	switch v := expressionValue.(type) {
	case nil:
		isEmpty = true
	case string:
		isEmpty = v == ""
	case int:
		isEmpty = v == 0
	case bool:
		isEmpty = !v
	default:
		isEmpty = expressionValue == nil
	}

	// Return default value if empty, otherwise return the expression value
	if isEmpty {
		return e.DefaultValue.GetValue(state)
	} else {
		return expressionValue, nil
	}
}

func (e *Elvis) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	val, err := e.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(val), nil
}

func (e *Elvis) ToStringAST() string {
	expression := ""
	if e.Expression != nil {
		expression = e.Expression.ToStringAST()
	}
	defaultVal := ""
	if e.DefaultValue != nil {
		defaultVal = e.DefaultValue.ToStringAST()
	}
	return expression + " ?: " + defaultVal
}

// Projection represents projection expressions like collection.![expression]
type Projection struct {
	*SpelNodeImpl
	ProjectionExpression SpelNode // The expression to project from each element
}

func NewProjection(projectionExpression SpelNode, startPos, endPos int) *Projection {
	children := []SpelNode{}
	if projectionExpression != nil {
		children = append(children, projectionExpression)
	}
	return &Projection{
		SpelNodeImpl:         NewSpelNodeImpl(startPos, endPos, children...),
		ProjectionExpression: projectionExpression,
	}
}

func (p *Projection) GetValue(state *ExpressionState) (interface{}, error) {
	// Placeholder implementation - would project the expression over a collection
	return nil, fmt.Errorf("projection evaluation not yet implemented")
}

func (p *Projection) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	val, err := p.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(val), nil
}

func (p *Projection) ToStringAST() string {
	projectionStr := ""
	if p.ProjectionExpression != nil {
		projectionStr = p.ProjectionExpression.ToStringAST()
	}
	return ".![" + projectionStr + "]"
}

// InlineMap represents a map literal like {key1:value1,key2:value2}
type InlineMap struct {
	*SpelNodeImpl
	KeyValuePairs []KeyValuePair
}

type KeyValuePair struct {
	Key   SpelNode
	Value SpelNode
}

func NewInlineMap(pairs []KeyValuePair, startPos, endPos int) *InlineMap {
	var children []SpelNode
	for _, pair := range pairs {
		children = append(children, pair.Key, pair.Value)
	}
	return &InlineMap{
		SpelNodeImpl:  NewSpelNodeImpl(startPos, endPos, children...),
		KeyValuePairs: pairs,
	}
}

func (i *InlineMap) GetValue(state *ExpressionState) (interface{}, error) {
	// Return a map of key-value pairs
	result := make(map[interface{}]interface{})
	for _, pair := range i.KeyValuePairs {
		key, err := pair.Key.GetValue(state)
		if err != nil {
			return nil, err
		}
		value, err := pair.Value.GetValue(state)
		if err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, nil
}

func (i *InlineMap) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	val, err := i.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(val), nil
}

func (i *InlineMap) ToStringAST() string {
	var pairs []string
	for _, pair := range i.KeyValuePairs {
		pairs = append(pairs, pair.Key.ToStringAST()+":"+pair.Value.ToStringAST())
	}
	return "{" + strings.Join(pairs, ",") + "}"
}
