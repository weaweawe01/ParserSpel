package ast

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// Tokenizer lexes input data into a stream of tokens that can then be parsed
type Tokenizer struct {
	expressionString string
	charsToProcess   []rune
	pos              int
	max              int
	tokens           []*Token
}

// Alternative textual operator names which must match enum constant names in TokenKind
// Note that AND and OR are also alternative textual names, but they are handled later
// This list must remain sorted since we use binary search
var alternativeOperatorNames = []string{
	"BETWEEN", "DIV", "EQ", "GE", "GT", "LE", "LT", "MOD", "NE", "NOT",
}

// Flags for character classification
const (
	IS_DIGIT    = 0x01
	IS_HEXDIGIT = 0x02
)

var flags [256]byte

func init() {
	// Initialize character classification flags
	for ch := '0'; ch <= '9'; ch++ {
		flags[ch] |= IS_DIGIT | IS_HEXDIGIT
	}
	for ch := 'A'; ch <= 'F'; ch++ {
		flags[ch] |= IS_HEXDIGIT
	}
	for ch := 'a'; ch <= 'f'; ch++ {
		flags[ch] |= IS_HEXDIGIT
	}

	// Sort alternative operator names for binary search
	sort.Strings(alternativeOperatorNames)
}

// NewTokenizer creates a new tokenizer for the given input data
func NewTokenizer(inputData string) *Tokenizer {
	charsToProcess := []rune(inputData + "\x00") // Add null terminator
	return &Tokenizer{
		expressionString: inputData,
		charsToProcess:   charsToProcess,
		pos:              0,
		max:              len(charsToProcess),
		tokens:           make([]*Token, 0),
	}
}

// Process tokenizes the input and returns a list of tokens
func (t *Tokenizer) Process() ([]*Token, error) {
	for t.pos < t.max {
		ch := t.charsToProcess[t.pos]

		if t.isAlphabetic(ch) {
			if err := t.lexIdentifier(); err != nil {
				return nil, err
			}
		} else {
			switch ch {
			case '+':
				if t.isTwoCharToken(INC) {
					t.pushPairToken(INC)
				} else {
					t.pushCharToken(PLUS)
				}
			case '_':
				if err := t.lexIdentifier(); err != nil { // '_' is another way to start an identifier
					return nil, err
				}
			case '-':
				if t.isTwoCharToken(DEC) {
					t.pushPairToken(DEC)
				} else {
					t.pushCharToken(MINUS)
				}
			case ':':
				t.pushCharToken(COLON)
			case '.':
				t.pushCharToken(DOT)
			case ',':
				t.pushCharToken(COMMA)
			case '*':
				t.pushCharToken(STAR)
			case '/':
				t.pushCharToken(DIV)
			case '%':
				t.pushCharToken(MOD)
			case '(':
				t.pushCharToken(LPAREN)
			case ')':
				t.pushCharToken(RPAREN)
			case '[':
				t.pushCharToken(LSQUARE)
			case '#':
				t.pushCharToken(HASH)
			case ']':
				t.pushCharToken(RSQUARE)
			case '{':
				t.pushCharToken(LCURLY)
			case '}':
				t.pushCharToken(RCURLY)
			case '@':
				t.pushCharToken(BEAN_REF)
			case '^':
				if t.isTwoCharToken(SELECT_FIRST) {
					t.pushPairToken(SELECT_FIRST)
				} else {
					t.pushCharToken(POWER)
				}
			case '!':
				if t.isTwoCharToken(NE) {
					t.pushPairToken(NE)
				} else if t.isTwoCharToken(PROJECT) {
					t.pushPairToken(PROJECT)
				} else {
					t.pushCharToken(NOT)
				}
			case '=':
				if t.isTwoCharToken(EQ) {
					t.pushPairToken(EQ)
				} else {
					t.pushCharToken(ASSIGN)
				}
			case '&':
				if t.isTwoCharToken(SYMBOLIC_AND) {
					t.pushPairToken(SYMBOLIC_AND)
				} else {
					t.pushCharToken(FACTORY_BEAN_REF)
				}
			case '|':
				if !t.isTwoCharToken(SYMBOLIC_OR) {
					return nil, t.raiseParseException(t.pos, "MISSING_CHARACTER", "|")
				}
				t.pushPairToken(SYMBOLIC_OR)
			case '?':
				if t.isTwoCharToken(SELECT) {
					t.pushPairToken(SELECT)
				} else if t.isTwoCharToken(ELVIS) {
					t.pushPairToken(ELVIS)
				} else if t.isTwoCharToken(SAFE_NAVI) {
					t.pushPairToken(SAFE_NAVI)
				} else {
					t.pushCharToken(QMARK)
				}
			case '$':
				if t.isTwoCharToken(SELECT_LAST) {
					t.pushPairToken(SELECT_LAST)
				} else {
					if err := t.lexIdentifier(); err != nil { // '$' is another way to start an identifier
						return nil, err
					}
				}
			case '>':
				if t.isTwoCharToken(GE) {
					t.pushPairToken(GE)
				} else {
					t.pushCharToken(GT)
				}
			case '<':
				if t.isTwoCharToken(LE) {
					t.pushPairToken(LE)
				} else {
					t.pushCharToken(LT)
				}
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				if err := t.lexNumericLiteral(ch == '0'); err != nil {
					return nil, err
				}
			case ' ', '\t', '\r', '\n':
				// Skip whitespace
				t.pos++
			case '\'':
				if err := t.lexQuotedStringLiteral(); err != nil {
					return nil, err
				}
			case '"':
				if err := t.lexDoubleQuotedStringLiteral(); err != nil {
					return nil, err
				}
			case 0:
				// Hit sentinel at end of value
				t.pos++ // will take us to the end
			case '\\':
				return nil, t.raiseParseException(t.pos, "UNEXPECTED_ESCAPE_CHAR")
			default:
				return nil, t.raiseParseException(t.pos+1, "UNSUPPORTED_CHARACTER", fmt.Sprintf("%c(%d)", ch, int(ch)))
			}
		}
	}
	return t.tokens, nil
}

