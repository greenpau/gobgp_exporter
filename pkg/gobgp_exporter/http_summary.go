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
	"net/http"
	"strings"
)

// Summary returns the content of the Exporter's default page.
func (e *Exporter) Summary(p string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate, no-store")
	token, authorized := e.authorize(r)
	if !authorized {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	var sb strings.Builder
	sb.WriteString(`<html>`)
	sb.WriteString(`<head><title>Prometheus Exporter for GoBGP</title></head>`)
	sb.WriteString(`<body>`)
	sb.WriteString(`<h1>Prometheus Exporter for GoBGP</h1>`)
	sb.WriteString(`<table border='1'>`)
	sb.WriteString(`<tr>`)
	sb.WriteString(`<th>Node</th>`)
	sb.WriteString(`<th>Last Result</th>`)
	sb.WriteString(`<th>Last Scrape</th>`)
	sb.WriteString(`<th>Metrics</th><tr>`)
	n := e.Node
	url := p + `?x-token=` + token
	sb.WriteString(`<tr>`)
	sb.WriteString(`<td>` + n.address + `</td>`)
	switch n.result {
	case "success":
		sb.WriteString(`<td style="background-color:lightgreen">` + n.result + `</td>`)
	case "failure":
		sb.WriteString(`<td style="background-color:tomato">` + n.result + `</td>`)
	default:
		sb.WriteString(`<td style="background-color:lightgray">` + n.result + `</td>`)
	}
	sb.WriteString(`<td>` + n.timestamp + `</td>`)
	sb.WriteString(`<td><a href='` + url + `'>Metrics</a></td>`)
	sb.WriteString(`</tr>`)
	sb.WriteString(`</table>`)
	sb.WriteString(`</body>`)
	sb.WriteString(`</html>`)
	w.Write([]byte(sb.String())) //nolint:errcheck
}
