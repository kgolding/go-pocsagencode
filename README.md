Golang port of the `POCSAG::Encode` Perl module extended to support Numeric as well as AlphaNumeric messages.

Example usage

```
package main

import (
	"log"

	"github.com/kgolding/go-pocsagencode"
)

func main() {
	messages := []*pocsagencode.Message{
		&pocsagencode.Message{1300100, FunctionA, "Hello Pager!", false},
	}

	log.Println("Sending", len(messages), "messages")
	var burst pocsagencode.Burst
	for len(messages) > 0 {
		burst, messages = pocsagencode.Generate(messages)
		// Options can be set as below for MaxLen and PreambleBits
		// burst, messages = pocsagencode.Generate(messages, pocsagencode.OptionPreambleBits(250))
		log.Println("Burst", burst.String())
		// Send Burst to the FSK modem here...
	}
	log.Println("Done")
}
```