// lexQuotedStringLiteral processes a single-quoted string literal
func (t *Tokenizer) lexQuotedStringLiteral() error {
	start := t.pos
	terminated := false

	for !terminated {
		t.pos++
		if t.isExhausted() {
			return t.raiseParseException(start, "NON_TERMINATING_QUOTED_STRING")
		}

		ch := t.charsToProcess[t.pos]
		if ch == '\'' {
			// May not be the end if the char after is also a '
			if t.pos+1 < len(t.charsToProcess) && t.charsToProcess[t.pos+1] == '\'' {
				t.pos++ // skip over that too, and continue
			} else {
				terminated = true
			}
		}
	}

	t.pos++
	data := t.subarray(start, t.pos)
	t.tokens = append(t.tokens, NewTokenWithData(LITERAL_STRING, data, start, t.pos))
	return nil
}

// lexDoubleQuotedStringLiteral processes a double-quoted string literal
func (t *Tokenizer) lexDoubleQuotedStringLiteral() error {
	start := t.pos
	terminated := false

	for !terminated {
		t.pos++
		if t.isExhausted() {
			return t.raiseParseException(start, "NON_TERMINATING_DOUBLE_QUOTED_STRING")
		}

		ch := t.charsToProcess[t.pos]
		if ch == '"' {
			// May not be the end if the char after is also a "
			if t.pos+1 < len(t.charsToProcess) && t.charsToProcess[t.pos+1] == '"' {
				t.pos++ // skip over that too, and continue
			} else {
				terminated = true
			}
		}
	}

	t.pos++
	data := t.subarray(start, t.pos)
	t.tokens = append(t.tokens, NewTokenWithData(LITERAL_STRING, data, start, t.pos))
	return nil
}

