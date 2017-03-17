package resp

import (
	"bufio"
	"fmt"
	"io"
)

const (
	bufferSize int = 32 * 1024
)

// Object represents a RESP data object
type Object interface{}

// Reader implements a RESP object reader for an io.Reader
type Reader struct {
	reader *bufio.Reader
}

// NewReader returns a new Reader
func NewReader(r io.Reader) *Reader {
	reader := Reader{reader: bufio.NewReaderSize(r, bufferSize)}
	return &reader
}

// ReadObject reads next RESP object
func (r *Reader) ReadObject() (string, error) {
	line, isPrefix, err := r.reader.ReadLine()
	if err != nil {
		return "", err
	}
	if isPrefix {
		// TODO define error
		return "", fmt.Errorf("data is too large")
	}

	// TODO parse explicitly
	return string(line[1:]), nil
}
