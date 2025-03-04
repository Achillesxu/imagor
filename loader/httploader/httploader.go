package httploader

import (
	"github.com/cshum/imagor"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

type HTTPLoader struct {
	// The Transport used to request images.
	// If nil, http.DefaultTransport is used.
	Transport http.RoundTripper

	// ForwardHeaders copy request headers to image request headers
	ForwardHeaders []string

	// OverrideHeaders override image request headers
	OverrideHeaders map[string]string

	// AllowedSources list of host names allowed to load from,
	// supports glob patterns such as *.google.com
	AllowedSources []string

	// MaxAllowedSize maximum bytes allowed for image
	MaxAllowedSize int

	// DefaultScheme default image URL scheme
	DefaultScheme string

	// UserAgent default user agent for image request.
	// Can be overridden by ForwardHeaders and OverrideHeaders
	UserAgent string
}

func New(options ...Option) *HTTPLoader {
	h := &HTTPLoader{
		OverrideHeaders: map[string]string{},
		DefaultScheme:   "https",
		UserAgent:       "Imagor",
	}
	for _, option := range options {
		option(h)
	}
	if s := strings.ToLower(h.DefaultScheme); s == "nil" {
		h.DefaultScheme = ""
	}
	return h
}

func (h *HTTPLoader) Load(r *http.Request, image string) ([]byte, error) {
	if r.Method != http.MethodGet || image == "" {
		return nil, imagor.ErrPass
	}
	u, err := url.Parse(image)
	if err != nil {
		return nil, imagor.ErrPass
	}
	if u.Host == "" || u.Scheme == "" {
		if h.DefaultScheme != "" {
			image = h.DefaultScheme + "://" + image
			if u, err = url.Parse(image); err != nil {
				return nil, imagor.ErrPass
			}
		} else {
			return nil, imagor.ErrPass
		}
	}
	if !h.isUrlAllowed(u) {
		return nil, imagor.ErrPass
	}
	client := &http.Client{Transport: h.Transport}
	if h.MaxAllowedSize > 0 {
		req, err := h.newRequest(r, http.MethodHead, image)
		if err != nil {
			return nil, err
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, imagor.ErrPass
		}
		_ = resp.Body.Close()
		if resp.StatusCode < 200 && resp.StatusCode > 206 {
			return nil, imagor.NewErrorFromStatusCode(resp.StatusCode)
		}
		contentLength, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
		if contentLength > h.MaxAllowedSize {
			return nil, imagor.ErrMaxSizeExceeded
		}
	}
	req, err := h.newRequest(r, http.MethodGet, image)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, imagor.ErrPass
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return buf, imagor.NewErrorFromStatusCode(resp.StatusCode)
	}
	return buf, nil
}

func (h *HTTPLoader) newRequest(r *http.Request, method, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(r.Context(), method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", h.UserAgent)
	for _, header := range h.ForwardHeaders {
		if header == "*" {
			req.Header = r.Header.Clone()
			req.Header.Del("Accept-Encoding") // fix compressions
			break
		}
		if _, ok := r.Header[header]; ok {
			req.Header.Set(header, r.Header.Get(header))
		}
	}
	for key, value := range h.OverrideHeaders {
		req.Header.Set(key, value)
	}
	return req, nil
}

func (h *HTTPLoader) isUrlAllowed(u *url.URL) bool {
	if len(h.AllowedSources) == 0 {
		return true
	}
	for _, source := range h.AllowedSources {
		if matched, e := path.Match(source, u.Host); matched && e == nil {
			return true
		}
	}
	return false
}
