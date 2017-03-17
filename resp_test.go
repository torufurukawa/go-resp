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
