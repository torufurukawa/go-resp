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
	// given: reader is ready to read large data
	data := append(make([]byte, bufferSize), []byte("\r\n")...)
	buf := bytes.NewReader(data)
	reader := NewReader(buf)

	// when: call ReadObject()
	_, err := reader.ReadObject()
	// then: error
	if err == nil {
		t.Errorf("err is nil, want error")
		return
	}
}