// lexNumericLiteral processes numeric literals (int, long, real, hex)
func (t *Tokenizer) lexNumericLiteral(firstCharIsZero bool) error {
	isReal := false
	start := t.pos

	if t.pos+1 < len(t.charsToProcess) {
		ch := t.charsToProcess[t.pos+1]
		isHex := ch == 'x' || ch == 'X'

		// Deal with hexadecimal
		if firstCharIsZero && isHex {
			t.pos += 2 // Skip '0' and 'x'/'X'
			hexStart := t.pos
			for t.pos < len(t.charsToProcess) && t.isHexadecimalDigit(t.charsToProcess[t.pos]) {
				t.pos++
			}

			if t.pos < len(t.charsToProcess) && t.isChar('L', 'l') {
				data := t.subarray(hexStart, t.pos)
				t.pushHexIntToken(data, true, start, t.pos)
				t.pos++
			} else {
				data := t.subarray(hexStart, t.pos)
				t.pushHexIntToken(data, false, start, t.pos)
			}
			return nil
		}
	}

	// Consume first part of number
	for t.pos < len(t.charsToProcess) && t.isDigit(t.charsToProcess[t.pos]) {
		t.pos++
	}

	// A '.' indicates this number is a real
	if t.pos < len(t.charsToProcess) {
		ch := t.charsToProcess[t.pos]
		if ch == '.' {
			isReal = true
			dotpos := t.pos
			t.pos++ // move past the dot
			// Continue consuming digits
			for t.pos < len(t.charsToProcess) && t.isDigit(t.charsToProcess[t.pos]) {
				t.pos++
			}
			if t.pos == dotpos+1 {
				// The number is something like '3.'. It is really an int but may be
				// part of something like '3.toString()'. Process it as an int and leave the dot.
				t.pos = dotpos
				data := t.subarray(start, t.pos)
				t.pushIntToken(data, false, start, t.pos)
				return nil
			}
		}
	}

	endOfNumber := t.pos

	// Check for long suffix
	if t.pos < len(t.charsToProcess) && t.isChar('L', 'l') {
		if isReal {
			return t.raiseParseException(start, "REAL_CANNOT_BE_LONG")
		}
		data := t.subarray(start, endOfNumber)
		t.pushIntToken(data, true, start, endOfNumber)
		t.pos++
		return nil
	}

	// Check for exponent
	if t.pos < len(t.charsToProcess) && t.isExponentChar(t.charsToProcess[t.pos]) {
		isReal = true
		t.pos++
		if t.pos < len(t.charsToProcess) {
			possibleSign := t.charsToProcess[t.pos]
			if t.isSign(possibleSign) {
				t.pos++
			}
		}

		// Exponent digits
		for t.pos < len(t.charsToProcess) && t.isDigit(t.charsToProcess[t.pos]) {
			t.pos++
		}

		isFloat := false
		if t.pos < len(t.charsToProcess) && t.isFloatSuffix(t.charsToProcess[t.pos]) {
			isFloat = true
			endOfNumber = t.pos + 1
			t.pos++
		} else if t.pos < len(t.charsToProcess) && t.isDoubleSuffix(t.charsToProcess[t.pos]) {
			endOfNumber = t.pos + 1
			t.pos++
		}
		data := t.subarray(start, t.pos)
		t.pushRealToken(data, isFloat, start, t.pos)
		return nil
	}

	// Check for float/double suffix
	if t.pos < len(t.charsToProcess) {
		ch := t.charsToProcess[t.pos]
		isFloat := false
		if t.isFloatSuffix(ch) {
			isReal = true
			isFloat = true
			endOfNumber = t.pos + 1
			t.pos++
		} else if t.isDoubleSuffix(ch) {
			isReal = true
			endOfNumber = t.pos + 1
			t.pos++
		}

		if isReal {
			data := t.subarray(start, endOfNumber)
			t.pushRealToken(data, isFloat, start, endOfNumber)
		} else {
			data := t.subarray(start, endOfNumber)
			t.pushIntToken(data, false, start, endOfNumber)
		}
	} else {
		data := t.subarray(start, endOfNumber)
		t.pushIntToken(data, false, start, endOfNumber)
	}

	return nil
}

// lexIdentifier processes identifiers and alternative operator names
func (t *Tokenizer) lexIdentifier() error {
	start := t.pos
	for t.pos < len(t.charsToProcess) && t.isIdentifier(t.charsToProcess[t.pos]) {
		t.pos++
	}

	subarray := t.subarray(start, t.pos)

	// Check if this is the alternative (textual) representation of an operator
	if len(subarray) >= 2 && len(subarray) <= 7 { // Support 2-7 character operators
		asString := strings.ToUpper(string(subarray))
		idx := sort.SearchStrings(alternativeOperatorNames, asString)
		if idx < len(alternativeOperatorNames) && alternativeOperatorNames[idx] == asString {
			tokenKind := t.getTokenKindByName(asString)
			t.pushOneCharOrTwoCharToken(tokenKind, start, subarray)
			return nil
		}
	}

	t.tokens = append(t.tokens, NewTokenWithData(IDENTIFIER, subarray, start, t.pos))
	return nil
}

// Helper functions for token creation
func (t *Tokenizer) pushIntToken(data []rune, isLong bool, start, end int) {
	if isLong {
		t.tokens = append(t.tokens, NewTokenWithData(LITERAL_LONG, data, start, end))
	} else {
		t.tokens = append(t.tokens, NewTokenWithData(LITERAL_INT, data, start, end))
	}
}

