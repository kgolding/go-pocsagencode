package main

import (
	"fmt"
	"log"

	"github.com/kgolding/go-pocsagencode"
)

func main() {
	messages := []*pocsagencode.Message{
		&pocsagencode.Message{1300100, "Hello Pager!", false},
	}

	for i := 0; i < 50; i++ {
		addr := uint32(1200000 + i*100)
		messages = append(messages, &pocsagencode.Message{addr, fmt.Sprintf("Hello pager number %d", addr), false})
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
