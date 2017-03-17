package resp

import (
	"bytes"
	"reflect"
	"testing"
)

func TestWrite(t *testing.T) {
	fixtures := []struct {
		obj      Object
		expected []byte
	}{
		{obj: NewSimpleString("OK"), expected: []byte("+OK\r\n")},
		{obj: NewError("ERR foo"), expected: []byte("-ERR foo\r\n")},
		{obj: NewInteger(123), expected: []byte(":123\r\n")},
		{obj: func() *BulkString {
			b := BulkString([]byte("PING"))
			return &b
		}(), expected: []byte("$4\r\nPING\r\n")},
	}

	for _, f := range fixtures {
		buf := new(bytes.Buffer)
		writer := NewWriter(buf)

		err := writer.WriteObject(f.obj)
		if err != nil {
			t.Errorf("err is %#v, want nil", err)
		}

		result := buf.Bytes()
		if !reflect.DeepEqual(result, f.expected) {
			t.Errorf("%#v is written, want %#v", result, f.expected)
		}
	}
}
