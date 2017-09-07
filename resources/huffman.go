package resources

import "fmt"

//https://github.com/scummvm/scummvm/blob/master/engines/sci/decompressor.cpp#L108
//http://sciwiki.sierrahelp.com//index.php?title=SCI_Specifications:_Chapter_2_-_Resource_files#Decompression_algorithm_HUFFMAN
func huffmanDecode(dat []byte, expectedLength uint16) []byte {
	numNodes := int(dat[0])
	terminator := dat[1]
	nodes := dat[2 : 2+numNodes*2]
	dat = dat[2+numNodes*2:]
	curByte := 0
	curBit := byte(7)
	out := make([]byte, 0, expectedLength)
	n := 0
	getBit := func() byte {
		n++
		if curByte >= len(dat) {
			fmt.Println(len(out), expectedLength, curByte, n)
		}
		v := (dat[curByte] >> curBit) & 0x01
		curBit--
		if curBit == 0xff {
			curBit = 7
			curByte++
		}
		return v
	}
	getByte := func() byte {
		var v byte
		for i := 0; i < 8; i++ {
			v = v<<1 | getBit()
		}
		return v
	}
	var getNextChar func(int) (byte, bool)
	getNextChar = func(idx int) (byte, bool) {
		val := nodes[idx*2]
		sibs := nodes[1+idx*2]
		// leaf node
		if sibs == 0 {
			return val, false
		}
		bit := getBit()
		if bit == 1 {
			//right branch
			right := int(sibs & 0x0f)
			if right == 0 {
				//literal
				return getByte(), true
			}
			return getNextChar(idx + int(right))
		}
		//left
		left := int(sibs&0xf0) >> 4
		return getNextChar(idx + int(left))
	}

	for {
		c, ok := getNextChar(0)
		if ok && c == terminator {
			break
		}
		out = append(out, c)
	}
	return out
}
