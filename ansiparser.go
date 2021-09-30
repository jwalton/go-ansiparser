// Package ansiparser is a package for parsing strings with ANSI or VT-100 control codes.
//
// This library is optimized for the case where a string contains no unicode
// characters, however it handles unicode characters correctly, and returns them
// as "CompletChar" tokens to make it easier to work out the printable length
// of a string.
//
package ansiparser

const bel = 7
const cr = 13
const cursorRight = 28
const cursonDown = 31
const st = "\u001B\\"

//go:generate stringer -type=TokenType

// TokenType represents the type of a parsed token.
type TokenType int

const (
	// String represents a token containing a plain string, where each char
	// is a single printable character.
	String TokenType = 0
	// EscapeCode represents an escape code.
	EscapeCode TokenType = 1
	// ComplexChar represents a single printable character which takes up more
	// than one char in the input.  Examples are codepoints that take up more than
	// two bytes, or multiple emoji joined together with zero-width-joiners.
	ComplexChar TokenType = 2
	// ZeroWidth represents a zero-width character (such as \n or BEL).
	ZeroWidth TokenType = 3
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

// Parse parses a string containing ANSI escape codes into a slice of one or more
// AnsiTokens.
func Parse(str string) []AnsiToken {
	tokens, consumed := parseASCII(str, false, "", "")
	if consumed < len(str) {
		// Didn't consume all tokens - string must contain non-ASCII characters.
		startFG := ""
		startBG := ""
		if len(tokens) > 0 {
			lastToken := tokens[len(tokens)-1]
			startFG = lastToken.FG
			startBG = lastToken.BG
		}
		unicodeTokens, _ := parseUnicode(str[consumed:], startFG, startBG)
		tokens = append(tokens, unicodeTokens...)
	}

	return tokens
}
