package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	pc "github.com/kgolding/go-pocsagencode"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <pager addr> <text message>\ne.g. %s 13100100 'Hello world!'\n", os.Args[0], os.Args[0])
		os.Exit(1)
	}
	addr, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid pager addr - must be a number")
		os.Exit(2)
	}
	messages := []*pc.Message{
		&pc.Message{
			Addr:    uint32(addr),
			Content: os.Args[2],
		},
	}
	burst, _ := pc.Generate(messages)

	binStrs := []string{}
	for _, b := range burst.Bytes() {
		binStrs = append(binStrs, fmt.Sprintf("0b%b", b))
	}
	fmt.Printf("Message: %s\n\n", burst.String())
	fmt.Printf("Hex bytes: % X\n\n", burst.Bytes())
	fmt.Printf("Binary: %s\n\n", strings.Join(binStrs, ", "))
}
