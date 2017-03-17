package resp

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestSimpleString(t *testing.T) {
	content := "OK"
	s := NewSimpleString(content)

	if string(*s) != content {
		t.Errorf("s is %#v, want %#v", s, content)
	}

	dump := s.Dump()
	expectedDump := []byte("+OK\r\n")
	if !reflect.DeepEqual(dump, expectedDump) {
		t.Errorf("s.Dump() is \n%#v, want \n%#v", dump, expectedDump)
	}
}

func TestError(t *testing.T) {
	content := "ERR foo bar"
	s := NewError(content)

	if string(*s) != content {
		t.Errorf("s is %#v, want %#v", s, content)
	}

	dump := s.Dump()
	expectedDump := []byte(fmt.Sprintf("-%s\r\n", content))
	if !reflect.DeepEqual(dump, expectedDump) {
		t.Errorf("s.Dump() is \n%#v, want \n%#v", dump, expectedDump)
	}
}

func TestWrite(t *testing.T) {
	fixtures := []struct {
		obj      Object
		expected []byte
	}{
		{obj: NewSimpleString("OK"), expected: []byte("+OK\r\n")},
		{obj: NewError("ERR foo"), expected: []byte("-ERR foo\r\n")},
		{obj: NewInteger(123), expected: []byte(":123\r\n")},
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

// TODO: test integer
func TestReadSuccess(t *testing.T) {
	fixtures := []struct {
		input []byte
		obj   Object
	}{
		{input: []byte("+OK\r\n"), obj: NewSimpleString("OK")},
		{input: []byte("-ERR\r\n"), obj: NewError("ERR")},
		{input: []byte(":123\r\n"), obj: NewInteger(123)},
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
		[]byte("\r\n"),
		[]byte("#\r\n"),
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
		obj  Object
		eq   func(a interface{}, b interface{}) bool
	}{
		{line: []byte("+OK"), obj: NewSimpleString("OK"), eq: reflect.DeepEqual},
		{line: []byte("-ERR"), obj: NewError("ERR"), eq: reflect.DeepEqual},
		{line: []byte(":123"), obj: NewInteger(123), eq: reflect.DeepEqual},
	}

	for _, f := range fixture {
		result, err := parseLine(f.line)
		if err != nil {
			t.Errorf("err is %#v, want nil", err)
		}
		if !f.eq(result, f.obj) {
			t.Errorf("result is %#v, want %#v", result, f.obj)
		}
	}
}
