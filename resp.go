package resp

import (
	"bufio"
	"fmt"
	"io"
)

const (
	bufferSize         = 32 * 1024
	simpleStringPrefix = '+'
	errorPrefix        = '-'
)

var objectSuffix = []byte("\r\n")

//
// RESP Objects
//

// Object represents a RESP data object reference
type Object interface {
	Dump() []byte
}

// SimpleString represents a RESP Simple String object
type SimpleString string

// NewSimpleString returns a new SimpleString object with content
func NewSimpleString(content string) *SimpleString {
	s := SimpleString(content)
	return &s
}

// Dump returns raw bytes representation
func (s *SimpleString) Dump() []byte {
	raw := []byte{byte(simpleStringPrefix)}
	raw = append(raw, []byte(*s)...)
	raw = append(raw, objectSuffix...)
	return raw
}

// Error represents a RESP Error object
// Note: this is not Go's error.
type Error string

// NewError returns a new Error object with content
func NewError(content string) *Error {
	e := Error(content)
	return &e
}

// Dump returns raw bytes representation
func (e *Error) Dump() []byte {
	raw := []byte{byte(errorPrefix)}
	raw = append(raw, []byte(*e)...)
	raw = append(raw, objectSuffix...)
	return raw
}

//
// Reader
//

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
func (r *Reader) ReadObject() (Object, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}

	return parseLine(line)
}

func (r *Reader) readLine() ([]byte, error) {
	line, isPrefix, err := r.reader.ReadLine()
	if err != nil {
		return nil, err
	}
	if isPrefix {
		return nil, fmt.Errorf("data is too large")
	}

	return line, nil
}

func parseLine(line []byte) (Object, error) {
	if len(line) == 0 {
		return nil, fmt.Errorf("prefix not found")
	}

	switch line[0] {
	case simpleStringPrefix:
		return NewSimpleString(string(line[1:])), nil
	case errorPrefix:
		return NewError(string(line[1:])), nil
	default:
		return nil, fmt.Errorf("unknown prefix %#v", line[0])
	}
}

//
// Writer
//

// Writer implements a RESP object writer for an io.Writer
type Writer struct {
	writer *io.Writer
}

// NewWriter returns a new Writer
func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: &w}
}

// WriteObject writes an object
func (w *Writer) WriteObject(obj Object) error {
	_, err := (*w.writer).Write(obj.Dump())
	return err
}
