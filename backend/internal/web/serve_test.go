package web

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

const testIndexHTML = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>SunApp</title>
  </head>
  <body>
    <div id="root"></div>
  </body>
</html>`

func testFS() fstest.MapFS {
	return fstest.MapFS{
		"index.html":           {Data: []byte(testIndexHTML)},
		"favicon.svg":          {Data: []byte("<svg></svg>")},
		"assets/index-abc.js":  {Data: []byte("console.log('app')")},
	}
}

func TestBuildIndexHTML_InjectsLogoLinkUrl(t *testing.T) {
	t.Setenv("LOGO_LINK_URL", "https://example.org/")

	out, err := buildIndexHTML(testFS())
	if err != nil {
		t.Fatalf("buildIndexHTML returned error: %v", err)
	}

	want := `window.__APP_CONFIG__={"logoLinkUrl":"https://example.org/"}`
	if !strings.Contains(string(out), want) {
		t.Fatalf("expected injected config %q in output, got: %s", want, out)
	}

	// The script must be injected before </head>, preserving the rest of the document.
	if !strings.Contains(string(out), "</head>") {
		t.Fatalf("expected </head> to remain in output")
	}
	if !strings.Contains(string(out), "<div id=\"root\"></div>") {
		t.Fatalf("expected body content to remain in output")
	}
}

func TestBuildIndexHTML_EmptyEnvVar(t *testing.T) {
	t.Setenv("LOGO_LINK_URL", "")

	out, err := buildIndexHTML(testFS())
	if err != nil {
		t.Fatalf("buildIndexHTML returned error: %v", err)
	}

	want := `window.__APP_CONFIG__={"logoLinkUrl":""}`
	if !strings.Contains(string(out), want) {
		t.Fatalf("expected empty config %q in output, got: %s", want, out)
	}
}

func TestBuildIndexHTML_EscapesHTMLChars(t *testing.T) {
	// A value containing HTML-breaking characters must be escaped by
	// json.Marshal so it cannot break out of the <script> element.
	t.Setenv("LOGO_LINK_URL", `</script><script>alert(1)</script>`)

	out, err := buildIndexHTML(testFS())
	if err != nil {
		t.Fatalf("buildIndexHTML returned error: %v", err)
	}

	if strings.Contains(string(out), "</script><script>alert(1)") {
		t.Fatalf("raw HTML chars were not escaped, got: %s", out)
	}
	if !strings.Contains(string(out), `\u003c`) {
		t.Fatalf("expected unicode-escaped HTML chars in output, got: %s", out)
	}
}

func TestHandler_ServesInjectedIndexForRoot(t *testing.T) {
	t.Setenv("LOGO_LINK_URL", "https://example.org/")

	h, err := newHandlerWithFS(testFS())
	if err != nil {
		t.Fatalf("newHandlerWithFS returned error: %v", err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("expected text/html content type, got %q", ct)
	}
	body := rr.Body.String()
	want := `window.__APP_CONFIG__={"logoLinkUrl":"https://example.org/"}`
	if !strings.Contains(body, want) {
		t.Fatalf("root response missing injected config %q, got: %s", want, body)
	}
}

func TestHandler_ServesInjectedIndexForSPAFallback(t *testing.T) {
	t.Setenv("LOGO_LINK_URL", "https://example.org/")

	h, err := newHandlerWithFS(testFS())
	if err != nil {
		t.Fatalf("newHandlerWithFS returned error: %v", err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/some/spa/route", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	want := `window.__APP_CONFIG__={"logoLinkUrl":"https://example.org/"}`
	if !strings.Contains(body, want) {
		t.Fatalf("SPA fallback response missing injected config %q, got: %s", want, body)
	}
}

func TestHandler_ServesInjectedIndexForDirectIndexHTML(t *testing.T) {
	t.Setenv("LOGO_LINK_URL", "https://example.org/")

	h, err := newHandlerWithFS(testFS())
	if err != nil {
		t.Fatalf("newHandlerWithFS returned error: %v", err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	want := `window.__APP_CONFIG__={"logoLinkUrl":"https://example.org/"}`
	if !strings.Contains(body, want) {
		t.Fatalf("direct /index.html response missing injected config %q, got: %s", want, body)
	}
}

func TestHandler_ServesStaticFilesDirectly(t *testing.T) {
	t.Setenv("LOGO_LINK_URL", "https://example.org/")

	h, err := newHandlerWithFS(testFS())
	if err != nil {
		t.Fatalf("newHandlerWithFS returned error: %v", err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/assets/index-abc.js", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	if string(body) != "console.log('app')" {
		t.Fatalf("expected static file content, got: %s", body)
	}
	// Static files must NOT contain the injected config.
	if strings.Contains(string(body), "__APP_CONFIG__") {
		t.Fatalf("static file response should not contain injected config")
	}
}

func TestHandler_ServesFavicon(t *testing.T) {
	t.Setenv("LOGO_LINK_URL", "https://example.org/")

	h, err := newHandlerWithFS(testFS())
	if err != nil {
		t.Fatalf("newHandlerWithFS returned error: %v", err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/favicon.svg", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	if string(body) != "<svg></svg>" {
		t.Fatalf("expected favicon content, got: %s", body)
	}
}
