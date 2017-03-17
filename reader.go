package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

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

	obj, err := r.parseLine(line)
	if err != nil {
		return nil, err
	}

	switch o := obj.(type) {
	case *BulkString:
		if err := r.readBulkStringData(o); err != nil {
			return nil, err
		}
	}

	return obj, nil
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

func (r *Reader) parseLine(line []byte) (Object, error) {
	if len(line) == 0 {
		return nil, fmt.Errorf("prefix not found")
	}

	switch line[0] {
	case simpleStringPrefix:
		return NewSimpleString(string(line[1:])), nil
	case errorPrefix:
		return NewError(string(line[1:])), nil
	case integerPrefix:
		v, err := toInt(string(line[1:]))
		if err != nil {
			return nil, err
		}
		return NewInteger(v), nil
	case bulkStringPrefix:
		v, err := toInt(string(line[1:]))
		if err != nil {
			return nil, err
		}
		return NewBulkString(v), nil
	case arrayPrefix:
		v, err := toInt(string(line[1:]))
		if err != nil {
			return nil, err
		}
		return NewArray(v), nil
	default:
		return nil, fmt.Errorf("unknown prefix %#v", line[0])
	}
}

func (r *Reader) readBulkStringData(b *BulkString) error {
	// read data content
	n, err := r.reader.Read([]byte(*b))
	if err != nil {
		return err
	}
	if n < len(*b) {
		return fmt.Errorf("too short data")
	}

	// read and discard training CR LF
	// Note: here single line is read and the line must be blank string
	// so that the training data is CR LF
	line, err := r.readLine()
	if err != nil {
		return err
	}
	if len(line) > 0 {
		return fmt.Errorf("too long data")
	}
	return nil
}

func toInt(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("cannot extract integer from %#v: %v", s, err)
	}
	return v, nil
}
