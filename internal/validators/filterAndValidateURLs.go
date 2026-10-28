package validators

import (
	"fmt"
	"strings"
)

// Normalize to a list of URLs, accepting either a single "url" or a batch in "urls". Blank entries are dropped and duplicates within the request are collapsed while preserving order.

func FilterAndValidateUrls(urls []string) ([]string, error) {
	seen := map[string]bool{}
	for _, u := range urls {
		u = strings.TrimSpace(u)
		if u == "" || seen[u] {
			continue
		}
		seen[u] = true
		urls = append(urls, u)
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("url is required\n")
	}

	for _, u := range urls {
		if !IsYouTubeURL(u) {
			return nil, fmt.Errorf("Invalid URL: %s\n", u)
		}
	}

	return urls, nil
}
