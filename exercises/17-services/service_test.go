package service

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func do(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestRoutes(t *testing.T) {
	h := NewHandler(MapQuotes{"AAPL": 18950})
	tests := []struct {
		name, method, path, body string
		wantStatus               int
		wantBody                 string // substring
	}{
		{"health", "GET", "/healthz", "", 200, "ok"},
		{"quote", "GET", "/quotes/AAPL", "", 200, `"price":18950`},
		{"unknown quote", "GET", "/quotes/ZZZ", "", 404, `"error"`},
		{"order ok", "POST", "/orders", `{"symbol":"AAPL","side":"buy","qty":10}`, 201, `"id":"ord-`},
		{"order bad side", "POST", "/orders", `{"symbol":"AAPL","side":"hold","qty":10}`, 400, `"error"`},
		{"order unknown field", "POST", "/orders", `{"symbol":"AAPL","side":"buy","qty":1,"x":1}`, 400, `"error"`},
		{"order bad json", "POST", "/orders", `{`, 400, `"error"`},
		{"wrong method", "DELETE", "/orders", "", 405, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := do(t, h, tt.method, tt.path, tt.body)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body %s)", rec.Code, tt.wantStatus, rec.Body)
			}
			if !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("body %q lacks %q", rec.Body, tt.wantBody)
			}
		})
	}
}

func TestOrderBodyLimit(t *testing.T) {
	h := NewHandler(MapQuotes{})
	big := `{"symbol":"` + strings.Repeat("A", 2<<20) + `","side":"buy","qty":1}`
	if rec := do(t, h, "POST", "/orders", big); rec.Code != 400 {
		t.Errorf("2 MiB body status = %d, want 400", rec.Code)
	}
}

func TestOrderIDsUnique(t *testing.T) {
	h := NewHandler(MapQuotes{})
	var mu sync.Mutex
	seen := map[string]bool{}
	var wg sync.WaitGroup
	for range 50 {
		wg.Go(func() {
			rec := do(t, h, "POST", "/orders", `{"symbol":"A","side":"buy","qty":1}`)
			var resp map[string]string
			_ = json.Unmarshal(rec.Body.Bytes(), &resp)
			mu.Lock()
			seen[resp["id"]] = true
			mu.Unlock()
		})
	}
	wg.Wait()
	if len(seen) != 50 {
		t.Errorf("%d unique IDs from 50 orders", len(seen))
	}
}

func TestGracefulShutdown(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	started := make(chan struct{})
	slow := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		close(started)
		time.Sleep(200 * time.Millisecond)
		_, _ = w.Write([]byte("finished"))
	})
	ctx, cancel := context.WithCancel(context.Background())
	served := make(chan error, 1)
	go func() { served <- Serve(ctx, ln, slow) }()

	body := make(chan string, 1)
	go func() {
		resp, err := http.Get("http://" + ln.Addr().String())
		if err != nil {
			body <- "error: " + err.Error()
			return
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		body <- string(b)
	}()

	<-started
	cancel() // shut down while the request is in flight

	select {
	case err := <-served:
		if err != nil {
			t.Fatalf("Serve = %v, want nil", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Serve did not return after cancel")
	}
	if got := <-body; got != "finished" {
		t.Errorf("in-flight request got %q, want it to finish", got)
	}
}
