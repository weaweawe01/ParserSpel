package ast

import (
	"fmt"
	"reflect"
	"strings"
)

// SpelNode represents a node in the SpEL Abstract Syntax Tree
type SpelNode interface {
	// GetValue evaluates this node and returns the result
	GetValue(state *ExpressionState) (interface{}, error)

	// GetTypedValue evaluates this node and returns a typed result
	GetTypedValue(state *ExpressionState) (*TypedValue, error)

	// IsWritable returns true if this node can be assigned a value
	IsWritable(state *ExpressionState) bool

	// SetValue assigns a value to this node
	SetValue(state *ExpressionState, value interface{}) error

	// ToStringAST returns the string representation of this AST node
	ToStringAST() string

	// GetStartPosition returns the start position in the expression string
	GetStartPosition() int

	// GetEndPosition returns the end position in the expression string
	GetEndPosition() int

	// GetChildren returns child nodes
	GetChildren() []SpelNode

	// IsCompilable returns true if this node can be compiled
	IsCompilable() bool
}

// SpelNodeImpl provides a base implementation for SpEL nodes
type SpelNodeImpl struct {
	StartPos int
	EndPos   int
	Children []SpelNode
}

// NewSpelNodeImpl creates a new base SpEL node
func NewSpelNodeImpl(startPos, endPos int, children ...SpelNode) *SpelNodeImpl {
	return &SpelNodeImpl{
		StartPos: startPos,
		EndPos:   endPos,
		Children: children,
	}
}

func (n *SpelNodeImpl) GetStartPosition() int {
	return n.StartPos
}

func (n *SpelNodeImpl) GetEndPosition() int {
	return n.EndPos
}

func (n *SpelNodeImpl) GetChildren() []SpelNode {
	return n.Children
}

func (n *SpelNodeImpl) IsCompilable() bool {
	return false // Default implementation
}

func (n *SpelNodeImpl) IsWritable(state *ExpressionState) bool {
	return false // Default implementation
}

func (n *SpelNodeImpl) SetValue(state *ExpressionState, value interface{}) error {
	return fmt.Errorf("cannot set value on %T", n)
}

// TypedValue represents a value with type information
type TypedValue struct {
	Value interface{}
	Type  string
}

func NewTypedValue(value interface{}) *TypedValue {
	return &TypedValue{
		Value: value,
		Type:  fmt.Sprintf("%T", value),
	}
}

// ExpressionState holds the evaluation context and configuration
type ExpressionState struct {
	EvaluationContext interface{} // Placeholder for evaluation context
	Configuration     *SpelParserConfiguration
	RootObject        *TypedValue
}

func NewExpressionState(config *SpelParserConfiguration) *ExpressionState {
	return &ExpressionState{
		Configuration: config,
		RootObject:    NewTypedValue(nil),
	}
}

func NewExpressionStateWithRoot(config *SpelParserConfiguration, rootObject *TypedValue) *ExpressionState {
	return &ExpressionState{
		Configuration: config,
		RootObject:    rootObject,
	}
}

// SpelParserConfiguration holds parser configuration
type SpelParserConfiguration struct {
	MaximumExpressionLength int
	AutoGrowCollections     bool
	AutoGrowNullReferences  bool
}

func NewSpelParserConfiguration() *SpelParserConfiguration {
	return &SpelParserConfiguration{
		MaximumExpressionLength: 10000,
		AutoGrowCollections:     false,
		AutoGrowNullReferences:  false,
	}
}

// ParserContext holds context information for parsing expressions
type ParserContext struct {
	IsTemplate       bool
	ExpressionPrefix string
	ExpressionSuffix string
}

func NewParserContext() *ParserContext {
	return &ParserContext{
		IsTemplate:       false,
		ExpressionPrefix: "#{",
		ExpressionSuffix: "}",
	}
}

func NewTemplateParserContext() *ParserContext {
	return &ParserContext{
		IsTemplate:       true,
		ExpressionPrefix: "#{",
		ExpressionSuffix: "}",
	}
}

func NewTemplateParserContextWithDelimiters(prefix, suffix string) *ParserContext {
	return &ParserContext{
		IsTemplate:       true,
		ExpressionPrefix: prefix,
		ExpressionSuffix: suffix,
	}
}

// PrintAST 递归打印AST节点结构（参数使用接口类型）
func PrintAST(node SpelNode, level int) {
	if node == nil {
		return
	}

	// 打印缩进
	indent := strings.Repeat("  ", level)

	// 获取节点类型名称
	nodeType := reflect.TypeOf(node).String()
	// 移除包名前缀，只保留类型名
	if idx := strings.LastIndex(nodeType, "."); idx >= 0 {
		nodeType = nodeType[idx+1:]
	}
	// 移除指针符号
	nodeType = strings.TrimPrefix(nodeType, "*")

	// 打印节点信息
	fmt.Printf("%s节点类型: %s, 表达式片段: '%s'\n",
		indent, nodeType, node.ToStringAST())

	// 递归打印子节点
	children := node.GetChildren()
	for _, child := range children {
		// fmt.Printf("%s子节点[%d]:\n", indent, i)
		PrintAST(child, level+1)
	}
}

// PrintASTWithTitle 打印带标题的 AST 树结构
func PrintASTWithTitle(node SpelNode, title string) {
	fmt.Printf("\n=== %s ===\n", title)
	PrintAST(node, 0)
	fmt.Println()
}
