package security

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	var tests = []struct {
		s, want string
	}{
		{"<\nscript", "<<br />script"},
		{"<script>", "&lt;script>"},
		{"<style>", "&lt;style>"},
		{"<style<script>>", "&lt;style&lt;script>>"},
	}

	for _, c := range tests {
		got := Escape(c.s)
		fmt.Printf("Escape(%q) == %q, want %q\n", c.s, got, c.want)
		if got != c.want {
			t.Errorf("Escape(%q) == %q, want %q", c.s, got, c.want)
		}
	}
	fmt.Printf("\n")
}
