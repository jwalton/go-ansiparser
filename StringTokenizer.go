package ansiparser

// StringTokenizer tokenizes a string.
type StringTokenizer struct {
	token    AnsiToken
	input    string
	position int
}

// NewStringTokenizer returns a new instance of StringTokenizer, which is used
// to tokenizer the input string.  Call `Next()` to see if there is a next token,
// and if this returns true the current token can be read from `Token()`.
func NewStringTokenizer(input string) *StringTokenizer {
	return &StringTokenizer{
		input:    input,
		position: 0,
	}
}

// Token returns the current token.
func (tokenizer *StringTokenizer) Token() AnsiToken {
	return tokenizer.token
}

// Next parses the next token from the input string.  Returns true if a token
// was found, false if the end of the input was reached.  The token is available
// via `Token()`.
func (tokenizer *StringTokenizer) Next() bool {
	str := tokenizer.input

	// The start of the token we are currently reading.
	currentStart := tokenizer.position

	// This function works by reading one character at a time.  Most of the time
	// we'll be reading a string, so we keep incremention `tokenizer.position`.
	// Whenever we run into a character that signals the start of an escape code,
	// we call this function which will generate a token if we have anything in
	// the string.  Otherwise we go ahead and handle the escape code.
	makeStringToken := func() bool {
		if currentStart == tokenizer.position || tokenizer.position > len(str) {
			return false
		}

		tokenizer.token = AnsiToken{
			Type:    String,
			Content: str[currentStart:tokenizer.position],
			FG:      tokenizer.token.FG,
			BG:      tokenizer.token.BG,
		}
		return true
	}

	for tokenizer.position < len(str) {
		c := str[tokenizer.position]
		if c > 127 {
			// Skip over any multi-byte UTF-8 characters.
			// This works because the first bit of any multi-byte UTF-8 character
			// is always 1 - the first byte starts with a 1, and all continuation
			// characters start with 10.
			tokenizer.position++

		} else if c == '\u001B' && (tokenizer.position+1) < len(str) && str[tokenizer.position+1] == '[' {
			// Control Sequence Introducer (CSI)
			if makeStringToken() {
				return true
			}

			escapeCode := parseASCIIEscapeCode(str[tokenizer.position:], tokenizer.token.FG, tokenizer.token.BG)
			tokenizer.token = escapeCode
			tokenizer.position += len(escapeCode.Content)
			return true
		} else if c == '\u001B' && (tokenizer.position+1) < len(str) && str[tokenizer.position+1] == ']' {
			// Operating System Command (OSC)
			if makeStringToken() {
				return true
			}

			escapeCode := parseASCIIOSC(str[tokenizer.position:], tokenizer.token.FG, tokenizer.token.BG)
			tokenizer.token = escapeCode
			tokenizer.position += len(escapeCode.Content)
			return true
		} else {
			// Add this character to the string we are reading...
			tokenizer.position++
		}
	}

	// If there's anything left over, return a string token.
	return makeStringToken()
}

func parseASCIIOSC(
	str string,
	prevFG string,
	prevBG string,
) AnsiToken {
	// Skip OSC
	i := 2

	done := func() bool {
		return i >= len(str) ||
			str[i] == bel ||
			str[i-1:i+1] == st
	}

	for !done() {
		i++
	}

	if i < len(str) {
		// Consume the termination character.
		i++
	}

	return AnsiToken{
		Type:    EscapeCode,
		Content: str[0:i],
		FG:      prevFG,
		BG:      prevBG,
	}
}

// parseASCIIEscapeCode parses an escape code from a string.
// Returns `end` which is the index of the first character after the escape code,
// and `token` which is the parsed token.
//
// `start` should be the index of the CSI ("\u001B[").
//
// This will return a token with FG = closeFgTag if this escape code clears
// the foreground color (even if it does it via a reset) and similarly with
// BG = closeBgTag if it clears if background color.  FG and BG will be set
// to the empty string if this token neither sets nor clears the color.
//
func parseASCIIEscapeCode(
	str string,
	prevFG string,
	prevBG string,
) (token AnsiToken) {
	token = AnsiToken{
		Type: EscapeCode,
	}

	var command byte

	// Skip the CSI
	var i = 2

	// Read parameter bytes
	for i < len(str) && str[i] >= 0x30 && str[i] <= 0x3F {
		i++
	}

	// Read intermediate bytes
	for i < len(str) && str[i] >= 0x20 && str[i] <= 0x2F {
		i++
	}

	// Read the final byte
	if i < len(str) && str[i] >= 40 && str[i] <= 0x7E {
		command = str[i]
		i++
	}

	token.Content = str[0:i]
	if command == 'm' {
		token.FG, token.BG = parseSGR(str[2:i-1], prevFG, prevBG)
	} else {
		token.FG = prevFG
		token.BG = prevBG
	}

	return token
}

// parseSGR parses an "select graphics rendition" string (e.g. "38;2;0;63;255" to
// set the forground color to rgb(0, 63, 255) or "1;93" to reset the foreground
// and background colors and then set the forground to bright yellow).
func parseSGR(
	sgr string,
	prevFG string,
	prevBG string,
) (fg string, bg string) {
	if len(sgr) == 0 {
		// Empty SGR is same as reset
		return "", ""
	}

	pos := 0
	readNextCommand := func() string {
		start := pos
		for pos < len(sgr) && sgr[pos] != ';' {
			pos++
		}
		answer := sgr[start:pos]
		if pos < len(sgr) {
			// Consume the ;
			pos++
		}
		return answer
	}

	fg = prevFG
	bg = prevBG

	parseSetColor := func(startPos int, command string) string {
		t := readNextCommand()
		if t == "5" {
			// Set ANSI 256 color
			c := readNextCommand()
			return sgr[startPos : startPos+len(command)+1+len(t)+1+len(c)]
		} else if t == "2" {
			// Set RGB color
			r := readNextCommand()
			g := readNextCommand()
			b := readNextCommand()
			return sgr[startPos : startPos+len(command)+1+
				len(t)+1+
				len(r)+1+
				len(g)+1+
				len(b)]
		} else {
			// ???
			return ""
		}
	}

	for pos < len(sgr) {
		startPos := pos
		command := readNextCommand()

		if command == "1" {
			// Reset
			fg = ""
			bg = ""
		} else if command == "39" {
			// Reset foreground
			fg = ""
		} else if command == "38" {
			// Set foreground color
			fg = parseSetColor(startPos, command)
		} else if len(command) == 2 && command[0] == '3' {
			// Set foreground to dim 4-bit color.
			fg = command
		} else if command == "90" ||
			command == "91" ||
			command == "92" ||
			command == "93" ||
			command == "94" ||
			command == "95" ||
			command == "96" ||
			command == "97" {
			// Set foreground to bright 4-bit color.
			fg = command
		} else if command == "49" {
			// Reset background
			bg = ""
		} else if command == "48" {
			// Set background
			bg = parseSetColor(startPos, command)
		} else if len(command) == 2 && command[0] == '4' {
			// Set background to dim 4-bit color.
			bg = command
		} else if command == "100" ||
			command == "101" ||
			command == "102" ||
			command == "103" ||
			command == "104" ||
			command == "105" ||
			command == "106" ||
			command == "107" {
			// Set backround to bright 4-bit color.
			bg = command
		} else {
			// ??? - Unknown command, skip.
		}
	}

	return fg, bg
}
