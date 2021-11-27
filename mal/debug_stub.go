//go:build !debug

package mal

import (
	"net/http"
)

func dumpRequest(req *http.Request) {}

func dumpResponse(resp *http.Response) {}
