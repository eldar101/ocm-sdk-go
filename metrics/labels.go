/*
Copyright (c) 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file contains functions that calculate the labels included in metrics.

package metrics

import (
	"net/http"
	"strconv"
	"strings"
)

// serviceLabel calculates the `service` for the given HTTP request.
func (t *roundTripper) serviceLabel(request *http.Request) string {
	path := request.URL.Path
	if !strings.HasPrefix(path, "/api/") {
		return ""
	}
	if strings.HasPrefix(path, "/api/accounts_mgmt") {
		return "ocm-accounts-service"
	} else if strings.HasPrefix(path, "/api/clusters_mgmt") {
		return "ocm-clusters-service"
	} else if strings.HasPrefix(path, "/api/authorizations") {
		return "ocm-authorizations-service"
	} else if strings.HasPrefix(path, "/api/service_logs") {
		return "ocm-logs-service"
	} else {
		parts := strings.Split(path, "/")
		if len(parts) > 3 {
			return "ocm-" + parts[3]
		}
		return ""
	}
}

// methodLabel calculates the `method` label from the HTTP method.
func (t *roundTripper) methodLabel(request *http.Request) string {
	return strings.ToUpper(request.Method)
}

// pathLabel calculates the `path` label from the URL path.
func (t *roundTripper) pathLabel(request *http.Request) string {
	// Remove leading and trailing slashes:
	path := request.URL.Path
	for len(path) > 0 && strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	for len(path) > 0 && strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-1]
	}

	// Clear segments that correspond to path variables:
	segments := strings.Split(path, "/")
	current := t.owner.paths
	for i, segment := range segments {
		next, ok := current[segment]
		if ok {
			current = next
			continue
		}
		next, ok = current["-"]
		if ok {
			segments[i] = "-"
			current = next
			continue
		}
		return "/-"
	}

	// Reconstruct the path joining the modified segments:
	return "/" + strings.Join(segments, "/")
}

// codeLabel calculates the `code` label from the given HTTP response.
func (t *roundTripper) codeLabel(response *http.Response) string {
	code := 0
	if response != nil {
		code = response.StatusCode
	}
	return strconv.Itoa(code)
}

// Names of the labels added to metrics:
const (
	serviceLabelName = "apiservice"
	codeLabelName    = "code"
	methodLabelName  = "method"
	pathLabelName    = "path"
)

// Array of labels added to call metrics:
var requestLabelNames = []string{
	serviceLabelName,
	codeLabelName,
	methodLabelName,
	pathLabelName,
}
