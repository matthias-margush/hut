package debug

import (
	"net/http"
	"net/http/httputil"
)

func DebugResponse(resp *http.Response, includeBody bool) string {
	dump, _ := httputil.DumpResponse(resp, includeBody)
	return string(dump)
}

func DebugRequest(req *http.Request, includeBody bool) string {
	dump, _ := httputil.DumpRequest(req, includeBody)
	return string(dump)
}
