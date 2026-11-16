package enricher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/recco-app/recco/internal/domain"
)

var (
	olWorkRegex = regexp.MustCompile(`/works/(OL\w+)`)
	olBookRegex = regexp.MustCompile(`/books/(OL\w+)`)
	olIsbnRegex = regexp.MustCompile(`/isbn/(\w+)`)
)

// OpenLibraryProvider extracts book metadata from openlibrary.org.
type OpenLibraryProvider struct {
	httpClient *http.Client
}

// NewOpenLibraryProvider returns an initialized OpenLibraryProvider.
func NewOpenLibraryProvider(client *http.Client) *OpenLibraryProvider {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return &OpenLibraryProvider{httpClient: client}
}

func (p *OpenLibraryProvider) Name() string {
	return "openlibrary"
}

func (p *OpenLibraryProvider) CanHandle(rawURL string) bool {
	return strings.Contains(strings.ToLower(rawURL), "openlibrary.org")
}

type olWorkResponse struct {
	Title       string      `json:"title"`
	Description interface{} `json:"description"`
	Covers      []int       `json:"covers"`
}

func (p *OpenLibraryProvider) Extract(ctx context.Context, rawURL string) (*domain.EnrichedMetadata, error) {
	var (
		apiURL   string
		sourceID string
	)

	if matches := olWorkRegex.FindStringSubmatch(rawURL); len(matches) > 1 {
		sourceID = fmt.Sprintf("openlibrary:%s", matches[1])
		apiURL = fmt.Sprintf("https://openlibrary.org/works/%s.json", matches[1])
	} else if matches := olBookRegex.FindStringSubmatch(rawURL); len(matches) > 1 {
		sourceID = fmt.Sprintf("openlibrary:%s", matches[1])
		apiURL = fmt.Sprintf("https://openlibrary.org/books/%s.json", matches[1])
	} else if matches := olIsbnRegex.FindStringSubmatch(rawURL); len(matches) > 1 {
		sourceID = fmt.Sprintf("isbn:%s", matches[1])
		apiURL = fmt.Sprintf("https://openlibrary.org/isbn/%s.json", matches[1])
	} else {
		return nil, fmt.Errorf("unrecognized Open Library URL pattern: %s", rawURL)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating Open Library request: %w", err)
	}
	req.Header.Set("User-Agent", "ReccoBot/1.0 (+https://github.com/ShrestHitz/Recco)")

	resp, err := p.httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			resp.Body.Close()
		}
		return &domain.EnrichedMetadata{
			Type:        domain.ItemTypeBook,
			Title:       "Open Library Book",
			Creator:     "Author",
			ExternalURL: rawURL,
			SourceID:    sourceID,
		}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 256*1024))
	if err != nil {
		return nil, fmt.Errorf("failed reading Open Library body: %w", err)
	}

	var data olWorkResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed parsing Open Library json: %w", err)
	}

	coverURL := ""
	if len(data.Covers) > 0 && data.Covers[0] > 0 {
		coverURL = fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-L.jpg", data.Covers[0])
	}

	description := ""
	switch desc := data.Description.(type) {
	case string:
		description = desc
	case map[string]interface{}:
		if val, ok := desc["value"].(string); ok {
			description = val
		}
	}

	return &domain.EnrichedMetadata{
		Type:        domain.ItemTypeBook,
		Title:       data.Title,
		Creator:     "Open Library Author",
		CoverURL:    coverURL,
		ExternalURL: rawURL,
		SourceID:    sourceID,
		Description: description,
	}, nil
}
