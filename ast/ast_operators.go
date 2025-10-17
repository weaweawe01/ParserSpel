package ast

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

// BinaryOperator represents a binary operator
type BinaryOperator struct {
	*SpelNodeImpl
	Left  SpelNode
	Right SpelNode
}

func NewBinaryOperator(left, right SpelNode, startPos, endPos int) *BinaryOperator {
	children := []SpelNode{left, right}
	return &BinaryOperator{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos, children...),
		Left:         left,
		Right:        right,
	}
}

// OpPlus represents addition operator
type OpPlus struct {
	*BinaryOperator
}

func NewOpPlus(left, right SpelNode, startPos, endPos int) *OpPlus {
	return &OpPlus{
		BinaryOperator: NewBinaryOperator(left, right, startPos, endPos),
	}
}

func (op *OpPlus) GetValue(state *ExpressionState) (interface{}, error) {
	leftVal, err := op.Left.GetValue(state)
	if err != nil {
		return nil, err
	}

	rightVal, err := op.Right.GetValue(state)
	if err != nil {
		return nil, err
	}

	// Handle string concatenation
	if leftStr, ok := leftVal.(string); ok {
		return leftStr + fmt.Sprintf("%v", rightVal), nil
	}
	if rightStr, ok := rightVal.(string); ok {
		return fmt.Sprintf("%v", leftVal) + rightStr, nil
	}

	// Handle numeric addition
	return addNumbers(leftVal, rightVal)
}

func (op *OpPlus) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := op.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (op *OpPlus) ToStringAST() string {
	return fmt.Sprintf("(%s + %s)", op.Left.ToStringAST(), op.Right.ToStringAST())
}

// OpMinus represents subtraction operator
type OpMinus struct {
	*BinaryOperator
}

func NewOpMinus(left, right SpelNode, startPos, endPos int) *OpMinus {
	return &OpMinus{
		BinaryOperator: NewBinaryOperator(left, right, startPos, endPos),
	}
}

func (op *OpMinus) GetValue(state *ExpressionState) (interface{}, error) {
	leftVal, err := op.Left.GetValue(state)
	if err != nil {
		return nil, err
	}

	rightVal, err := op.Right.GetValue(state)
	if err != nil {
		return nil, err
	}

	return subtractNumbers(leftVal, rightVal)
}

func (op *OpMinus) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := op.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (op *OpMinus) ToStringAST() string {
	return fmt.Sprintf("(%s - %s)", op.Left.ToStringAST(), op.Right.ToStringAST())
}

// OpMultiply represents multiplication operator
type OpMultiply struct {
	*BinaryOperator
}

func NewOpMultiply(left, right SpelNode, startPos, endPos int) *OpMultiply {
	return &OpMultiply{
		BinaryOperator: NewBinaryOperator(left, right, startPos, endPos),
	}
}

func (op *OpMultiply) GetValue(state *ExpressionState) (interface{}, error) {
	leftVal, err := op.Left.GetValue(state)
	if err != nil {
		return nil, err
	}

	rightVal, err := op.Right.GetValue(state)
	if err != nil {
		return nil, err
	}

	return multiplyNumbers(leftVal, rightVal)
}

func (op *OpMultiply) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := op.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (op *OpMultiply) ToStringAST() string {
	return fmt.Sprintf("(%s * %s)", op.Left.ToStringAST(), op.Right.ToStringAST())
}

// OpDivide represents division operator
type OpDivide struct {
	*BinaryOperator
}

func NewOpDivide(left, right SpelNode, startPos, endPos int) *OpDivide {
	return &OpDivide{
		BinaryOperator: NewBinaryOperator(left, right, startPos, endPos),
	}
}

func (op *OpDivide) GetValue(state *ExpressionState) (interface{}, error) {
	leftVal, err := op.Left.GetValue(state)
	if err != nil {
		return nil, err
	}

	rightVal, err := op.Right.GetValue(state)
	if err != nil {
		return nil, err
	}

	return divideNumbers(leftVal, rightVal)
}

func (op *OpDivide) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := op.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (op *OpDivide) ToStringAST() string {
	return fmt.Sprintf("(%s / %s)", op.Left.ToStringAST(), op.Right.ToStringAST())
}

// Comparison operators

// OpEQ represents equality operator
type OpEQ struct {
	*BinaryOperator
}

func NewOpEQ(left, right SpelNode, startPos, endPos int) *OpEQ {
	return &OpEQ{
		BinaryOperator: NewBinaryOperator(left, right, startPos, endPos),
	}
}