func (t *Tokenizer) pushHexIntToken(data []rune, isLong bool, start, end int) error {
	if len(data) == 0 {
		if isLong {
			return t.raiseParseException(start, "NOT_A_LONG", t.expressionString[start:end+1])
		} else {
			return t.raiseParseException(start, "NOT_AN_INTEGER", t.expressionString[start:end])
		}
	}

	if isLong {
		t.tokens = append(t.tokens, NewTokenWithData(LITERAL_HEXLONG, data, start, end))
	} else {
		t.tokens = append(t.tokens, NewTokenWithData(LITERAL_HEXINT, data, start, end))
	}
	return nil
}

func (t *Tokenizer) pushRealToken(data []rune, isFloat bool, start, end int) {
	if isFloat {
		t.tokens = append(t.tokens, NewTokenWithData(LITERAL_REAL_FLOAT, data, start, end))
	} else {
		t.tokens = append(t.tokens, NewTokenWithData(LITERAL_REAL, data, start, end))
	}
}

func (t *Tokenizer) subarray(start, end int) []rune {
	if end > len(t.charsToProcess) {
		end = len(t.charsToProcess)
	}
	if start > end {
		start = end
	}
	return t.charsToProcess[start:end]
}

func (t *Tokenizer) isTwoCharToken(kind TokenKind) bool {
	tokenChars := kind.TokenChars()
	return len(tokenChars) == 2 &&
		t.pos+1 < len(t.charsToProcess) &&
		t.charsToProcess[t.pos] == rune(tokenChars[0]) &&
		t.charsToProcess[t.pos+1] == rune(tokenChars[1])
}

func (t *Tokenizer) pushCharToken(kind TokenKind) {
	t.tokens = append(t.tokens, NewToken(kind, t.pos, t.pos+1))
	t.pos++
}

func (t *Tokenizer) pushPairToken(kind TokenKind) {
	t.tokens = append(t.tokens, NewToken(kind, t.pos, t.pos+2))
	t.pos += 2
}

func (t *Tokenizer) pushOneCharOrTwoCharToken(kind TokenKind, pos int, data []rune) {
	t.tokens = append(t.tokens, NewTokenWithData(kind, data, pos, pos+kind.GetLength()))
}

// Character classification methods
func (t *Tokenizer) isIdentifier(ch rune) bool {
	return t.isAlphabetic(ch) || t.isDigit(ch) || ch == '_' || ch == '$'
}

func (t *Tokenizer) isChar(a, b rune) bool {
	if t.pos >= len(t.charsToProcess) {
		return false
	}
	ch := t.charsToProcess[t.pos]
	return ch == a || ch == b
}

func (t *Tokenizer) isExponentChar(ch rune) bool {
	return ch == 'e' || ch == 'E'
}

func (t *Tokenizer) isFloatSuffix(ch rune) bool {
	return ch == 'f' || ch == 'F'
}

func (t *Tokenizer) isDoubleSuffix(ch rune) bool {
	return ch == 'd' || ch == 'D'
}

func (t *Tokenizer) isSign(ch rune) bool {
	return ch == '+' || ch == '-'
}

func (t *Tokenizer) isDigit(ch rune) bool {
	if ch > 255 {
		return false
	}
	return (flags[ch] & IS_DIGIT) != 0
}

func (t *Tokenizer) isAlphabetic(ch rune) bool {
	return unicode.IsLetter(ch)
}

func (t *Tokenizer) isHexadecimalDigit(ch rune) bool {
	if ch > 255 {
		return false
	}
	return (flags[ch] & IS_HEXDIGIT) != 0
}

func (t *Tokenizer) isExhausted() bool {
	return t.pos == t.max-1
}

func (t *Tokenizer) raiseParseException(start int, msg string, inserts ...interface{}) error {
	return fmt.Errorf("parse exception at position %d: %s %v", start, msg, inserts)
}

// getTokenKindByName returns the TokenKind for a given operator name
func (t *Tokenizer) getTokenKindByName(name string) TokenKind {
	switch name {
	case "BETWEEN":
		return BETWEEN
	case "DIV":
		return DIV
	case "EQ":
		return EQ
	case "GE":
		return GE
	case "GT":
		return GT
	case "LE":
		return LE
	case "LT":
		return LT
	case "MOD":
		return MOD
	case "NE":
		return NE
	case "NOT":
		return NOT
	default:
		return IDENTIFIER // fallback
	}
}
