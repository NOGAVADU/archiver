package vlc

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type BinaryChunk string

type BinaryChunks []BinaryChunk

const chunkSize = 8

func NewBinChunks(data []byte) BinaryChunks {
	res := make(BinaryChunks, 0, len(data))

	for _, part := range data {
		res = append(res, NewBinChunk(part))
	}

	return res
}

func NewBinChunk(code byte) BinaryChunk {
	return BinaryChunk(fmt.Sprintf("%08b", code))
}

func (bcs BinaryChunks) Bytes() []byte {
	res := make([]byte, 0, len(bcs))

	for _, bc := range bcs {
		res = append(res, bc.Byte())
	}

	return res
}

func (bc BinaryChunk) Byte() byte {
	num, err := strconv.ParseUint(string(bc), 2, chunkSize)
	if err != nil {
		panic("can't parse binary chunk: " + err.Error())
	}

	return byte(num)
}

// Join joins chunks into one line and returns as string
func (bcs BinaryChunks) Join() string {
	var buf strings.Builder

	for _, bc := range bcs {
		buf.WriteString(string(bc))
	}

	return buf.String()
}

// splitByChunks splits binary string by chunks with given size,
// i.g. '1010101010110101' -> '1010010 01011011 01100110'
func splitByChunks(bStr string, chunkSize int) BinaryChunks {
	strLen := utf8.RuneCountInString(bStr)
	chunksCount := strLen / chunkSize

	if strLen/chunkSize != 0 {
		chunksCount++
	}

	res := make(BinaryChunks, 0, chunksCount)

	var buf strings.Builder

	for i, ch := range bStr {
		buf.WriteString(string(ch))

		if (i+1)%chunkSize == 0 {
			res = append(res, BinaryChunk(buf.String()))
			buf.Reset()
		}
	}

	if buf.Len() != 0 {
		lastChunk := buf.String()

		lastChunk += strings.Repeat("0", chunkSize-len(lastChunk))

		res = append(res, BinaryChunk(lastChunk))
	}

	return res
}
