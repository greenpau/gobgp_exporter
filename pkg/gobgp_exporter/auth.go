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
	"fmt"
	"github.com/prometheus/common/log"
	"net/http"
)

// AddAuthenticationToken adds an authentication token for accessing
// the exporter itself.
func (e *Exporter) AddAuthenticationToken(s string) error {
	if s == "" {
		return fmt.Errorf("invalid empty token")
	}
	e.Tokens[s] = true
	return nil
}

func (e *Exporter) authorize(r *http.Request) (string, bool) {
	if _, exists := e.Tokens["anonymous"]; exists {
		return "anonymous", true
	}
	var invalidToken bool
	tokens := []string{"x_token", "x-token", "X-Token"}
	for _, t := range tokens {
		token := r.Header.Get(t)
		if token != "" {
			if _, exists := e.Tokens[token]; exists {
				return token, true
			}
			invalidToken = true
		}
	}

	for _, t := range tokens {
		token := r.URL.Query().Get(t)
		if token != "" {
			if _, exists := e.Tokens[token]; exists {
				return token, true
			}
			invalidToken = true
		}
	}

	if invalidToken {
		log.Warnf("unauthorized access from %q due to invalid token", r.RemoteAddr)
	} else {
		log.Warnf("unauthorized access from %q due to the lack of auth token", r.RemoteAddr)
	}

	return "", false
}
