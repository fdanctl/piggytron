package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func staticTestFS() fstest.MapFS {
	return fstest.MapFS{
		"static/app.js":                {Data: []byte("console.log(1)")},
		"static/app.js.br":             {Data: []byte("brotli-js")},
		"static/app.js.gz":             {Data: []byte("gzip-js")},
		"static/app.css":               {Data: []byte("body{color:red}")},
		"static/app.css.gz":            {Data: []byte("gzip-css")},
		"static/htmx-4.0.0.min.js":     {Data: []byte("htmx")},
		"static/icone.svg":             {Data: []byte("<svg/>")},
		"static/asset.templar-test":    {Data: []byte("%PDF-1.4 fake")},
		"static/asset.templar-test.br": {Data: []byte("<html>not a pdf</html>")},
	}
}

func acceptEncoding(value string) *string {
	return &value
}

func TestStaticHandlerNegotiation(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		acceptEncoding  *string
		wantStatus      int
		wantEncoding    string
		wantBody        string
		wantContentType string
	}{
		{
			name:            "absent header prefers brotli",
			path:            "/static/app.js",
			wantStatus:      http.StatusOK,
			wantEncoding:    "br",
			wantBody:        "brotli-js",
			wantContentType: "text/javascript; charset=utf-8",
		},
		{
			name:           "brotli wins on a higher quality",
			path:           "/static/app.js",
			acceptEncoding: acceptEncoding("gzip;q=0.5, br;q=0.9, identity;q=0"),
			wantStatus:     http.StatusOK,
			wantEncoding:   "br",
			wantBody:       "brotli-js",
		},
		{
			name:           "gzip wins on a higher quality",
			path:           "/static/app.js",
			acceptEncoding: acceptEncoding("gzip;q=0.9, br;q=0.5, identity;q=0"),
			wantStatus:     http.StatusOK,
			wantEncoding:   "gzip",
			wantBody:       "gzip-js",
		},
		{
			name:           "brotli disabled falls back to gzip",
			path:           "/static/app.js",
			acceptEncoding: acceptEncoding("br;q=0, gzip"),
			wantStatus:     http.StatusOK,
			wantEncoding:   "gzip",
			wantBody:       "gzip-js",
		},
		{
			name:           "missing brotli representation falls back to identity",
			path:           "/static/app.css",
			acceptEncoding: acceptEncoding("br"),
			wantStatus:     http.StatusOK,
			wantEncoding:   "",
			wantBody:       "body{color:red}",
		},
		{
			name:           "zero wildcard excludes every encoding",
			path:           "/static/app.css",
			acceptEncoding: acceptEncoding("*;q=0"),
			wantStatus:     http.StatusNotAcceptable,
		},
		{
			name:           "identity and the only available coding disabled",
			path:           "/static/app.css",
			acceptEncoding: acceptEncoding("identity;q=0, gzip;q=0"),
			wantStatus:     http.StatusNotAcceptable,
		},
		{
			name:           "invalid quality is ignored",
			path:           "/static/app.js",
			acceptEncoding: acceptEncoding("gzip;q=abc"),
			wantStatus:     http.StatusOK,
			wantEncoding:   "",
			wantBody:       "console.log(1)",
		},
		{
			name:            "unknown extension is sniffed from the identity file",
			path:            "/static/asset.templar-test",
			acceptEncoding:  acceptEncoding("br"),
			wantStatus:      http.StatusOK,
			wantEncoding:    "br",
			wantBody:        "<html>not a pdf</html>",
			wantContentType: "application/pdf",
		},
		{
			name:       "missing original returns not found",
			path:       "/static/missing.js",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "path traversal returns not found",
			path:       "/static/../main.go",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "double slash returns not found",
			path:       "/static//etc/passwd",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewStaticHandler(staticTestFS(), false)

			r := httptest.NewRequest(http.MethodGet, tt.path, nil)
			if tt.acceptEncoding != nil {
				r.Header.Set("Accept-Encoding", *tt.acceptEncoding)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			if got := w.Header().Get("Content-Encoding"); got != tt.wantEncoding {
				t.Errorf("Content-Encoding = %q, want %q", got, tt.wantEncoding)
			}

			if got := w.Body.String(); got != tt.wantBody {
				t.Errorf("body = %q, want %q", got, tt.wantBody)
			}

			if tt.wantContentType != "" {
				if got := w.Header().Get("Content-Type"); got != tt.wantContentType {
					t.Errorf("Content-Type = %q, want %q", got, tt.wantContentType)
				}
			}

			if got := w.Header().Get("X-Content-Type-Options"); got != "nosniff" {
				t.Errorf("X-Content-Type-Options = %q, want nosniff", got)
			}

			if got := w.Header().Get("Vary"); got != "Accept-Encoding" {
				t.Errorf("Vary = %q, want Accept-Encoding", got)
			}
		})
	}
}

