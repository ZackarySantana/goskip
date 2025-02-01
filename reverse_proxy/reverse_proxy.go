package skip_reverse_proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// Options are high level options for the reverse proxy.
type Options struct {
	// GetUUIDFromPath is a function that takes a path and returns the UUID.
	// If not provided, the default implementation will extract the UUID from the path.
	GetUUIDFromPath func(path string) string
}

// WithGetUUIDFromPath sets the function to extract the UUID from the path.
func WithGetUUIDFromPath(getUUIDFromPath func(path string) string) func(*Options) {
	return func(o *Options) {
		o.GetUUIDFromPath = getUUIDFromPath
	}
}

// New creates a reverse proxy that redirects requests
// served to it to the given URL. The given URL must have a path
// that has a placeholder for the UUID (e.g. %s). For most Skip
// services, this will be "/v1/streams/%s".
func New(url *url.URL, options ...func(*Options)) (*httputil.ReverseProxy, error) {
	opts := &Options{
		GetUUIDFromPath: func(path string) string {
			streamsIndex := strings.Index(path, "/streams/")
			if streamsIndex == -1 {
				return "invalid_path"
			}
			if len(path) <= streamsIndex+9 {
				return "invalid_path"
			}
			return path[streamsIndex+9:]
		},
	}

	for _, o := range options {
		o(opts)
	}

	rp := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetXForwarded()

			uuid := opts.GetUUIDFromPath(pr.In.URL.Path)

			url := *url
			pr.Out.URL = &url
			pr.Out.URL.Path = fmt.Sprintf(pr.Out.URL.Path, uuid)
		},
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Minute,
				KeepAlive: 30 * time.Minute,
			}).Dial,
			TLSHandshakeTimeout: 60 * time.Second,
		},
		FlushInterval: 100 * time.Millisecond,
	}

	return rp, nil
}
