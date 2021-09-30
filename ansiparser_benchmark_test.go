package ansiparser

import (
	"testing"
)

func BenchmarkParseSGR(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseSGR("38;2;0;63;255", "", "")
	}
}

func BenchmarkParseAscii(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Parse("hello world")
	}
}

func BenchmarkParseUnicode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Parse("hello 👍🏼 world")
	}
}
