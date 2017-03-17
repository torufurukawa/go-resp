package resp

import "testing"
import "bytes"

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
	if obj != "OK" {
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

func TestParse(t *testing.T) {
	line := []byte("+OK")
	result, err := parse(line)
	if err != nil {
		t.Errorf("err is %#v, want nil", err)
	}
	if result != "OK" {
		t.Errorf("result is %#v, want OK", result)
	}
}
