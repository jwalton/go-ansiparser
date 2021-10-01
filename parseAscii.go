package ansiparser

// This is the "fast" parser which is used when all characters in `str` are part
// of the ASCII character set.  In this case, we know every char in the string
// either takes up one character on screen, or is part of an ANSI escape code.
//
// When `parseASCII()` reaches the end of the string, or encounters a non-ASCII
// character, it will return any tokens it managed to parse, and the length
// of the string that it managed to consume.
func parseASCII(str string, force bool, startFG string, startBG string) (tokens []AnsiToken, consumed int) {
	currentFG := startFG
	currentBG := startBG

	tokens = make([]AnsiToken, 0, 1)
	currentStart := -1

	finishStringToken := func(end int) {
		if currentStart == -1 {
			return
		}

		tokens = append(tokens, AnsiToken{
			Type:    String,
			Content: str[currentStart:end],
			FG:      currentFG,
			BG:      currentBG,
		})

		currentStart = -1
	}

	i := 0
	for i < len(str) {
		c := str[i]
		if c > 127 && !force {
			// Non-ascii character. Abort.
			// FIXME: Rewind one byte because the previous char may be
			// modified by this char.
			finishStringToken(i)
			break
		} else if c == '\u001B' && (i+1) < len(str) && str[i+1] == '[' {
			finishStringToken(i)

			// Control Sequence Introducer (CSI)
			escapeCode := parseASCIIEscapeCode(str[i:], currentFG, currentBG)
			currentFG = escapeCode.FG
			currentBG = escapeCode.BG
			tokens = append(tokens, escapeCode)
			i += len(escapeCode.Content)
		} else if c == '\u001B' && (i+1) < len(str) && str[i+1] == ']' {
			finishStringToken(i)

			// Operating System Command (OSC)
			escapeCode := parseASCIIOSC(str[i:], currentFG, currentBG)
			tokens = append(tokens, escapeCode)
			i += len(escapeCode.Content)
		} else if c == bel {
			finishStringToken(i)
			tokens = append(tokens, AnsiToken{
				Type:    ZeroWidth,
				Content: str[i : i+1],
				FG:      currentFG,
				BG:      currentBG,
			})
			i++
		} else {
			if currentStart == -1 {
				// Start reading a string token
				currentStart = i
			} else {
				// Continue reading a string token...
			}
			i++
		}
	}

	finishStringToken(i)

	return tokens, i
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
			// Eat the ;
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
			// Set backgrond to dim 4-bit color.
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
