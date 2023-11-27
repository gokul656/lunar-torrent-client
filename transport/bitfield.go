package main

type Bitfield []byte

/*
* Check whether the Bitfield has the piece
*
* Assume we've 16 Bits, it can be stored in 2 Bytes
*
* piece := 6
* bf := BitField{1111 1000, 0}
*
* First, we've to calculate ByteIndex
* byteIndex := 6 / 8  => returns 0
*
* Next, we've to check inside the bit array:
*
* bf[0] = 1111 1000
*
* Now, let's verify if the 6th piece is available in the bitfield:
*
* 1. Right Shift Operation:
*    - Calculate the offset within the byte: offset := 6 % 8 => returns 6
*    - Calculate the shift count: shiftCount := 7 - offset => 7 - 6 = 1
*    - Perform a right shift by 1 position: bf[0] >> 1 => 0111 1100
*
* 2. Bitwise AND with 1:
*    - The rightmost bit after the shift is 0.
*    - 0111 1100 & 0000 0001 => 0000 0000
*
* 3. Comparison with 0:
*    - The result is 0, indicating that the 6th piece is not available in the bitfield.
*
* Therefore, the expression bf[byteIndex] >> uint(7-offset) & 1 != 0 evaluates to false.
 */
func (bf Bitfield) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8
	if index < 0 || byteIndex >= len(bf) {
		return false
	}

	return bf[byteIndex]>>(7-offset)&1 != 0
}
