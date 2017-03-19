package resp

import "io"

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
