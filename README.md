# ansiparser

[![PkgGoDev](https://pkg.go.dev/badge/github.com/jwalton/go-ansiparser)](https://pkg.go.dev/github.com/jwalton/go-ansiparser?readme=expanded#section-readme)
[![Build Status](https://github.com/jwalton/go-ansiparser/workflows/Build/badge.svg)](https://github.com/jwalton/go-ansiparser/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/jwalton/go-ansiparser)](https://goreportcard.com/report/github.com/jwalton/go-ansiparser)
[![Release](https://img.shields.io/github/release/jwalton/go-ansiparser.svg?style=flat-square)](https://github.com/jwalton/go-ansiparser/releases/latest)

ansiparser is a golang library for parsing strings contianing ANSI or VT-100 escape codes.. It will correctly parse 8 bit, 16 bit, and truecolor escape codes out of strings. This library is optimized for the case where a string contains no unicode characters, however it handles unicode characters correctly, and returns them as "CompletChar" tokens to make it easier to work out the printable length of a string.

A quick example:

```go
import (
    "github.com/jwalton/go-ansiparser"
)

func main() {
    result := ansiparser.Parse("hello \u001B[31m👍🏼 \u001B[39mworld")
    // Result will be:
    //
    // []AnsiToken{
    //     {Type: ansiparser.String, Content: "hello ", FG: "", BG: ""},
    //     {Type: ansiparser.EscapeCode, Content: "\u001B[31m", FG: "31", BG: ""},
    //     {Type: ansiparser.ComplexChar, Content: "👍🏼", FG: "31", BG: ""},
    //     {Type: ansiparser.String, Content: " ", FG: "31", BG: ""},
    //     {Type: ansiparser.EscapeCode, Content: "\u001B[39m", FG: "", BG: ""},
    //     {Type: ansiparser.String, Content: "world", FG: "", BG: ""},
    // }
}
```

Calling `ansiparser.Parse()` will return an array of tokens of type:

- `String` for a bare string, with the FG and BG colors set appropriately. The printable length of the string will equal the `len()` of the string.
- `ComplexChar` for any case where more than one character in the string results in only a single printable character in the output.
- `EscapeCode` for any characters that are part of an ANSI escape sequence. These are always 0-width strings when output to a terminal.
- `ZeroWidth` for the `BEL` character.

## Related

- [ansistyles](https://github.com/jwalton/gchalk/tree/master/pkg/ansistyles) - A low level library for generating ANSI escape codes, ported from Node.js's [ansi-styles](https://github.com/chalk/ansi-styles).
- [supportscolor](https://github.com/jwalton/go-supportscolor) - Detect whether a terminal supports color, ported from Node.js's [supports-color](https://github.com/chalk/supports-color).
