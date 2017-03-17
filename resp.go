package resp

import (
	"bufio"
	"fmt"
	"io"
)

const (
	bufferSize         = 32 * 1024
	simpleStringPrefix = '+'
)

// Object represents a RESP data object
type Object interface{}

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
func (r *Reader) ReadObject() (string, error) {
	line, err := r.readLine()
	if err != nil {
		return "", err
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

func parseLine(line []byte) (string, error) {
	if len(line) == 0 {
		return "", fmt.Errorf("prefix not found")
	}

	switch line[0] {
	case simpleStringPrefix:
		return string(line[1:]), nil
	default:
		return "", fmt.Errorf("unknown prefix %#v", line[0])
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
func (w *Writer) WriteObject(obj string) error {
	rawData := []byte{byte(simpleStringPrefix)}
	rawData = append(rawData, []byte(obj)...)
	rawData = append(rawData, []byte("\r\n")...)
	_, err := (*w.writer).Write(rawData)
	return err
}
