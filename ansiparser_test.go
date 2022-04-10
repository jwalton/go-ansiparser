package ansiparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseString(t *testing.T) {
	result := Parse("hello world")

	assert.Equal(t, []AnsiToken{
		{
			Type:    String,
			Content: "hello world",
			FG:      "",
			BG:      "",
			IsASCII: true,
		},
	}, result)
}

func TestUnicodeString(t *testing.T) {
	result := Parse("hello 👍🏼 world")

	assert.Equal(t, []AnsiToken{
		{
			Type:    String,
			Content: "hello 👍🏼 world",
			FG:      "",
			BG:      "",
			IsASCII: false,
		},
	}, result)
}

func TestUnicodeStringWithANSI(t *testing.T) {
	result := Parse("hello \u001B[31m👍🏼 \u001B[39mworld")

	assert.Equal(t, []AnsiToken{
		{Type: String, Content: "hello ", FG: "", BG: "", IsASCII: true},
		{Type: EscapeCode, Content: "\u001B[31m", FG: "31", BG: "", IsASCII: true},
		{Type: String, Content: "👍🏼 ", FG: "31", BG: "", IsASCII: false},
		{Type: EscapeCode, Content: "\u001B[39m", FG: "", BG: "", IsASCII: true},
		{Type: String, Content: "world", FG: "", BG: "", IsASCII: true},
	}, result)
}

func TestStringWithANSI(t *testing.T) {
	result := Parse("hello \u001B[31mred\u001B[39m world")

	assert.Equal(t, []AnsiToken{
		{
			Type:    String,
			Content: "hello ",
			FG:      "",
			BG:      "",
			IsASCII: true,
		},
		{
			Type:    EscapeCode,
			Content: "\u001B[31m",
			FG:      "31",
			BG:      "",
			IsASCII: true,
		},
		{
			Type:    String,
			Content: "red",
			FG:      "31",
			BG:      "",
			IsASCII: true,
		},
		{
			Type:    EscapeCode,
			Content: "\u001B[39m",
			FG:      "",
			BG:      "",
			IsASCII: true,
		},
		{
			Type:    String,
			Content: " world",
			FG:      "",
			BG:      "",
			IsASCII: true,
		},
	}, result)
}

func TestStringWithOSC(t *testing.T) {
	result := Parse("hello \u001B]8;;http://thedreaming.org\u001B\\link\u001B]8;;\u001B\\")

	assert.Equal(t, []AnsiToken{
		{
			Type:    String,
			Content: "hello ",
			FG:      "",
			BG:      "",
			IsASCII: true,
		},
		{
			Type:    EscapeCode,
			Content: "\u001B]8;;http://thedreaming.org\u001B\\",
			FG:      "",
			BG:      "",
			IsASCII: true,
		},
		{
			Type:    String,
			Content: "link",
			FG:      "",
			BG:      "",
			IsASCII: true,
		},
		{
			Type:    EscapeCode,
			Content: "\u001B]8;;\u001B\\",
			FG:      "",
			BG:      "",
			IsASCII: true,
		},
	}, result)
}

func TestStringWithCursorMovement(t *testing.T) {
	result := Parse("hello \u001B[31m\u001B[1Cworld\u001B[39m")

	assert.Equal(t, []AnsiToken{
		{
			Type:    String,
			Content: "hello ",
			FG:      "",
			BG:      "",
			IsASCII: true,
		},
		{
			Type:    EscapeCode,
			Content: "\u001B[31m",
			FG:      "31",
			BG:      "",
			IsASCII: true,
		},
		{
			Type:    EscapeCode,
			Content: "\u001B[1C",
			FG:      "31",
			BG:      "",
			IsASCII: true,
		},
		{
			Type:    String,
			Content: "world",
			FG:      "31",
			BG:      "",
			IsASCII: true,
		},
		{
			Type:    EscapeCode,
			Content: "\u001B[39m",
			FG:      "",
			BG:      "",
			IsASCII: true,
		},
	}, result)
}

func TestReset(t *testing.T) {
	result := Parse("\u001B[31;42mhello\u001B[1m world")

	assert.Equal(t, []AnsiToken{
		{
			Type:    EscapeCode,
			Content: "\u001B[31;42m",
			FG:      "31",
			BG:      "42",
			IsASCII: true,
		},
		{
			Type:    String,
			Content: "hello",
			FG:      "31",
			BG:      "42",
			IsASCII: true,
		},
		{
			Type:    EscapeCode,
			Content: "\u001B[1m",
			FG:      "",
			BG:      "",
			IsASCII: true,
		},
		{
			Type:    String,
			Content: " world",
			FG:      "",
			BG:      "",
			IsASCII: true,
		},
	}, result)
}

func TestRGB(t *testing.T) {
	result := Parse("\u001B[38;2;0;30;255;48;2;255;90;0mhello")

	assert.Equal(t, []AnsiToken{
		{
			Type:    EscapeCode,
			Content: "\u001B[38;2;0;30;255;48;2;255;90;0m",
			FG:      "38;2;0;30;255",
			BG:      "48;2;255;90;0",
			IsASCII: true,
		},
		{
			Type:    String,
			Content: "hello",
			FG:      "38;2;0;30;255",
			BG:      "48;2;255;90;0",
			IsASCII: true,
		},
	}, result)
}

func TestTokenizer(t *testing.T) {
	tokenizer := NewStringTokenizer("hello \u001B[31m👍🏼 \u001B[39mworld")

	assert.Equal(t, true, tokenizer.Next())
	assert.Equal(t,
		AnsiToken{Type: String, Content: "hello ", FG: "", BG: "", IsASCII: true},
		tokenizer.Token(),
	)

	assert.Equal(t, true, tokenizer.Next())
	assert.Equal(t,
		AnsiToken{Type: EscapeCode, Content: "\u001B[31m", FG: "31", BG: "", IsASCII: true},
		tokenizer.Token(),
	)

	assert.Equal(t, true, tokenizer.Next())
	assert.Equal(t,
		AnsiToken{Type: String, Content: "👍🏼 ", FG: "31", BG: "", IsASCII: false},
		tokenizer.Token(),
	)

	assert.Equal(t, true, tokenizer.Next())
	assert.Equal(t,
		AnsiToken{Type: EscapeCode, Content: "\u001B[39m", FG: "", BG: "", IsASCII: true},
		tokenizer.Token(),
	)

	assert.Equal(t, true, tokenizer.Next())
	assert.Equal(t,
		AnsiToken{Type: String, Content: "world", FG: "", BG: "", IsASCII: true},
		tokenizer.Token(),
	)

	assert.Equal(t, false, tokenizer.Next())
}
