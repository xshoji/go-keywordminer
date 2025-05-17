package fetcher

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchURL_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>Hello</body></html>"))
	}))
	defer ts.Close()

	res, err := FetchURL(ts.URL, 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.URL != ts.URL {
		t.Errorf("expected URL %s, got %s", ts.URL, res.URL)
	}
	if string(res.Body) != "<html><body>Hello</body></html>" {
		t.Errorf("unexpected body: %s", string(res.Body))
	}
}

func TestFetchURL_Timeout(t *testing.T) {
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// no response, simulate timeout
	}))
	// Listen but do not start serving, so connection will hang
	ts.Listener.Close() // force connection error

	_, err := FetchURL(ts.URL, 1) // 1秒タイムアウト
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}
