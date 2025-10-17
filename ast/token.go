package ast

import "fmt"

// Token represents a holder for a kind of token, the associated data,
// and its position in the input data stream (start/end).
type Token struct {
	Kind     TokenKind
	Data     *string // Nullable string pointer to match Java's @Nullable String
	StartPos int
	EndPos   int
}

// NewToken creates a new token without data (for tokens like TRUE or '+')
func NewToken(tokenKind TokenKind, startPos, endPos int) *Token {
	return &Token{
		Kind:     tokenKind,
		Data:     nil,
		StartPos: startPos,
		EndPos:   endPos,
	}
}

// NewTokenWithData creates a new token with data
func NewTokenWithData(tokenKind TokenKind, tokenData []rune, startPos, endPos int) *Token {
	var data *string
	if tokenData != nil {
		str := string(tokenData)
		data = &str
	}

	return &Token{
		Kind:     tokenKind,
		Data:     data,
		StartPos: startPos,
		EndPos:   endPos,
	}
}

// GetKind returns the token kind
func (t *Token) GetKind() TokenKind {
	return t.Kind
}

// IsIdentifier checks if the token is an identifier
func (t *Token) IsIdentifier() bool {
	return t.Kind == IDENTIFIER
}

// IsNumericRelationalOperator checks if the token is a numeric relational operator
func (t *Token) IsNumericRelationalOperator() bool {
	return t.Kind == GT || t.Kind == GE || t.Kind == LT ||
		t.Kind == LE || t.Kind == EQ || t.Kind == NE
}

// StringValue returns the string value of the token's data
func (t *Token) StringValue() string {
	if t.Data != nil {
		return *t.Data
	}
	return ""
}

// AsInstanceOfToken returns a new INSTANCEOF token with the same position
func (t *Token) AsInstanceOfToken() *Token {
	return NewToken(INSTANCEOF, t.StartPos, t.EndPos)
}

// AsMatchesToken returns a new MATCHES token with the same position
func (t *Token) AsMatchesToken() *Token {
	return NewToken(MATCHES, t.StartPos, t.EndPos)
}

// AsBetweenToken returns a new BETWEEN token with the same position
func (t *Token) AsBetweenToken() *Token {
	return NewToken(BETWEEN, t.StartPos, t.EndPos)
}

// String returns the string representation of the token
func (t *Token) String() string {
	result := fmt.Sprintf("[%s", t.Kind.String())

	if t.Kind.HasPayload() && t.Data != nil {
		result += fmt.Sprintf(":%s", *t.Data)
	}

	result += fmt.Sprintf("](%d,%d)", t.StartPos, t.EndPos)
	return result
}
