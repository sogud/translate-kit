package translatekit

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestGoogle429StopsRepeatedRequestsAndCanCancelPacing(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1); w.WriteHeader(429) }))
	defer server.Close()
	s := New(Options{})
	s.googleURL = server.URL
	for i := 0; i < 2; i++ {
		result, err := s.Translate(context.Background(), Request{Config: Config{Provider: "google", TargetLanguage: "zh-CN"}, Text: "# Hello\n\nText.", Revision: "r"})
		if err == nil || result.Text != "" {
			t.Fatal("Rate limit accepted", result, err)
		}
	}
	if calls.Load() != 1 {
		t.Fatal("Cooldown sent repeated requests", calls.Load())
	}
	s.webMu <- struct{}{}
	s.webBlocked = nil
	s.webNext = time.Now().Add(time.Hour)
	<-s.webMu
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := s.googleRequest(ctx, server.URL)
	if !errors.Is(err, context.Canceled) || calls.Load() != 1 {
		t.Fatal("Canceled pacing sent request", err)
	}
	s.webMu <- struct{}{}
	s.webNext = time.Time{}
	s.webBlocked = map[string]time.Time{"google": time.Now().Add(time.Minute)}
	<-s.webMu
	_, _ = s.webRequest(context.Background(), "youdao", http.MethodGet, server.URL, "", nil, nil)
	if calls.Load() != 2 {
		t.Fatal("Google cooldown blocked another provider")
	}
}

func TestLostBatchBoundariesAreRejected(t *testing.T) {
	output, err := translateWebBatches([]string{"Heading", "Text"}, 3000, func(string) (string, error) { return "merged translation", nil }, func(int) {})
	if err == nil || output != nil {
		t.Fatal("Lost line boundaries accepted", output, err)
	}
}