func TestStaticHandlerHead(t *testing.T) {
	handler := NewStaticHandler(staticTestFS(), false)

	r := httptest.NewRequest(http.MethodHead, "/static/app.js", nil)
	r.Header.Set("Accept-Encoding", "br")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if got := w.Header().Get("Content-Encoding"); got != "br" {
		t.Errorf("Content-Encoding = %q, want br", got)
	}
	if w.Body.Len() != 0 {
		t.Errorf("HEAD body = %q, want empty", w.Body.String())
	}
}

func TestStaticHandlerMethodNotAllowed(t *testing.T) {
	handler := NewStaticHandler(staticTestFS(), false)

	r := httptest.NewRequest(http.MethodPost, "/static/app.js", nil)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestStaticHandlerCacheControl(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		isDev bool
		want  string
	}{
		{
			name:  "dev disables caching",
			path:  "/static/app.js",
			isDev: true,
			want:  cacheControlDev,
		},
		{
			name: "versioned bundle is immutable",
			path: "/static/app.js?v=abc123",
			want: cacheControlImmutable,
		},
		{
			name: "unversioned bundle revalidates sooner",
			path: "/static/app.js",
			want: cacheControlDefault,
		},
		{
			name: "versioned library is immutable",
			path: "/static/htmx-4.0.0.min.js",
			want: cacheControlImmutable,
		},
		{
			name: "unversioned asset uses the default",
			path: "/static/icone.svg",
			want: cacheControlDefault,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewStaticHandler(staticTestFS(), tt.isDev)

			r := httptest.NewRequest(http.MethodGet, tt.path, nil)

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
			}

			if got := w.Header().Get("Cache-Control"); got != tt.want {
				t.Errorf("Cache-Control = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseAcceptEncoding(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		want   EncodingPrefs
	}{
		{
			name:   "absent header accepts everything",
			values: nil,
			want:   EncodingPrefs{Brotli: 1, Gzip: 1, Identity: 1},
		},
		{
			name:   "empty header disables content codings",
			values: []string{""},
			want:   EncodingPrefs{Identity: 1},
		},
		{
			name:   "explicit qualities",
			values: []string{"gzip;q=0.5, br;q=0.8"},
			want:   EncodingPrefs{Brotli: 0.8, Gzip: 0.5, Identity: 1},
		},
		{
			name:   "wildcard fills unlisted codings",
			values: []string{"gzip, *;q=0.5"},
			want:   EncodingPrefs{Brotli: 0.5, Gzip: 1, Identity: 1},
		},
		{
			name:   "zero wildcard excludes identity",
			values: []string{"*;q=0"},
			want:   EncodingPrefs{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseAcceptEncoding(tt.values); got != tt.want {
				t.Errorf("parseAcceptEncoding(%q) = %#v, want %#v", tt.values, got, tt.want)
			}
		})
	}
}

func TestIsImmutableAsset(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "versioned js bundle", path: "/static/app.js?v=abc", want: true},
		{name: "versioned css bundle", path: "/static/app.css?v=abc", want: true},
		{name: "unversioned js bundle", path: "/static/app.js", want: false},
		{name: "versioned library", path: "/static/htmx-4.0.0.min.js", want: true},
		{name: "unversioned asset", path: "/static/icone.svg", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tt.path, nil)
			fp := strings.TrimPrefix(r.URL.Path, "/")

			if got := isImmutableAsset(fp, r); got != tt.want {
				t.Errorf("isImmutableAsset(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
