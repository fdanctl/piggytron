package handlers

import (
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/fdanctl/piggytron/internal/interface/http/assets"
)

const (
	// cacheControlImmutable is safe only for URLs whose path or query changes
	// whenever the contents change.
	cacheControlImmutable = "public, max-age=31536000, immutable"

	// cacheControlDefault is used for assets that can change without their
	// URL changing.
	cacheControlDefault = "public, max-age=3600"

	// cacheControlDev disables caching so a rebuild is visible immediately.
	cacheControlDev = "no-store"
)

type StaticHandler struct {
	static fs.FS
	isDev  bool
}

func NewStaticHandler(static fs.FS, isDev bool) *StaticHandler {
	return &StaticHandler{
		static: static,
		isDev:  isDev,
	}
}

type EncodingPrefs struct {
	Brotli   float64
	Gzip     float64
	Identity float64
}

func parseAcceptEncoding(values []string) EncodingPrefs {
	// No header: any content coding is acceptable.
	if values == nil {
		return EncodingPrefs{
			Brotli:   1,
			Gzip:     1,
			Identity: 1,
		}
	}

	// Header present but empty: don't use content codings.
	prefs := EncodingPrefs{
		Identity: 1,
	}

	var wildcard float64
	var hasWildcard bool

	var hasBrotli bool
	var hasGzip bool
	var hasIdentity bool

	for _, value := range values {
		for _, entry := range strings.Split(value, ",") {
			parts := strings.SplitN(strings.TrimSpace(entry), ";", 2)
			encoding := strings.TrimSpace(parts[0])

			if encoding == "" {
				continue
			}

			q := 1.0

			if len(parts) == 2 {
				params := strings.TrimSpace(parts[1])
				name, value, found := strings.Cut(params, "=")

				if !found || !strings.EqualFold(strings.TrimSpace(name), "q") {
					continue
				}

				parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
				if err != nil || parsed < 0 || parsed > 1 {
					continue
				}

				q = parsed
			}

			switch strings.ToLower(encoding) {
			case "br":
				prefs.Brotli = q
				hasBrotli = true

			case "gzip":
				prefs.Gzip = q
				hasGzip = true

			case "identity":
				prefs.Identity = q
				hasIdentity = true

			case "*":
				wildcard = q
				hasWildcard = true
			}
		}
	}

	// Apply wildcard only to encodings not explicitly listed.
	if hasWildcard {
		if !hasBrotli {
			prefs.Brotli = wildcard
		}

		if !hasGzip {
			prefs.Gzip = wildcard
		}

		// Identity is special: wildcard q=0 excludes it,
		// unless identity was explicitly listed.
		if wildcard == 0 && !hasIdentity {
			prefs.Identity = 0
		}
	}

	return prefs
}

func (h *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodHead:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *StaticHandler) Get(w http.ResponseWriter, r *http.Request) {
	fp := strings.TrimPrefix(r.URL.Path, "/")

	if !fs.ValidPath(fp) {
		http.NotFound(w, r)
		return
	}

	// Check that the original asset exists.
	if _, err := fs.Stat(h.static, fp); err != nil {
		http.NotFound(w, r)
		return
	}

	prefs := parseAcceptEncoding(r.Header.Values("Accept-Encoding"))

	// Add available representations in tie-break preference order.
	type candidate struct {
		encoding string
		path     string
		quality  float64
	}

	candidates := []candidate{
		{"br", fp + ".br", prefs.Brotli},
		{"gzip", fp + ".gz", prefs.Gzip},
		{"", fp, prefs.Identity},
	}

	bestQuality := -1.0
	var selected candidate

	for _, c := range candidates {
		if c.quality <= 0 || c.quality <= bestQuality {
			continue
		}

		if _, err := fs.Stat(h.static, c.path); err != nil {
			continue
		}

		selected = c
		bestQuality = c.quality
	}

	if bestQuality < 0 {
		http.Error(w, "Not Acceptable", http.StatusNotAcceptable)
		return
	}

	if contentType := h.contentType(fp); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Add("Vary", "Accept-Encoding")
	w.Header().Set("Cache-Control", h.cacheControl(fp, r))

	if selected.encoding != "" {
		w.Header().Set("Content-Encoding", selected.encoding)
	}

	http.ServeFileFS(w, r, h.static, selected.path)
}

// contentType returns the type of the identity representation. Sniffing the
// compressed bytes would name the content coding instead of the media type,
// so unknown extensions fall back to the original file.
func (h *StaticHandler) contentType(fp string) string {
	if contentType := mime.TypeByExtension(path.Ext(fp)); contentType != "" {
		return contentType
	}

	file, err := h.static.Open(fp)
	if err != nil {
		return ""
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return ""
	}

	return http.DetectContentType(buf[:n])
}

// cacheControl picks a policy from the kind of asset being served.
func (h *StaticHandler) cacheControl(fp string, r *http.Request) string {
	if h.isDev {
		return cacheControlDev
	}

	if isImmutableAsset(fp, r) {
		return cacheControlImmutable
	}

	return cacheControlDefault
}

// isImmutableAsset reports whether the asset is guaranteed never to change
// under the same URL.
func isImmutableAsset(fp string, r *http.Request) bool {
	base := path.Base(fp)

	switch base {
	case path.Base(assets.JSFile), path.Base(assets.CSSFile):
		return r.URL.Query().Get("v") != ""
	}

	// Third-party libraries carry their version in the file name.
	return strings.HasPrefix(base, "htmx-") || strings.HasPrefix(base, "echarts-")
}
