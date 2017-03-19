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
		{input: []byte("$4\r\nPING\r\n"), obj: NewBulkString([]byte("PING"))},
		{input: []byte("$0\r\n\r\n"), obj: NewBulkString([]byte(""))},
		{input: []byte("$-1\r\n"), obj: NewBulkString(nil)},
		{input: []byte("*1\r\n$4\r\nPING\r\n"),
			obj: NewArray(NewBulkString([]byte("PING")))},
		{input: []byte("*2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n"),
			obj: NewArray(
				NewBulkString([]byte("LLEN")),
				NewBulkString([]byte("mylist")))},
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
	}{
		{line: []byte("+OK"), obj: NewSimpleString("OK")},
		{line: []byte("-ERR"), obj: NewError("ERR")},
		{line: []byte(":123"), obj: NewInteger(123)},
		{line: []byte("$4"), obj: NewBulkStringSize(4)},
		{line: []byte("*2"), obj: NewArraySize(2)},
	}

	r := NewReader(new(bytes.Buffer))
	for _, f := range fixture {
		result, err := r.parseLine(f.line)
		if err != nil {
			t.Errorf("err is %#v, want nil", err)
		}
		if !reflect.DeepEqual(result, f.obj) {
			t.Errorf("result is %#v, want %#v", result, f.obj)
		}
	}
}
