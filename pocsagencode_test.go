package pocsagencode

import (
	"log"
	"os"
	"testing"
)

func init() {
	// SetLogger(log.New(os.Stdout, "POCSAG ", log.LstdFlags))
	// Comment below to enable debug logging
	SetLogger(nil)
}

func Test_Encode_NumericPadding(t *testing.T) {
	SetLogger(log.New(os.Stdout, "POCSAG ", log.LstdFlags))
	defer SetLogger(nil)

	enc, left := Generate([]*Message{
		&Message{1300100, FunctionA, "123", true},
	})
	if len(left) != 0 {
		t.Errorf("expect no message left, got %v", left)
	}

	expect := Burst{
		// 18 words, 576 bits of preamble
		0xAAAAAAAA, 0xAAAAAAAA,
		0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA,
		0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA,
		// The real data starts here
		0x7CD215D8, 0x7A89C197, 0x7A89C197, 0x7A89C197, 0x7A89C197, 0x7A89C197, 0x7A89C197, 0x7A89C197,
		0x7A89C197, 0x4F5A0109, 0xC2619CE1, 0x7A89C197, 0x7A89C197,
	}

	if len(enc) != len(expect) {
		t.Errorf("expected:\n%s\ngot:\n%s\n", expect, enc)
	} else {
		for i, w := range expect {
			if w != enc[i] {
				t.Errorf("expected:%X got:%X at index %d\n", w, enc[i], i)
			}
		}
	}
}

func Test_Encode_Numeric(t *testing.T) {

	enc, left := Generate([]*Message{
		&Message{1300100, FunctionA, "12[3]", true},
	})
	if len(left) != 0 {
		t.Errorf("expect no message left, got %v", left)
	}

	expect := Burst{
		// 18 words, 576 bits of preamble
		0xAAAAAAAA, 0xAAAAAAAA,
		0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA,
		0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA,
		// The real data starts here
		0x7CD215D8, 0x7A89C197, 0x7A89C197, 0x7A89C197, 0x7A89C197, 0x7A89C197, 0x7A89C197, 0x7A89C197,
		0x7A89C197, 0x4F5A0109, 0xC27E3D14, 0x7A89C197, 0x7A89C197,
	}

	if len(enc) != len(expect) {
		t.Errorf("expected:\n%s\ngot:\n%s\n", expect, enc)
	} else {
		for i, w := range expect {
			if w != enc[i] {
				t.Errorf("expected:%X got:%X at index %d\n", w, enc[i], i)
			}
		}
	}
}

func Test_Encode_Alpha(t *testing.T) {
	// SetLogger(log.New(os.Stdout, "POCSAG ", log.LstdFlags))
	enc, left := Generate([]*Message{
		&Message{1300100, FunctionA, "happy christmas!", false},
	})
	if len(left) != 0 {
		t.Errorf("expect no message left, got %v", left)
	}

	expect := Burst{
		// 18 words, 576 bits of preamble
		0xAAAAAAAA, 0xAAAAAAAA,
		0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA,
		0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA, 0xAAAAAAAA,
		// The real dat starts here
		0x7CD215D8, 0x7A89C197, 0x7A89C197, 0x7A89C197, 0x7A89C197, 0x7A89C197, 0x7A89C197, 0x7A89C197,
		0x7A89C197, 0x4F5A0109, 0x8B861C9F, 0xC3CF04CD, 0xD8C5A3C6, 0xF979C8DE, 0xBDB878F3, 0x9E110386,
		0x7A89C197, 0x7CD215D8, 0x7A89C197,
	}

	if len(enc) != len(expect) {
		t.Errorf("expected:\n%s\ngot:\n%s\n", expect, enc)
	} else {
		for i, w := range expect {
			if w != enc[i] {
				t.Errorf("expected:%X got:%X at index %d\n", w, enc[i], i)
			}
		}
	}
}
