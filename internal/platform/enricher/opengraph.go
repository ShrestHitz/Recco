package enricher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
	"github.com/recco-app/recco/internal/domain"
)

// OpenGraphProvider extracts standard OpenGraph and HTML metadata from web pages.
type OpenGraphProvider struct {
	httpClient *http.Client
}

// NewOpenGraphProvider returns an initialized OpenGraphProvider.
func NewOpenGraphProvider(client *http.Client) *OpenGraphProvider {
	if client == nil {
		client = &http.Client{Timeout: 6 * time.Second}
	}
	return &OpenGraphProvider{httpClient: client}
}

func (p *OpenGraphProvider) Name() string {
	return "opengraph"
}

func (p *OpenGraphProvider) CanHandle(rawURL string) bool {
	// Fallback provider capable of handling any HTTP/HTTPS URL
	return strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://")
}

func (p *OpenGraphProvider) Extract(ctx context.Context, rawURL string) (*domain.EnrichedMetadata, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating HTTP request: %w", err)
	}

	// Browser-like user agent to avoid bot blockers on TMDB / IMDb
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed fetching URL for OpenGraph extraction: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("URL returned non-success HTTP status %d", resp.StatusCode)
	}

	// Limit reader to 512KB to protect backend from memory exhaustion
	limitReader := io.LimitReader(resp.Body, 512*1024)

	var (
		ogTitle       string
		ogImage       string
		ogDescription string
		ogSiteName    string
		pageTitle     string
		author        string
	)

	tokenizer := html.NewTokenizer(limitReader)
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}

		tok := tokenizer.Token()
		if tok.Data == "head" && tt == html.EndTagToken {
			break // Stop parsing once head tag closes
		}

		if tok.Data == "title" && tt == html.StartTagToken {
			if tokenizer.Next() == html.TextToken {
				pageTitle = strings.TrimSpace(tokenizer.Token().Data)
			}
		}

		if tok.Data == "meta" && (tt == html.StartTagToken || tt == html.SelfClosingTagToken) {
			var prop, name, content string
			for _, attr := range tok.Attr {
				switch strings.ToLower(attr.Key) {
				case "property":
					prop = strings.ToLower(attr.Val)
				case "name":
					name = strings.ToLower(attr.Val)
				case "content":
					content = attr.Val
				}
			}

			switch prop {
			case "og:title":
				ogTitle = content
			case "og:image":
				ogImage = content
			case "og:description":
				ogDescription = content
			case "og:site_name":
				ogSiteName = content
			}

			switch name {
			case "author":
				author = content
			case "description":
				if ogDescription == "" {
					ogDescription = content
				}
			}
		}
	}

	title := ogTitle
	if title == "" {
		title = pageTitle
	}
	if title == "" {
		title = "Cultural Recommendation"
	}

	creator := author
	if creator == "" {
		creator = ogSiteName
	}

	mediaType := inferMediaTypeFromURL(rawURL)

	return &domain.EnrichedMetadata{
		Type:        mediaType,
		Title:       title,
		Creator:     creator,
		CoverURL:    ogImage,
		ExternalURL: rawURL,
		SourceID:    rawURL,
		Description: ogDescription,
	}, nil
}

func inferMediaTypeFromURL(rawURL string) domain.ItemType {
	lower := strings.ToLower(rawURL)
	switch {
	case strings.Contains(lower, "themoviedb.org"), strings.Contains(lower, "imdb.com"),
		strings.Contains(lower, "letterboxd.com"), strings.Contains(lower, "rottentomatoes.com"):
		return domain.ItemTypeMovie
	case strings.Contains(lower, "anilist.co"), strings.Contains(lower, "myanimelist.net"):
		return domain.ItemTypeAnime
	case strings.Contains(lower, "goodreads.com"), strings.Contains(lower, "openlibrary.org"):
		return domain.ItemTypeBook
	case strings.Contains(lower, "spotify.com/episode"), strings.Contains(lower, "spotify.com/show"),
		strings.Contains(lower, "podcasts.apple.com"):
		return domain.ItemTypePodcast
	case strings.Contains(lower, "youtube.com"), strings.Contains(lower, "youtu.be"):
		return domain.ItemTypeYouTube
	default:
		return domain.ItemTypeMovie
	}
}
