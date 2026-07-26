package ai

import (
	"context"
	"encoding/xml"
	"errors"
	"html"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type externalIntelligenceConfig struct {
	Enabled        bool
	FeedURLs       []string
	AllowedDomains []string
	Timeout        time.Duration
	CacheTTL       time.Duration
	MaxResults     int
}

type externalIntelligenceService struct {
	config externalIntelligenceConfig
	client *http.Client
	mu     sync.Mutex
	cache  map[string]externalIntelligenceCacheEntry
}

type externalIntelligenceCacheEntry struct {
	expiresAt time.Time
	result    externalIntelligenceResult
}

type externalIntelligenceResult struct {
	Enabled     bool
	Query       string
	GeneratedAt time.Time
	Notice      string
	Sources     []externalSource
}

type externalSource struct {
	Title       string
	URL         string
	SourceHost  string
	PublishedAt string
	Snippet     string
	Score       int
}

type rssFeed struct {
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type atomFeed struct {
	Entries []atomEntry `xml:"entry"`
}

type atomEntry struct {
	Title   string     `xml:"title"`
	Summary string     `xml:"summary"`
	Content string     `xml:"content"`
	Updated string     `xml:"updated"`
	Links   []atomLink `xml:"link"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

func newExternalIntelligenceServiceFromEnv() *externalIntelligenceService {
	return newExternalIntelligenceService(externalIntelligenceConfig{
		Enabled:        getBoolEnv("AI_EXTERNAL_INTELLIGENCE_ENABLED", false),
		FeedURLs:       splitCSVEnv("AI_EXTERNAL_FEED_URLS"),
		AllowedDomains: splitCSVEnv("AI_EXTERNAL_ALLOWED_DOMAINS"),
		Timeout:        time.Duration(getIntEnv("AI_EXTERNAL_TIMEOUT_SECONDS", 5)) * time.Second,
		CacheTTL:       time.Duration(getIntEnv("AI_EXTERNAL_CACHE_TTL_MINUTES", 60)) * time.Minute,
		MaxResults:     getIntEnv("AI_EXTERNAL_MAX_RESULTS", 5),
	})
}

func newExternalIntelligenceService(config externalIntelligenceConfig) *externalIntelligenceService {
	if config.Timeout <= 0 {
		config.Timeout = 5 * time.Second
	}
	if config.CacheTTL <= 0 {
		config.CacheTTL = 60 * time.Minute
	}
	if config.MaxResults <= 0 {
		config.MaxResults = 5
	}
	if len(config.AllowedDomains) == 0 {
		config.AllowedDomains = hostsFromURLs(config.FeedURLs)
	}
	return &externalIntelligenceService{
		config: config,
		client: &http.Client{Timeout: config.Timeout},
		cache:  make(map[string]externalIntelligenceCacheEntry),
	}
}

func (s *externalIntelligenceService) Search(ctx context.Context, query string, now time.Time) externalIntelligenceResult {
	result := externalIntelligenceResult{Enabled: s != nil && s.config.Enabled, Query: query, GeneratedAt: now}
	if s == nil || !s.config.Enabled {
		result.Notice = "External intelligence belum aktif. Aktifkan AI_EXTERNAL_INTELLIGENCE_ENABLED dan konfigurasi AI_EXTERNAL_FEED_URLS untuk memakai sumber luar database."
		return result
	}
	if len(s.config.FeedURLs) == 0 {
		result.Notice = "External intelligence aktif, tetapi AI_EXTERNAL_FEED_URLS belum dikonfigurasi."
		return result
	}

	cacheKey := normalizeExternalQuery(query) + "|" + strings.Join(s.config.FeedURLs, ",")
	if cached, ok := s.getCached(cacheKey, now); ok {
		return cached
	}

	terms := externalQueryTerms(query)
	var sources []externalSource
	var notices []string
	for _, feedURL := range s.config.FeedURLs {
		items, err := s.fetchFeed(ctx, feedURL)
		if err != nil {
			notices = append(notices, err.Error())
			continue
		}
		for _, source := range items {
			source.Score = scoreExternalSource(source, terms)
			if source.Score > 0 || len(terms) == 0 {
				sources = append(sources, source)
			}
		}
	}

	sort.SliceStable(sources, func(i, j int) bool {
		if sources[i].Score == sources[j].Score {
			return sources[i].PublishedAt > sources[j].PublishedAt
		}
		return sources[i].Score > sources[j].Score
	})
	if len(sources) > s.config.MaxResults {
		sources = sources[:s.config.MaxResults]
	}

	result.Sources = sources
	if len(notices) > 0 {
		result.Notice = strings.Join(notices, "; ")
	}
	if len(sources) == 0 && result.Notice == "" {
		result.Notice = "Tidak ada sumber eksternal yang relevan dari feed yang dikonfigurasi."
	}
	s.setCached(cacheKey, result, now)
	return result
}

func (s *externalIntelligenceService) fetchFeed(ctx context.Context, feedURL string) ([]externalSource, error) {
	parsed, err := url.Parse(feedURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, errors.New("external feed URL tidak valid")
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return nil, errors.New("external feed URL harus memakai http atau https")
	}
	if !hostAllowed(parsed.Host, s.config.AllowedDomains) {
		return nil, errors.New("external feed host tidak ada di allowlist")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, errors.New("external feed request tidak valid")
	}
	req.Header.Set("Accept", "application/rss+xml, application/atom+xml, application/xml, text/xml")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, errors.New("external feed tidak dapat diakses")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("external feed mengembalikan status non-2xx")
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, errors.New("external feed tidak dapat dibaca")
	}
	if sources := parseRSSSources(body, parsed.Host); len(sources) > 0 {
		return sources, nil
	}
	return parseAtomSources(body, parsed.Host), nil
}

func parseRSSSources(body []byte, host string) []externalSource {
	var feed rssFeed
	if err := xml.Unmarshal(body, &feed); err != nil || len(feed.Channel.Items) == 0 {
		return nil
	}
	sources := make([]externalSource, 0, len(feed.Channel.Items))
	for _, item := range feed.Channel.Items {
		sources = append(sources, externalSource{
			Title:       cleanExternalText(item.Title, 140),
			URL:         strings.TrimSpace(item.Link),
			SourceHost:  hostnameOnly(host),
			PublishedAt: strings.TrimSpace(item.PubDate),
			Snippet:     cleanExternalText(item.Description, 260),
		})
	}
	return sources
}

func parseAtomSources(body []byte, host string) []externalSource {
	var feed atomFeed
	if err := xml.Unmarshal(body, &feed); err != nil || len(feed.Entries) == 0 {
		return nil
	}
	sources := make([]externalSource, 0, len(feed.Entries))
	for _, entry := range feed.Entries {
		link := ""
		for _, candidate := range entry.Links {
			if candidate.Rel == "" || candidate.Rel == "alternate" {
				link = candidate.Href
				break
			}
		}
		summary := entry.Summary
		if strings.TrimSpace(summary) == "" {
			summary = entry.Content
		}
		sources = append(sources, externalSource{
			Title:       cleanExternalText(entry.Title, 140),
			URL:         strings.TrimSpace(link),
			SourceHost:  hostnameOnly(host),
			PublishedAt: strings.TrimSpace(entry.Updated),
			Snippet:     cleanExternalText(summary, 260),
		})
	}
	return sources
}

func (s *externalIntelligenceService) getCached(key string, now time.Time) (externalIntelligenceResult, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cached, ok := s.cache[key]
	if !ok || now.After(cached.expiresAt) {
		return externalIntelligenceResult{}, false
	}
	return cached.result, true
}

func (s *externalIntelligenceService) setCached(key string, result externalIntelligenceResult, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[key] = externalIntelligenceCacheEntry{expiresAt: now.Add(s.config.CacheTTL), result: result}
}

func scoreExternalSource(source externalSource, terms []string) int {
	haystack := strings.ToLower(source.Title + " " + source.Snippet + " " + source.SourceHost)
	score := 0
	for _, term := range terms {
		if strings.Contains(haystack, term) {
			score += 2
		}
	}
	if strings.Contains(strings.ToLower(source.Title), "health") || strings.Contains(strings.ToLower(source.Title), "pharma") {
		score++
	}
	return score
}

func externalQueryTerms(query string) []string {
	normalized := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(strings.ToLower(query), " ")
	stopWords := map[string]bool{
		"dan": true, "yang": true, "saya": true, "anda": true, "dari": true, "untuk": true, "buatkan": true,
		"grafik": true, "chart": true, "data": true, "the": true, "and": true, "for": true, "with": true,
	}
	seen := map[string]bool{}
	var terms []string
	for _, term := range strings.Fields(normalized) {
		if len(term) < 3 || stopWords[term] || seen[term] {
			continue
		}
		seen[term] = true
		terms = append(terms, term)
	}
	return terms
}

func normalizeExternalQuery(query string) string {
	return strings.Join(externalQueryTerms(query), " ")
}

func cleanExternalText(value string, limit int) string {
	text := html.UnescapeString(value)
	text = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(text, " ")
	text = strings.Join(strings.Fields(text), " ")
	if limit > 0 && len(text) > limit {
		return strings.TrimSpace(text[:limit]) + "..."
	}
	return text
}

func hostAllowed(host string, allowedDomains []string) bool {
	host = strings.ToLower(hostnameOnly(host))
	for _, allowed := range allowedDomains {
		allowed = strings.ToLower(hostnameOnly(strings.TrimSpace(allowed)))
		if allowed == "" {
			continue
		}
		if host == allowed || strings.HasSuffix(host, "."+allowed) {
			return true
		}
	}
	return false
}

func hostnameOnly(host string) string {
	if parsedHost, _, err := net.SplitHostPort(host); err == nil {
		return parsedHost
	}
	return strings.TrimSpace(host)
}

func hostsFromURLs(rawURLs []string) []string {
	seen := map[string]bool{}
	var hosts []string
	for _, rawURL := range rawURLs {
		parsed, err := url.Parse(strings.TrimSpace(rawURL))
		if err != nil || parsed.Host == "" {
			continue
		}
		host := hostnameOnly(parsed.Host)
		if host == "" || seen[host] {
			continue
		}
		seen[host] = true
		hosts = append(hosts, host)
	}
	return hosts
}

func splitCSVEnv(key string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}
	return items
}

func getBoolEnv(key string, fallback bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if value == "" {
		return fallback
	}
	return value == "true" || value == "1" || value == "yes" || value == "on"
}

func getIntEnv(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
