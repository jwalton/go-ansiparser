// Package ansiparser is a package for parsing strings with ANSI or VT-100 control codes.
//
package ansiparser

const bel = 7
const st = "\u001B\\"

//go:generate stringer -type=TokenType

// TokenType represents the type of a parsed token.
type TokenType int

const (
	// String represents a token containing a plain string.
	String TokenType = 0
	// EscapeCode represents an escape code.
	EscapeCode TokenType = 1
)

// AnsiToken represents a substring parsed from a string containing ANSI escape
// codes.
type AnsiToken struct {
	// Type is the type of this token.
	Type TokenType
	// Content is a slice of the original string which represents the content of
	// this token.
	Content string
	// The foreground color of the text represented by this token, as ANSI codes
	// (e.g. "31" for red, or "38;2;255;20;20" for an RGB color), or an empty
	// string if this is uncolored.  If Type is EscapeCode and this clears the
	// current foreground color, this will be "39".
	FG string
	// The background color of the text represented by this token, as ANSI codes
	// (e.g. "31" for red, or "38;2;255;20;20" for an RGB color), or an empty
	// string if this is uncolored.  If Type is EscapeCode and this clears the
	// current foreground color, this will be "49".
	BG string
}

// PrintLength returns the number of characters this token would take up on-screen.
func (token AnsiToken) PrintLength() int {
	// FIXME: Need to handle variable-width chars
	switch token.Type {
	case String:
		return len(token.Content)
	case EscapeCode:
		return 0
	}

	return 0
}

// Parse parses a string containing ANSI escape codes into a slice of one or more
// AnsiTokens.
func Parse(str string) []AnsiToken {
	tokens := make([]AnsiToken, 0, 1)

	tokenizer := NewStringTokenizer(str)
	for tokenizer.Next() {
		tokens = append(tokens, tokenizer.Token())
	}

	return tokens
}
