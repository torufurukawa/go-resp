package resp

import (
	"bytes"
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	fixtures := []struct {
		input []byte
		obj   Object
	}{
		{input: []byte("+OK\r\n"), obj: NewSimpleString("OK")},
		{input: []byte("-ERR\r\n"), obj: NewError("ERR")},
		{input: []byte(":123\r\n"), obj: NewInteger(123)},
		{input: []byte("$4\r\nPING\r\n"),
			obj: func() *BulkString {
				b := BulkString([]byte("PING"))
				return &b
			}()},
	}

	for _, f := range fixtures {
		reader := NewReader(bytes.NewReader(f.input))
		obj, err := reader.ReadObject()
		if err != nil {
			t.Errorf("err is %#v, want nil", err)
		}
		if !reflect.DeepEqual(obj, f.obj) {
			t.Errorf("obj is %#v, want %#v", obj, f.obj)
		}
	}
}

func TestReadError(t *testing.T) {
	inputs := [][]byte{
		append(make([]byte, bufferSize), []byte("\r\n")...),
		[]byte("\r\n"),               // no data
		[]byte("#\r\n"),              // invalid prefix
		[]byte("$4\r\nP\r\n"),        // too short
		[]byte("$4\r\nPI\r\n"),       // too short
		[]byte("$4\r\nPINGPING\r\n"), // too long
	}

	for _, input := range inputs {
		reader := NewReader(bytes.NewReader(input))
		_, err := reader.ReadObject()
		if err == nil {
			t.Errorf("err is nil, want non-nil error")
		}
	}
}

func TestParseLine(t *testing.T) {
	fixture := []struct {
		line []byte
		obj  interface{}
		eq   func(a interface{}, b interface{}) bool
	}{
		{line: []byte("+OK"), obj: NewSimpleString("OK"), eq: reflect.DeepEqual},
		{line: []byte("-ERR"), obj: NewError("ERR"), eq: reflect.DeepEqual},
		{line: []byte(":123"), obj: NewInteger(123), eq: reflect.DeepEqual},
		{line: []byte("$4"), obj: NewBulkString(4), eq: reflect.DeepEqual},
	}

	r := NewReader(new(bytes.Buffer))
	for _, f := range fixture {
		result, err := r.parseLine(f.line)
		if err != nil {
			t.Errorf("err is %#v, want nil", err)
		}
		if !f.eq(result, f.obj) {
			t.Errorf("result is %#v, want %#v", result, f.obj)
		}
	}
}