func (op *OpEQ) GetValue(state *ExpressionState) (interface{}, error) {
	leftVal, err := op.Left.GetValue(state)
	if err != nil {
		return nil, err
	}

	rightVal, err := op.Right.GetValue(state)
	if err != nil {
		return nil, err
	}

	return reflect.DeepEqual(leftVal, rightVal), nil
}

func (op *OpEQ) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := op.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (op *OpEQ) ToStringAST() string {
	return fmt.Sprintf("(%s == %s)", op.Left.ToStringAST(), op.Right.ToStringAST())
}

// OpNE represents not equal operator
type OpNE struct {
	*BinaryOperator
}

func NewOpNE(left, right SpelNode, startPos, endPos int) *OpNE {
	return &OpNE{
		BinaryOperator: NewBinaryOperator(left, right, startPos, endPos),
	}
}

func (op *OpNE) GetValue(state *ExpressionState) (interface{}, error) {
	leftVal, err := op.Left.GetValue(state)
	if err != nil {
		return nil, err
	}

	rightVal, err := op.Right.GetValue(state)
	if err != nil {
		return nil, err
	}

	return !reflect.DeepEqual(leftVal, rightVal), nil
}

func (op *OpNE) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := op.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (op *OpNE) ToStringAST() string {
	return fmt.Sprintf("(%s != %s)", op.Left.ToStringAST(), op.Right.ToStringAST())
}

// OpGT represents greater than operator
type OpGT struct {
	*BinaryOperator
}

func NewOpGT(left, right SpelNode, startPos, endPos int) *OpGT {
	return &OpGT{
		BinaryOperator: NewBinaryOperator(left, right, startPos, endPos),
	}
}

func (op *OpGT) GetValue(state *ExpressionState) (interface{}, error) {
	leftVal, err := op.Left.GetValue(state)
	if err != nil {
		return nil, err
	}

	rightVal, err := op.Right.GetValue(state)
	if err != nil {
		return nil, err
	}

	return compareNumbers(leftVal, rightVal) > 0, nil
}

func (op *OpGT) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := op.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (op *OpGT) ToStringAST() string {
	return fmt.Sprintf("(%s > %s)", op.Left.ToStringAST(), op.Right.ToStringAST())
}

// OpLT represents less than operator
type OpLT struct {
	*BinaryOperator
}

func NewOpLT(left, right SpelNode, startPos, endPos int) *OpLT {
	return &OpLT{
		BinaryOperator: NewBinaryOperator(left, right, startPos, endPos),
	}
}

func (op *OpLT) GetValue(state *ExpressionState) (interface{}, error) {
	leftVal, err := op.Left.GetValue(state)
	if err != nil {
		return nil, err
	}

	rightVal, err := op.Right.GetValue(state)
	if err != nil {
		return nil, err
	}

	return compareNumbers(leftVal, rightVal) < 0, nil
}

func (op *OpLT) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := op.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (op *OpLT) ToStringAST() string {
	return fmt.Sprintf("(%s < %s)", op.Left.ToStringAST(), op.Right.ToStringAST())
}

// Logical operators

// OpAnd represents logical AND operator
type OpAnd struct {
	*BinaryOperator
}

func NewOpAnd(left, right SpelNode, startPos, endPos int) *OpAnd {
	return &OpAnd{
		BinaryOperator: NewBinaryOperator(left, right, startPos, endPos),
	}
}

func (op *OpAnd) GetValue(state *ExpressionState) (interface{}, error) {
	leftVal, err := op.Left.GetValue(state)
	if err != nil {
		return nil, err
	}

	// Short-circuit evaluation
	if !isTruthy(leftVal) {
		return false, nil
	}

	rightVal, err := op.Right.GetValue(state)
	if err != nil {
		return nil, err
	}

	return isTruthy(rightVal), nil
}

func (op *OpAnd) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := op.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (op *OpAnd) ToStringAST() string {
	return fmt.Sprintf("(%s && %s)", op.Left.ToStringAST(), op.Right.ToStringAST())
}

// OpOr represents logical OR operator
type OpOr struct {
	*BinaryOperator
}

func NewOpOr(left, right SpelNode, startPos, endPos int) *OpOr {
	return &OpOr{
		BinaryOperator: NewBinaryOperator(left, right, startPos, endPos),
	}
}

func (op *OpOr) GetValue(state *ExpressionState) (interface{}, error) {
	leftVal, err := op.Left.GetValue(state)
	if err != nil {
		return nil, err
	}

	// Short-circuit evaluation
	if isTruthy(leftVal) {
		return true, nil
	}

	rightVal, err := op.Right.GetValue(state)
	if err != nil {
		return nil, err
	}

	return isTruthy(rightVal), nil
}

