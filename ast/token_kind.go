package ast

import "fmt"

// TokenKind represents the different types of tokens in SpEL expressions
type TokenKind int

const (
	// Ordered by priority - operands first
	LITERAL_INT TokenKind = iota
	LITERAL_LONG
	LITERAL_HEXINT
	LITERAL_HEXLONG
	LITERAL_STRING
	LITERAL_REAL
	LITERAL_REAL_FLOAT
	LPAREN // "("
	RPAREN // ")"
	COMMA  // ","
	IDENTIFIER
	COLON            // ":"
	HASH             // "#"
	RSQUARE          // "]"
	LSQUARE          // "["
	LCURLY           // "{"
	RCURLY           // "}"
	DOT              // "."
	PLUS             // "+"
	STAR             // "*"
	MINUS            // "-"
	SELECT_FIRST     // "^["
	SELECT_LAST      // "$["
	QMARK            // "?"
	PROJECT          // "!["
	DIV              // "/"
	GE               // ">="
	GT               // ">"
	LE               // "<="
	LT               // "<"
	EQ               // "=="
	NE               // "!="
	MOD              // "%"
	NOT              // "!"
	ASSIGN           // "="
	INSTANCEOF       // "instanceof"
	MATCHES          // "matches"
	BETWEEN          // "between"
	SELECT           // "?["
	POWER            // "^"
	ELVIS            // "?:"
	SAFE_NAVI        // "?."
	BEAN_REF         // "@"
	FACTORY_BEAN_REF // "&"
	SYMBOLIC_OR      // "||"
	SYMBOLIC_AND     // "&&"
	INC              // "++"
	DEC              // "--"
)

// tokenCharMap maps TokenKind to their string representation
var tokenCharMap = map[TokenKind]string{
	LPAREN:           "(",
	RPAREN:           ")",
	COMMA:            ",",
	COLON:            ":",
	HASH:             "#",
	RSQUARE:          "]",
	LSQUARE:          "[",
	LCURLY:           "{",
	RCURLY:           "}",
	DOT:              ".",
	PLUS:             "+",
	STAR:             "*",
	MINUS:            "-",
	SELECT_FIRST:     "^[",
	SELECT_LAST:      "$[",
	QMARK:            "?",
	PROJECT:          "![",
	DIV:              "/",
	GE:               ">=",
	GT:               ">",
	LE:               "<=",
	LT:               "<",
	EQ:               "==",
	NE:               "!=",
	MOD:              "%",
	NOT:              "!",
	ASSIGN:           "=",
	INSTANCEOF:       "instanceof",
	MATCHES:          "matches",
	BETWEEN:          "between",
	SELECT:           "?[",
	POWER:            "^",
	ELVIS:            "?:",
	SAFE_NAVI:        "?.",
	BEAN_REF:         "@",
	FACTORY_BEAN_REF: "&",
	SYMBOLIC_OR:      "||",
	SYMBOLIC_AND:     "&&",
	INC:              "++",
	DEC:              "--",
}

// tokenNameMap maps TokenKind to their name strings
var tokenNameMap = map[TokenKind]string{
	LITERAL_INT:        "LITERAL_INT",
	LITERAL_LONG:       "LITERAL_LONG",
	LITERAL_HEXINT:     "LITERAL_HEXINT",
	LITERAL_HEXLONG:    "LITERAL_HEXLONG",
	LITERAL_STRING:     "LITERAL_STRING",
	LITERAL_REAL:       "LITERAL_REAL",
	LITERAL_REAL_FLOAT: "LITERAL_REAL_FLOAT",
	LPAREN:             "LPAREN",
	RPAREN:             "RPAREN",
	COMMA:              "COMMA",
	IDENTIFIER:         "IDENTIFIER",
	COLON:              "COLON",
	HASH:               "HASH",
	RSQUARE:            "RSQUARE",
	LSQUARE:            "LSQUARE",
	LCURLY:             "LCURLY",
	RCURLY:             "RCURLY",
	DOT:                "DOT",
	PLUS:               "PLUS",
	STAR:               "STAR",
	MINUS:              "MINUS",
	SELECT_FIRST:       "SELECT_FIRST",
	SELECT_LAST:        "SELECT_LAST",
	QMARK:              "QMARK",
	PROJECT:            "PROJECT",
	DIV:                "DIV",
	GE:                 "GE",
	GT:                 "GT",
	LE:                 "LE",
	LT:                 "LT",
	EQ:                 "EQ",
	NE:                 "NE",
	MOD:                "MOD",
	NOT:                "NOT",
	ASSIGN:             "ASSIGN",
	INSTANCEOF:         "INSTANCEOF",
	MATCHES:            "MATCHES",
	BETWEEN:            "BETWEEN",
	SELECT:             "SELECT",
	POWER:              "POWER",
	ELVIS:              "ELVIS",
	SAFE_NAVI:          "SAFE_NAVI",
	BEAN_REF:           "BEAN_REF",
	FACTORY_BEAN_REF:   "FACTORY_BEAN_REF",
	SYMBOLIC_OR:        "SYMBOLIC_OR",
	SYMBOLIC_AND:       "SYMBOLIC_AND",
	INC:                "INC",
	DEC:                "DEC",
}

// String returns the string representation of TokenKind
func (tk TokenKind) String() string {
	name, exists := tokenNameMap[tk]
	if !exists {
		return "UNKNOWN"
	}

	chars, hasChars := tokenCharMap[tk]
	if hasChars {
		return fmt.Sprintf("%s(%s)", name, chars)
	}
	return name
}

// TokenChars returns the character representation of the token
func (tk TokenKind) TokenChars() string {
	chars, exists := tokenCharMap[tk]
	if exists {
		return chars
	}
	return ""
}

// HasPayload returns true if the token has additional data beyond its kind
func (tk TokenKind) HasPayload() bool {
	chars := tk.TokenChars()
	return len(chars) == 0
}

// GetLength returns the length of the token's character representation
func (tk TokenKind) GetLength() int {
	return len(tk.TokenChars())
}
