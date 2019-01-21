package cmd

import (
	"path"
	"strings"
)

type sharedObject struct {
	filename string
	content  []byte
}

func (s *sharedObject) decode(data []byte) {
	// 1. decode the filename [x | name | y | ext]
	nameLength := int(data[0])
	filename := string(data[1 : nameLength+1])

	extLength := int(data[1+nameLength])

	ext := string(data[1+nameLength+1 : 1+nameLength+1+extLength])
	parts := []string{filename, ext}
	fname := strings.Join(parts, ".")

	s.filename = fname
	s.content = data[nameLength+extLength+3:]

}

// [x | name | y | ext]
// CUBA.JPG
// ../CUBA.JPG
// img/CUBA.JPG

// encode a filename into bytes
func encodeFileName(filePath string) []byte {
	_, file := path.Split(filePath)
	ext := strings.Trim(path.Ext(file), ".")
	split := strings.Split(file, ".")
	filename := strings.Join(split[:len(split)-1], "")

	encoded := make([]byte, 0)
	encoded = append(encoded, byte(len(filename)))
	encoded = append(encoded, []byte(filename)...)
	encoded = append(encoded, byte(len(ext)))
	encoded = append(encoded, []byte(ext)...)
	encoded = append(encoded, byte(0))

	return encoded
}
