package resp

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSimpleString(t *testing.T) {
	content := "OK"
	s := NewSimpleString(content)

	if s.String() != content {
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

	if s.String() != content {
		t.Errorf("s is %#v, want %#v", s, content)
	}

	dump := s.Dump()
	expectedDump := []byte(fmt.Sprintf("-%s\r\n", content))
	if !reflect.DeepEqual(dump, expectedDump) {
		t.Errorf("s.Dump() is \n%#v, want \n%#v", dump, expectedDump)
	}
}

func TestBulkString(t *testing.T) {
	v := []byte("foo")
	s := NewBulkString(v)
	fixtures := []struct {
		result   []byte
		desc     string
		expected []byte
	}{
		{result: s.bytes, expected: v, desc: "s.bytes"},
		{result: s.Bytes(), expected: v, desc: "s.Bytes()"},
	}
	for _, f := range fixtures {
		fmt.Printf(f.desc)
		if string(f.result) != string(f.expected) {
			t.Errorf("%s is %#v, want %#v", f.desc, f.result, f.expected)
		}
	}
}

func TestBulkStringString(t *testing.T) {
	v := "foo"
	s := NewBulkString([]byte(v))
	result := s.String()
	if result != v {
		t.Errorf("s.String() returns %#v, want %#v", result, v)
	}
}

func TestNullBulkString(t *testing.T) {
	s := NewBulkString(nil)
	fixtures := []struct {
		result []byte
		desc   string
	}{
		{result: s.bytes, desc: "s.bytes"},
		{result: s.Bytes(), desc: "s.Bytes()"},
	}
	for _, f := range fixtures {
		if f.result != nil {
			t.Errorf("%s is %#v, want nil", f.desc, f.result)
		}
	}
}
