package pocsagencode

import (
	"log"
)

// Message is a single POCSAG Alphanumeric message
type Message struct {
	Addr    uint32
	Content string
}

var logger *log.Logger

// SetLogger can be passed a *log.Logger to enable log output
// Example pocsagencoder.SetLogger(log.New(os.Stdout, "POCSAG ", log.LstdFlags))
func SetLogger(Logger *log.Logger) {
	logger = Logger
}

func debugf(format string, args ...interface{}) {
	if logger != nil {
		logger.Printf(format, args...)
	}
}

// The POCSAG transmission starts with 576 bit reversals (101010...).
// That's 576/8 == 72 bytes of 0xAA.
var pocsagPreambleWord uint32 = 0xAAAAAAAA

// The Frame Synchronisation (FS) code is 32 bits:
// 01111100 11010010 00010101 11011000
var pocsagFrameSyncWord uint32 = 0x7CD215D8

// The Idle Codeword:
// 01111010 10001001 11000001 10010111
var pocsagIdleWord uint32 = 0x7A89C197

// calcBchAndParity calculates the binary checksum and parity for a codeword
func calcBchAndParity(cw uint32) uint32 {

	// make sure the 11 LSB are 0.
	cw &= 0xFFFFF800

	parity := 0

	// calculate bch
	localCw := cw
	for bit := 1; bit <= 21; bit++ {
		if cw&0x80000000 > 0 {
			cw ^= 0xED200000
		}
		cw = cw << 1
	}
	localCw |= (cw >> 21)
	// at this point $local_cw has codeword with bch

	// calculate parity
	cw = localCw
	for bit := 1; bit <= 32; bit++ {
		if cw&0x80000000 > 0 {
			parity++
		}
		cw = cw << 1
	}

	// turn last bit to 1 depending on parity
	cw_with_parity := localCw
	if parity%2 != 0 {
		cw_with_parity = localCw + 1
	}

	debugf("  bch_and_parity returning %X\n", cw_with_parity)
	return cw_with_parity
}

//
//	Given the numeric destination address and function, generate an address codeword.
//

// sub _address_codeword($$)
func addressCodeword(inAddr uint32, function byte) uint32 {
	// POCSAG recommendation 1.3.2
	// The three least significant bits are not transmitted but
	// serve to define the frame in which the address codeword
	// must be transmitted.
	// So we take them away.
	// shift address to right by two bits to remove the least significant bits
	addr := inAddr >> 3

	// truncate address to 18 bits
	addr &= 0x3FFFF

	// truncate function to 2 bits
	function &= 0x3

	// codeword without parity
	codeword := addr<<13 | uint32(function)<<11

	debugf("  generated address codeword for %d function %d: %X\n", inAddr, function, codeword)

	return calcBchAndParity(codeword)
}

// appendMessageCodeword appends a message content codeword to the message, calculating bch+parity for it
func appendMessageCodeword(word uint32) uint32 {
	return calcBchAndParity(word | 1<<31)
}

// reverseBits reverses the bits in a byte. Used to encode characters in a text message,
//since the opposite order is used when transmitting POCSAG text.
func reverseBits(in byte) byte {
	out := byte(0)

	for i := byte(0); i < 7; i++ {
		out |= ((in >> i) & 0x01) << (6 - i)
	}

	return out
}

// appendContentText appends text message content to the transmission blob
func appendContentText(content string) (int, Burst) {
	out := make(Burst, 0)
	debugf("appendContentText: %s", content)

	bitpos := 0
	word := uint32(0)
	leftbits := 0
	pos := 0

	// walk through characters in message
	for i, r := range content {
		// make sure it's 7 bits
		char := byte(r & 0x7f)

		debugf("  char %d: %d [%X]\n", i, char, char)

		char = reverseBits(char)

		//  if the bits won't fit:
		if bitpos+7 > 20 {
			space := 20 - bitpos
			//  leftbits least significant bits of $char are left over in the next word
			leftbits = 7 - space
			debugf("  bits of char won't fit since bitpos is %d, got %d bits free, leaving %d bits in next word", bitpos, space, leftbits)
		}

		word |= (uint32(char) << uint(31-7-bitpos))

		bitpos += 7

		if bitpos >= 20 {
			debugf("   appending word: %X\n", word)
			out = append(out, appendMessageCodeword(word))
			pos++
			word = 0
			bitpos = 0
		}

		if leftbits > 0 {
			word |= (uint32(char) << uint(31-leftbits))
			bitpos = leftbits
			leftbits = 0
		}
	}

	if bitpos > 0 {
		debugf("  got %d bits in word at end of text, word: %X", bitpos, word)
		step := 0
		for bitpos < 20 {
			if step == 2 {
				word |= (1 << uint(30-bitpos))
			}
			bitpos++
			step++
			if step == 7 {
				step = 0
			}
		}
		out = append(out, appendMessageCodeword(word))
		pos++
	}

	return pos, out
}

// appendMessage appends a single message to the end of the transmission blob.
func appendMessage(startpos int, msg *Message) (int, Burst) {
	// expand the parameters of the message
	addr := msg.Addr
	function := byte(0)
	type_ := 'a'
	content := msg.Content

	debugf("append_message: addr %d function %d type %d content %s", addr, function, type_, content)

	// the starting frame is selected based on the three least significant bits
	frameAddr := addr & 7
	frameAddrCw := frameAddr * 2 // or << 2 ?

	debugf("  frame_addr is %d, current position %d", frameAddr, startpos)

	// append idle codewords, until we're in the right frame for this address
	tx := make(Burst, 0)
	pos := 0
	for uint32(startpos+pos)%16 != frameAddrCw {
		debugf("   inserting IDLE codewords in position %d (%d)", startpos+pos, (startpos+pos)%16)
		tx = append(tx, pocsagIdleWord)
		pos++
	}

	// Then, append the address codeword, containing the function and the address
	// (sans 3 least significant bits, which are indicated by the starting frame,
	// which the receiver is waiting for)
	tx = append(tx, addressCodeword(addr, function))
	pos++

	// Next, append the message contents
	contentEncLen, contentEnc := appendContentText(content)

	tx = append(tx, contentEnc...)
	pos += contentEncLen

	// Return the current frame position and the binary string to be appended
	return pos, tx
}

