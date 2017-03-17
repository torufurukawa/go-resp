package resp

import "testing"
import "bytes"
import "reflect"

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

func TestWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf)
	err := writer.WriteObject(NewSimpleString("OK"))
	if err != nil {
		t.Errorf("err is %#v, want nil", err)
	}

	expected := []byte("+OK\r\n")
	if !reflect.DeepEqual(buf.Bytes(), expected) {
		t.Errorf("buf is %#v, want %#v", buf, expected)
	}
}

func TestRead(t *testing.T) {
	// given: reader is ready to read "+OK\r\n"
	buf := bytes.NewReader([]byte("+OK\r\n"))
	reader := NewReader(buf)

	// when: call ReadObject()
	obj, err := reader.ReadObject()
	if err != nil {
		t.Errorf("err is %#v, want nil", err)
		return
	}

	// then: "OK" is returned
	if *(obj.(*SimpleString)) != SimpleString("OK") {
		t.Errorf("obj is %#v, want OK", obj)
	}
}

func TestReadError(t *testing.T) {
	inputs := [][]byte{
		append(make([]byte, bufferSize), []byte("\r\n")...),
		[]byte("\r\n"),
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
	line := []byte("+OK")
	result, err := parseLine(line)
	if err != nil {
		t.Errorf("err is %#v, want nil", err)
	}
	if *(result.(*SimpleString)) != SimpleString("OK") {
		t.Errorf("result is %#v, want OK", result)
	}
}
