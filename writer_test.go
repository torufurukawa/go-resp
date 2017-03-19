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
		{obj: NewBulkString([]byte("PING")), expected: []byte("$4\r\nPING\r\n")},
		{obj: NewBulkString([]byte("")), expected: []byte("$0\r\n\r\n")},
		{obj: NewArray(
			NewBulkString([]byte("LLEN")),
			NewBulkString([]byte("mylist"))),
			expected: []byte("*2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n")},
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
			t.Errorf("\n%#v is written, want \n%#v", result, f.expected)
		}
	}
}
