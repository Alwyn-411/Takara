// Using ForExEndpoint v2 endpoint for conversion rates
package forex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type response struct {
	Base  string  `json:"base"`
	Quote string  `json:"quote"`
	Rate  float64 `json:"rate"`
	Date  string  `json:"date"`
}

type cache struct {
	rate      string
	fetchedAt time.Time
}

type ForEx struct {
	client   *http.Client
	cache    map[string]cache
	mu       sync.RWMutex
	endpoint string
	ttl      time.Duration
}

func NewAccessor() *ForEx {
	endpoint := os.Getenv("FOREX_ENDPOINT")

	return &ForEx{
		client:   &http.Client{Timeout: 10 * time.Second},
		cache:    make(map[string]cache),
		endpoint: endpoint,
		ttl:      1 * time.Hour, // rates don't change intraday
	}
}

// GetRate fetches the current rate for a currency pair
func (s *ForEx) GetRate(base, quote string) (string, error) {
	if base == quote {
		return "1", nil
	}

	cacheKey := base + "#" + quote
	s.mu.RLock()
	if cached, ok := s.cache[cacheKey]; ok && time.Since(cached.fetchedAt) < s.ttl {
		s.mu.RUnlock()
		return cached.rate, nil
	}
	s.mu.RUnlock()

	url := fmt.Sprintf("https://%s/v2/rate/%s/%s", s.endpoint, base, quote)
	resp, err := s.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("ForEx request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("ForEx returned %d", resp.StatusCode)
	}

	var data response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("decode failed: %w", err)
	}

	rateStr := fmt.Sprintf("%f", data.Rate)

	s.mu.Lock()
	s.cache[cacheKey] = cache{rate: rateStr, fetchedAt: time.Now()}
	s.mu.Unlock()

	return rateStr, nil
}

// GetHistoricalRate fetches the rate on a specific date
func (s *ForEx) GetHistoricalRate(base, quote, date string) (string, error) {
	url := fmt.Sprintf(
		"https://%s/v2/rate/%s/%s?date=%s",
		s.endpoint, base, quote, date,
	)

	resp, err := s.client.Get(url)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("ForEx returned %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var data response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	return fmt.Sprintf("%f", data.Rate), nil
}
