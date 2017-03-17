package resp

import (
	"bufio"
	"io"
)

// Object represents a RESP data object
type Object interface{}

// Reader implements a RESP object reader for an io.Reader
type Reader struct {
	reader *bufio.Reader
}

// NewReader returns a new Reader
func NewReader(r io.Reader) *Reader {
	reader := Reader{reader: bufio.NewReader(r)}
	return &reader
}

// ReadObject reads next RESP object
func (r *Reader) ReadObject() (string, error) {
	line, isPrefix, err := r.reader.ReadLine()
	if err != nil {
		return "", err
	}
	// TODO handle case of isPrefix == true
	_ = isPrefix

	// TODO parse
	return string(line[1:]), nil
}
