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
	"github.com/prometheus/client_golang/prometheus"
	"sort"
	"strings"
)

// GetMetricsTable returns markdown-formatted table with the name, help, and
// labels of the metrics produced by the exporter.
func (n *RouterNode) GetMetricsTable() string {
	var out strings.Builder
	metrics := make(chan *prometheus.Desc, 100000)
	n.Describe(metrics)
	close(metrics)
	out.WriteString("| **Metric** | **Description** | **Labels** |\n")
	out.WriteString("| ------ | ------- | ------ |\n")
	for metric := range metrics {
		m := strings.Split(strings.Replace(metric.String(), "Desc{", "", 1), ":")
		descr := make(map[string]string)
		var k string
		for _, e := range m {
			if k == "variableLabels" {
				descr[k] = strings.TrimSpace(e)
				descr[k] = strings.Trim(descr[k], "}")
				descr[k] = strings.Trim(descr[k], "]")
				descr[k] = strings.Trim(descr[k], "[")

				break
			}
			switch {
			case strings.HasSuffix(e, "fqName"):
				k = "name"
			case strings.HasSuffix(e, ", help"):
				descr[k] = strings.Replace(e, ", help", "", 1)
				descr[k] = strings.TrimSpace(descr[k])
				descr[k] = strings.Trim(descr[k], `"`)
				k = "help"
			case strings.HasSuffix(e, ", constLabels"):
				descr[k] = strings.Replace(e, ", constLabels", "", 1)
				descr[k] = strings.TrimSpace(descr[k])
				descr[k] = strings.Trim(descr[k], `"`)
				k = "constLabels"
			case strings.HasSuffix(e, ", variableLabels"):
				descr[k] = strings.Replace(e, ", variableLabels", "", 1)
				descr[k] = strings.TrimSpace(descr[k])
				descr[k] = strings.Trim(descr[k], `"`)
				k = "variableLabels"
			default:
				out.WriteString(e + "\n")
			}
		}

		labels := []string{}
		for _, label := range strings.Split(descr[k], " ") {
			if label == "" {
				continue
			}
			labels = append(labels, label)
		}
		sort.Strings(labels)
		if len(labels) > 0 {
			out.WriteString(fmt.Sprintf("| `%s` | %s | `%s` |\n", descr["name"], descr["help"], strings.Join(labels, "`, `")))
		} else {
			out.WriteString(fmt.Sprintf("| `%s` | %s | |\n", descr["name"], descr["help"]))
		}
	}
	return out.String()
}
