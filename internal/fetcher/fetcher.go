package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const UserAgent = "Mozilla/5.0 (compatible; KeywordBot/1.0)"

// FetchResult はHTTP取得結果を格納します
// （今後の拡張用に構造体でラップ）
type FetchResult struct {
	URL  string
	Body []byte
}

// FetchURL は指定URLからHTTPレスポンスボディを取得します
func FetchURL(url string, timeoutSeconds int) (*FetchResult, error) {
	client := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create HTTP request for URL '%s': %w", url, err)
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to access URL '%s': %w", url, err)
	}
	defer resp.Body.Close()

	finalURL := resp.Request.URL.String()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body from URL '%s': %w", finalURL, err)
	}

	return &FetchResult{
		URL:  finalURL,
		Body: body,
	}, nil
}
