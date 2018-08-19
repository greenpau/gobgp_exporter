package main

import "testing"

func TestNewExporter(t *testing.T) {
	cases := []struct {
		address string
		ok      bool
	}{
		{address: "", ok: false},
		{address: "localaddress:50051", ok: false},
		{address: "127.0.0.1:50051", ok: true},
		{address: "http://localaddress:50051", ok: false},
		{address: "fuuuu://localaddress:50051", ok: false},
	}

	for _, test := range cases {
		_, err := NewExporter(gobgpOpts{address: test.address})
		if test.ok && err != nil {
			t.Errorf("expected no error w/ %q, but got %q", test.address, err)
		}
		if !test.ok && err == nil {
			t.Errorf("expected error w/ %q, but got %q", test.address, err)
		}
	}
}
