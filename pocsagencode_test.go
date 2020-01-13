package pocsagencode

import (
	"testing"
)

func Test_Encode(t *testing.T) {
	// SetLogger(log.New(os.Stdout, "POCSAG ", log.LstdFlags))
	enc, left := Generate([]*Message{
		&Message{1300100, "happy christmas!"},
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
	t.Log(enc)
	t.Log(enc.Bytes())
}
