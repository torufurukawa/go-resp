package resp

import (
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