func (op *OpOr) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := op.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (op *OpOr) ToStringAST() string {
	return fmt.Sprintf("(%s || %s)", op.Left.ToStringAST(), op.Right.ToStringAST())
}

// Unary operators

// UnaryOperator represents a unary operator
type UnaryOperator struct {
	*SpelNodeImpl
	Child SpelNode
}

func NewUnaryOperator(child SpelNode, startPos, endPos int) *UnaryOperator {
	children := []SpelNode{child}
	return &UnaryOperator{
		SpelNodeImpl: NewSpelNodeImpl(startPos, endPos, children...),
		Child:        child,
	}
}

// OperatorNot represents logical NOT operator
type OperatorNot struct {
	*UnaryOperator
}

func NewOperatorNot(child SpelNode, startPos, endPos int) *OperatorNot {
	return &OperatorNot{
		UnaryOperator: NewUnaryOperator(child, startPos, endPos),
	}
}

func (op *OperatorNot) GetValue(state *ExpressionState) (interface{}, error) {
	childVal, err := op.Child.GetValue(state)
	if err != nil {
		return nil, err
	}

	return !isTruthy(childVal), nil
}

func (op *OperatorNot) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := op.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (op *OperatorNot) ToStringAST() string {
	return fmt.Sprintf("(!%s)", op.Child.ToStringAST())
}

// OperatorMatches represents the matches operator for regex
type OperatorMatches struct {
	*BinaryOperator
}

func NewOperatorMatches(left, right SpelNode, startPos, endPos int) *OperatorMatches {
	return &OperatorMatches{
		BinaryOperator: NewBinaryOperator(left, right, startPos, endPos),
	}
}

func (op *OperatorMatches) GetValue(state *ExpressionState) (interface{}, error) {
	leftVal, err := op.Left.GetValue(state)
	if err != nil {
		return nil, err
	}

	rightVal, err := op.Right.GetValue(state)
	if err != nil {
		return nil, err
	}

	str := fmt.Sprintf("%v", leftVal)
	pattern := fmt.Sprintf("%v", rightVal)

	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %v", err)
	}

	return matched, nil
}

func (op *OperatorMatches) GetTypedValue(state *ExpressionState) (*TypedValue, error) {
	value, err := op.GetValue(state)
	if err != nil {
		return nil, err
	}
	return NewTypedValue(value), nil
}

func (op *OperatorMatches) ToStringAST() string {
	return fmt.Sprintf("(%s matches %s)", op.Left.ToStringAST(), op.Right.ToStringAST())
}

// Helper functions for numeric operations

func addNumbers(left, right interface{}) (interface{}, error) {
	return performNumericOperation(left, right, func(a, b float64) float64 { return a + b })
}

func subtractNumbers(left, right interface{}) (interface{}, error) {
	return performNumericOperation(left, right, func(a, b float64) float64 { return a - b })
}

func multiplyNumbers(left, right interface{}) (interface{}, error) {
	return performNumericOperation(left, right, func(a, b float64) float64 { return a * b })
}

func divideNumbers(left, right interface{}) (interface{}, error) {
	return performNumericOperation(left, right, func(a, b float64) float64 {
		if b == 0 {
			return 0 // Handle division by zero
		}
		return a / b
	})
}

func performNumericOperation(left, right interface{}, op func(float64, float64) float64) (interface{}, error) {
	// For string operands in multiplication, treat as string repetition or concatenation
	if _, isLeftString := left.(string); isLeftString {
		if _, isRightString := right.(string); isRightString {
			// Both strings - not a numeric operation
			return nil, fmt.Errorf("cannot perform numeric operation on two strings")
		}
	}

	leftNum, err := toNumber(left)
	if err != nil {
		return nil, err
	}

	rightNum, err := toNumber(right)
	if err != nil {
		return nil, err
	}

	result := op(leftNum, rightNum)

	// Try to preserve integer type if possible
	if float64(int64(result)) == result {
		return int64(result), nil
	}

	return result, nil
}

func compareNumbers(left, right interface{}) int {
	leftNum, err := toNumber(left)
	if err != nil {
		return 0
	}

	rightNum, err := toNumber(right)
	if err != nil {
		return 0
	}

	if leftNum > rightNum {
		return 1
	} else if leftNum < rightNum {
		return -1
	}
	return 0
}

func toNumber(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to number", value)
	}
}

func isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	case float32:
		return v != 0
	case float64:
		return v != 0
	case string:
		return v != ""
	default:
		return true // Non-nil objects are truthy
	}
}
