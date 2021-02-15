package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	pc "github.com/kgolding/go-pocsagencode"
)

func main() {
	var num bool
	var debug bool

	flag.BoolVar(&num, "num", false, "send as numeric message")
	flag.BoolVar(&debug, "v", false, "verbose logging")
	flag.Parse()

	if len(flag.Args()) != 2 {
		me := filepath.Base(os.Args[0])
		fmt.Printf(`
Usage: %s [-num] <pager addr> <text or numeric message>
  Examples:  %s 13100100 'Hello world!'
             %s -num 13100100 9876\n`, me, me, me)
		os.Exit(1)
	}
	addr, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		fmt.Println("Invalid pager addr - must be a number")
		os.Exit(2)
	}
	messages := []*pc.Message{
		&pc.Message{
			Addr:      uint32(addr),
			Content:   flag.Arg(0),
			IsNumeric: num,
		},
	}

	if debug {
		pc.SetLogger(log.New(os.Stdout, "", log.Lshortfile))
	}

	burst, _ := pc.Generate(messages)

	fmt.Printf("Message: %s\n\n", burst.String())
	fmt.Printf("Hex bytes: % X\n\n", burst.Bytes())
	fmt.Printf("Binary: %s\n\n", burst.BinStr())
}