// insertSCS inserts Synchronisation Codewords before every 8 POCSAG frames
// (frame is SC+ 64 bytes of address and message codewords)
func insertSCS(tx Burst) Burst {
	out := make(Burst, 0)

	// each batch is SC + 8 frames, each frame is 2 codewords,
	// each codeword is 32 bits, so we must insert an SC
	// every (8*2*32) bits == 64 bytes
	txLen := len(tx)
	for i := 0; i < txLen; i += 16 {
		//  put in the CW and 64 the next 64 bytes
		out = append(out, pocsagFrameSyncWord)
		end := i + 16
		if end > txLen {
			end = txLen
		}
		out = append(out, tx[i:end]...)
	}

	return out
}

// selectMsg selects the optimal next message to be appended, trying to
// minimize the amount of idle codewords transmitted
func selectMsg(pos int, msgListRef []*Message) int {
	currentPick := -1
	currentDist := 0
	posFrame := uint32(pos/2) % 8

	debugf("select_msg pos %d: %d", pos, posFrame)

	for i := 0; i < len(msgListRef); i++ {
		addr := msgListRef[i].Addr
		frameAddr := addr & 7
		distance := int(frameAddr - posFrame)
		if distance < 0 {
			distance += 8
		}

		debugf("  considering list item %d: %d - frame addr %d distance %d\n", i, addr, frameAddr, distance)

		if frameAddr == posFrame {
			debugf("  exact match %d: %d - frame addr %d\n", i, addr, frameAddr)
			return i
		}

		if currentPick == -1 {
			debugf("  first option %d: %d - frame addr %d distance %d\n", i, addr, frameAddr, distance)
			currentPick = i
			currentDist = distance
			continue
		}

		if distance < currentDist {
			debugf("  better option %d: %d - frame addr %d distance %d", i, addr, frameAddr, distance)
			currentPick = i
			currentDist = distance
		}
	}

	return currentPick
}

// Generate a transmission from an array of given messages, to fit with the maximum lenght
// The function returns the an array of Uint32 to be keyed over the air in FSK, and
// any messages which did not fit in the transmission, given the maximum
// transmission length (in bytes) given in the first parameter. They can be passed
// in the next Generate() call and sent in the next brrraaaap.
func Generate(messages []*Message, optionFns ...OptionFn) (Burst, []*Message) {
	options := &Options{
		MaxLen:       3000,
		PreambleBits: 576,
	}
	for _, opt := range optionFns {
		opt(options)
	}

	txWithoutScs := make(Burst, 0)
	debugf("generate_transmission, maxlen: %d", options.MaxLen)

	pos := 0
	for len(messages) > 0 {
		//  figure out an optimal next message to minimize the amount of required idle codewords
		//  TODO: do a deeper search, considering the length of the message and a possible
		//  optimal next recipient
		optimalNextMsg := selectMsg(pos, messages)
		msg := messages[optimalNextMsg]
		messages = append(messages[:optimalNextMsg], messages[optimalNextMsg+1:]...)

		appendLen, x := appendMessage(pos, msg)

		nextLen := pos + appendLen + 2
		//  initial sync codeword + one for every 16 codewords
		nextLen += 1 + int((nextLen-1)/16)
		nextLenBytes := nextLen * 4
		debugf("after this message of %d codewords, burst will be %d codewords and %d bytes long\n", appendLen, nextLen, nextLenBytes)

		if nextLenBytes > int(options.MaxLen) {
			if pos == 0 {
				debugf("burst would become too large (%d > %d) with first message alone - discarding!", nextLenBytes, options.MaxLen)
			} else {
				debugf("burst would become too large (%d > %d) - returning msg in queue", nextLenBytes, options.MaxLen)
				messages = append([]*Message{msg}, messages...)
				break
			}
		} else {
			txWithoutScs = append(txWithoutScs, x...)
			pos += appendLen
		}
	}

	// if the burst is empty, return it as completely empty
	if pos == 0 {
		return Burst{}, messages
	}

	// append a couple of IDLE codewords, otherwise many pagers will
	// happily decode the junk in the end and show it to the recipient
	txWithoutScs = append(txWithoutScs, pocsagIdleWord)
	txWithoutScs = append(txWithoutScs, pocsagIdleWord)

	burstLen := len(txWithoutScs)
	debugf("transmission without SCs: %d bytes, %d codewords\n%X\n", burstLen*4, burstLen, txWithoutScs)

	// put SC every 8 frames
	burst := insertSCS(txWithoutScs)

	burstLen = len(burst)
	debugf("transmission with SCs: %d bytes, %d codewords\n%X\n", burstLen*4, burstLen, burst)

	if options.PreambleBits > 0 {
		preambleWords := options.PreambleBits / 32
		if options.PreambleBits%32 > 0 {
			preambleWords++
		}
		preamble := make(Burst, preambleWords)
		for i := range preamble {
			preamble[i] = pocsagPreambleWord
		}
		burst = append(preamble, burst...)
	}

	return burst, messages
}
