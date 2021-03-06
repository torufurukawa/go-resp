package resp

import "fmt"

// Object represents a RESP data object reference
type Object interface {
	Dump() []byte
}

// SimpleString represents a RESP Simple String object
type SimpleString struct {
	content string
}

// NewSimpleString returns a new SimpleString object with content
func NewSimpleString(content string) *SimpleString {
	s := SimpleString{content: content}
	return &s
}

// String returns the content string
func (s *SimpleString) String() string {
	return s.content
}

// Dump returns raw bytes representation
func (s *SimpleString) Dump() []byte {
	raw := []byte{byte(simpleStringPrefix)}
	raw = append(raw, []byte(s.String())...)
	raw = append(raw, objectSuffix...)
	return raw
}

// Error represents a RESP Error object
// Note: this is not Go's error.
type Error struct {
	content string
}

// NewError returns a new Error object with content
func NewError(content string) *Error {
	e := Error{content: content}
	return &e
}

// String returns the content string
func (e *Error) String() string {
	return e.content
}

// Dump returns raw bytes representation
func (e *Error) Dump() []byte {
	raw := []byte{byte(errorPrefix)}
	raw = append(raw, []byte(e.String())...)
	raw = append(raw, objectSuffix...)
	return raw
}

// Integer represents a RESP Integer object
type Integer struct {
	content int
}

// NewInteger returns a new Integer object with v
func NewInteger(v int) *Integer {
	i := Integer{content: v}
	return &i
}

// Int returns content as int
func (i *Integer) Int() int {
	return i.content
}

// Dump returns raw bytes representation
func (i *Integer) Dump() []byte {
	return []byte(fmt.Sprintf(":%d\r\n", i.Int()))
}

// BulkString represents a RESP Bulk String object
type BulkString struct {
	bytes []byte
}

// NewBulkString returns a new BulkString object whose length is size
func NewBulkString(dat []byte) *BulkString {
	return &BulkString{bytes: dat}
}

// NewBulkStringSize returns a new BulkString object whose length is size.
// The length is the bytes of internal binary data.
func NewBulkStringSize(size int) *BulkString {
	b := new(BulkString)
	b.bytes = nil
	if size > -1 {
		b.bytes = make([]byte, size)
	}
	return b
}

// String returns the content string
func (b *BulkString) String() string {
	return string(b.bytes)
}

// Dump returns raw bytes representation
func (b *BulkString) Dump() []byte {
	if b.bytes == nil {
		return []byte("$-1\r\n")
	}
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(b.bytes), string(b.bytes)))
}

// Length returns length of internal data in bytes
func (b *BulkString) Length() int {
	return len(b.bytes)
}

// Bytes returns internal data in byte array form
func (b *BulkString) Bytes() []byte {
	return b.bytes
}

// SetBytes sets internal data
func (b *BulkString) SetBytes(bytes []byte) {
	b.bytes = bytes
}

// Array represents a RESP Array object
type Array struct {
	objects []Object
}

// NewArraySize returns a new Array object whose length is size
func NewArraySize(size int) *Array {
	a := Array{}
	a.objects = make([]Object, size)
	return &a
}

// NewArray returns a new Array object filled by objs
func NewArray(objs ...Object) *Array {
	return &Array{objects: objs}
}

// Dump returns raw bytes representation
func (a *Array) Dump() []byte {
	data := []byte{arrayPrefix}
	data = append(data, []byte(fmt.Sprintf("%d%s", a.Length(), objectSuffix))...)
	for _, o := range a.objects {
		data = append(data, o.Dump()...)
	}
	return data
}

// Length returns the number of objects contained inside
func (a *Array) Length() int {
	return len(a.objects)
}

// SetObject sets i-th object to o.  An error will be returned if i is out
// of range.
func (a *Array) SetObject(i int, o Object) error {
	if len(a.objects) <= i {
		return fmt.Errorf("%d is out of range", i)
	}

	a.objects[i] = o
	return nil
}

// Objects returns internal data as an array of Objects
func (a *Array) Objects() []Object {
	return a.objects
}
