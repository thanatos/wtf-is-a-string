package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	// Strings are just byte strings, not text strings:
	var s string = "\xed\xa0\xbd\x00"
	fmt.Printf("Hello, 世界 %s\n", s)
	const nihongo = "日本語"
	showString(nihongo)
	showString(s)
	
	// A rune is just an int32
	// This isn't a scalar value:
	var r rune = 0xd83d
	fmt.Printf("Rune? %v\n", r)
	// This isn't even a code point…
	var r2 rune = 0xfffffff
	fmt.Printf("Rune? %v\n", r2)
}

func showString(s string) {
	fmt.Printf("For string \"%s\" (Valid? %v):\n", s, utf8.ValidString(s))
	// range will happily iterate through non-UTF-8 strings, emitting replacement chars if it can't decode
	for index, runeValue := range s {
		fmt.Printf("  %#U starts at byte position %d\n", runeValue, index)
	}
}
