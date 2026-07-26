package ai

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestExternalIntelligenceSearchRSSFeed(t *testing.T) {
	rssBody := `<?xml version="1.0"?>
<rss version="2.0">
  <channel>
    <item>
      <title>Indonesia pharma market trend grows</title>
      <link>https://example.test/pharma-market</link>
      <description>Demand for healthcare products increases in urban hospitals.</description>
      <pubDate>Sun, 26 Jul 2026 10:00:00 GMT</pubDate>
    </item>
    <item>
      <title>Unrelated shipping update</title>
      <link>https://example.test/shipping</link>
      <description>Logistics news.</description>
      <pubDate>Sat, 25 Jul 2026 10:00:00 GMT</pubDate>
    </item>
  </channel>
</rss>`

	service := newExternalIntelligenceService(externalIntelligenceConfig{
		Enabled:    true,
		FeedURLs:   []string{"https://allowed.example/feed.xml"},
		Timeout:    time.Second,
		CacheTTL:   time.Minute,
		MaxResults: 2,
	})
	service.client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "allowed.example" {
			t.Fatalf("unexpected host %s", req.URL.Host)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(rssBody)),
			Request:    req,
		}, nil
	})

	result := service.Search(context.Background(), "tren pasar pharma healthcare", time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC))

	if !result.Enabled {
		t.Fatal("expected external intelligence to be enabled")
	}
	if len(result.Sources) == 0 {
		t.Fatalf("expected at least one external source, got notice %q", result.Notice)
	}
	if !strings.Contains(strings.ToLower(result.Sources[0].Title), "pharma") {
		t.Fatalf("expected pharma source to rank first, got %#v", result.Sources)
	}
}

func TestExternalIntelligenceRejectsDisallowedHost(t *testing.T) {
	service := newExternalIntelligenceService(externalIntelligenceConfig{
		Enabled:        true,
		FeedURLs:       []string{"https://not-allowed.example/feed.xml"},
		AllowedDomains: []string{"allowed.example"},
		Timeout:        time.Second,
		CacheTTL:       time.Minute,
		MaxResults:     1,
	})

	result := service.Search(context.Background(), "market trend", time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC))

	if len(result.Sources) != 0 {
		t.Fatalf("expected no sources from disallowed host, got %#v", result.Sources)
	}
	if !strings.Contains(result.Notice, "allowlist") {
		t.Fatalf("expected allowlist notice, got %q", result.Notice)
	}
}

func TestExternalIntelligenceDisabledNotice(t *testing.T) {
	service := newExternalIntelligenceService(externalIntelligenceConfig{Enabled: false})

	result := service.Search(context.Background(), "market trend", time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC))

	if result.Enabled {
		t.Fatal("expected external intelligence to be disabled")
	}
	if !strings.Contains(result.Notice, "AI_EXTERNAL_INTELLIGENCE_ENABLED") {
		t.Fatalf("expected configuration notice, got %q", result.Notice)
	}
}
