// Package assets derives cache-busting URLs for the static assets embedded in
// the server binary.
package assets

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"io/fs"
)

// FS holds the static files served by every binary, so the assets are shared
// by every command instead of copied into each cmd directory.
//
//go:embed static/*
var FS embed.FS

// Paths of the generated bundles inside the embedded file system.
const (
	JSFile  = "static/app.js"
	CSSFile = "static/app.css"
)

// URLs holds the URLs of the static assets referenced by the page layout.
type URLs struct {
	JS  string
	CSS string
}

// New computes content-addressed URLs for the given assets.
//
// The version query parameter is derived from the contents of the file, so it
// changes if and only if the bytes change. That makes it safe to cache the
// assets forever: a new build is served from a new URL.
func New(static fs.FS) URLs {
	return URLs{
		JS:  versionedURL(static, JSFile),
		CSS: versionedURL(static, CSSFile),
	}
}

func versionedURL(static fs.FS, name string) string {
	return "/" + name + "?v=" + contentHash(static, name)
}

// contentHash returns a short digest of the asset contents.
//
// A missing asset yields a stable marker instead of a random value: the
// handler will answer 404 either way, but the URL stays deterministic so the
// failure is easy to spot in the rendered page.
func contentHash(static fs.FS, name string) string {
	content, err := fs.ReadFile(static, name)
	if err != nil {
		return "missing"
	}

	sum := sha256.Sum256(content)

	return hex.EncodeToString(sum[:])[:10]
}
