package ansiparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsciiString(t *testing.T) {
	result := Parse("hello world")

	assert.Equal(t, []AnsiToken{
		{
			Type:    String,
			Content: "hello world",
			FG:      "",
			BG:      "",
		},
	}, result)
}

func TestUnicodeString(t *testing.T) {
	result := Parse("hello 👍🏼 world")

	assert.Equal(t, []AnsiToken{
		{
			Type:    String,
			Content: "hello ",
			FG:      "",
			BG:      "",
		},
		{
			Type:    ComplexChar,
			Content: "👍🏼",
			FG:      "",
			BG:      "",
		},
		{
			Type:    String,
			Content: " world",
			FG:      "",
			BG:      "",
		},
	}, result)
}

func TestUnicodeStringWithANSI(t *testing.T) {
	result := Parse("hello \u001B[31m👍🏼 \u001B[39mworld")

	assert.Equal(t, []AnsiToken{
		{Type: String, Content: "hello ", FG: "", BG: ""},
		{Type: EscapeCode, Content: "\u001B[31m", FG: "31", BG: ""},
		{Type: ComplexChar, Content: "👍🏼", FG: "31", BG: ""},
		{Type: String, Content: " ", FG: "31", BG: ""},
		{Type: EscapeCode, Content: "\u001B[39m", FG: "", BG: ""},
		{Type: String, Content: "world", FG: "", BG: ""},
	}, result)
}

func TestAsciiStringWithANSI(t *testing.T) {
	result := Parse("hello \u001B[31mred\u001B[39m world")

	assert.Equal(t, []AnsiToken{
		{
			Type:    String,
			Content: "hello ",
			FG:      "",
			BG:      "",
		},
		{
			Type:    EscapeCode,
			Content: "\u001B[31m",
			FG:      "31",
			BG:      "",
		},
		{
			Type:    String,
			Content: "red",
			FG:      "31",
			BG:      "",
		},
		{
			Type:    EscapeCode,
			Content: "\u001B[39m",
			FG:      "",
			BG:      "",
		},
		{
			Type:    String,
			Content: " world",
			FG:      "",
			BG:      "",
		},
	}, result)
}

func TestAsciiStringWithOSC(t *testing.T) {
	result := Parse("hello \u001B]8;;http://thedreaming.org\u001B\\link\u001B]8;;\u001B\\")

	assert.Equal(t, []AnsiToken{
		{
			Type:    String,
			Content: "hello ",
			FG:      "",
			BG:      "",
		},
		{
			Type:    EscapeCode,
			Content: "\u001B]8;;http://thedreaming.org\u001B\\",
			FG:      "",
			BG:      "",
		},
		{
			Type:    String,
			Content: "link",
			FG:      "",
			BG:      "",
		},
		{
			Type:    EscapeCode,
			Content: "\u001B]8;;\u001B\\",
			FG:      "",
			BG:      "",
		},
	}, result)
}

func TestAsciiStringWithCursorMovement(t *testing.T) {
	result := Parse("hello \u001B[31m\u001B[1Cworld\u001B[39m")

	assert.Equal(t, []AnsiToken{
		{
			Type:    String,
			Content: "hello ",
			FG:      "",
			BG:      "",
		},
		{
			Type:    EscapeCode,
			Content: "\u001B[31m",
			FG:      "31",
			BG:      "",
		},
		{
			Type:    EscapeCode,
			Content: "\u001B[1C",
			FG:      "31",
			BG:      "",
		},
		{
			Type:    String,
			Content: "world",
			FG:      "31",
			BG:      "",
		},
		{
			Type:    EscapeCode,
			Content: "\u001B[39m",
			FG:      "",
			BG:      "",
		},
	}, result)
}

func TestAsciiReset(t *testing.T) {
	result := Parse("\u001B[31;42mhello\u001B[1m world")

	assert.Equal(t, []AnsiToken{
		{
			Type:    EscapeCode,
			Content: "\u001B[31;42m",
			FG:      "31",
			BG:      "42",
		},
		{
			Type:    String,
			Content: "hello",
			FG:      "31",
			BG:      "42",
		},
		{
			Type:    EscapeCode,
			Content: "\u001B[1m",
			FG:      "",
			BG:      "",
		},
		{
			Type:    String,
			Content: " world",
			FG:      "",
			BG:      "",
		},
	}, result)
}

func TestAsciiRGB(t *testing.T) {
	result := Parse("\u001B[38;2;0;30;255;48;2;255;90;0mhello")

	assert.Equal(t, []AnsiToken{
		{
			Type:    EscapeCode,
			Content: "\u001B[38;2;0;30;255;48;2;255;90;0m",
			FG:      "38;2;0;30;255",
			BG:      "48;2;255;90;0",
		},
		{
			Type:    String,
			Content: "hello",
			FG:      "38;2;0;30;255",
			BG:      "48;2;255;90;0",
		},
	}, result)
}
