// Copyright 2018 Paul Greenberg (greenpau@outlook.com)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exporter

import (
	"testing"

	"github.com/prometheus/common/promlog"
)

func TestNewExporter(t *testing.T) {
	allowedLogLevel := &promlog.AllowedLevel{}
	if err := allowedLogLevel.Set("debug"); err != nil {
		t.Fatalf("%s", err)
	}

	promlogConfig := &promlog.Config{
		Level: allowedLogLevel,
	}

	logger := promlog.New(promlogConfig)

	cases := []struct {
		address string
		ok      bool
	}{
		{address: "127.0.0.1:50051", ok: false},
		{address: "", ok: false},
		{address: "127.0.0.1:500511", ok: false},
		{address: "localaddress:50051", ok: false},
		{address: "http://localaddress:50051", ok: false},
		{address: "fuuuu://localaddress:50051", ok: false},
		{address: "dns:///localhost:50051", ok: false},
		{address: "[::1]:50051", ok: false},
		{address: "::1:50051", ok: false},
	}
	pollTimeout := 2
	for _, test := range cases {
		opts := Options{
			Timeout: pollTimeout,
			Address: test.address,
			Logger:  logger,
		}
		_, err := NewExporter(opts)
		if test.ok && err != nil {
			t.Errorf("expected no error w/ %q, but got %q", test.address, err)
		}
		if !test.ok && err == nil {
			t.Errorf("expected error w/ %q, but got %q", test.address, err)
		}
	}
}
