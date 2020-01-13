package pocsagencode

import (
	"encoding/binary"
	"fmt"
)

type Burst []uint32

// String return a formated multiline string with the
func (b Burst) String() string {
	s := ""
	preambleCount := 0
	preambleOver := false
	for _, w := range b {
		if !preambleOver {
			if w == pocsagPreambleWord {
				preambleCount++
				continue
			} else {
				s += fmt.Sprintf("[%d bits of 1010101010... (0xAA) preamble] ", preambleCount*32)
				preambleOver = true
			}
		}
		if w == pocsagFrameSyncWord {
			s += "[Batch start/sync] "
		}
		s += fmt.Sprintf("%X ", w)
	}
	return s
}

// Bytes returns a []byte
func (b Burst) Bytes() []byte {
	buf := make([]byte, len(b)*4)
	for i, w := range b {
		binary.BigEndian.PutUint32(buf[i*4:], w)
	}
	return buf
}
