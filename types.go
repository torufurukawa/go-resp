package resp

import "fmt"

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

// Integer represents a RESP Integer object
type Integer int

// NewInteger returns a new Integer object with v
func NewInteger(v int) *Integer {
	i := Integer(v)
	return &i
}

// Dump returns raw bytes representation
func (i *Integer) Dump() []byte {
	return []byte(fmt.Sprintf(":%d\r\n", int(*i)))
}

// BulkString represents a RESP Bulk String object
type BulkString []byte

// NewBulkString returns a new BulkString object whose length is size
func NewBulkString(size int) *BulkString {
	b := make(BulkString, size)
	return &b
}

// Dump returns raw bytes representation
func (b *BulkString) Dump() []byte {
	data := []byte(*b)
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(data), string(data)))
}

// Array represents a RESP Array object
type Array []Object

// NewArray returns a new Array object whose length is size
func NewArray(size int) *Array {
	a := Array(make([]Object, size))
	return &a
}

// Dump returns raw bytes representation
func (a *Array) Dump() []byte {
	return nil
}
