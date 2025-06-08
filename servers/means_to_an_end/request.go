package means_to_an_end

import "encoding/binary"

type Request []byte

func (r *Request) Decode() (operation rune, input1 int32, input2 int32) {
	operation = rune((*r)[0])
	input1 = int32(binary.BigEndian.Uint32((*r)[1:5]))
	input2 = int32(binary.BigEndian.Uint32((*r)[5:9]))
	return operation, input1, input2
}
