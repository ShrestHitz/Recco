package enricher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/recco-app/recco/internal/domain"
)

var (
	ytWatchRegex  = regexp.MustCompile(`(?:v=|/v/|youtu\.be/|/embed/|/shorts/)([^&?#/]+)`)
)

// YouTubeProvider extracts metadata from YouTube links using official oEmbed.
type YouTubeProvider struct {
	httpClient *http.Client
}

// NewYouTubeProvider returns an initialized YouTubeProvider.
func NewYouTubeProvider(client *http.Client) *YouTubeProvider {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return &YouTubeProvider{httpClient: client}
}

func (p *YouTubeProvider) Name() string {
	return "youtube"
}

func (p *YouTubeProvider) CanHandle(rawURL string) bool {
	lower := strings.ToLower(rawURL)
	return strings.Contains(lower, "youtube.com") || strings.Contains(lower, "youtu.be")
}

type ytOEmbedResponse struct {
	Title        string `json:"title"`
	AuthorName   string `json:"author_name"`
	AuthorURL    string `json:"author_url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

func (p *YouTubeProvider) Extract(ctx context.Context, rawURL string) (*domain.EnrichedMetadata, error) {
	videoID := extractYouTubeID(rawURL)
	if videoID == "" {
		return nil, fmt.Errorf("unable to extract video ID from YouTube URL: %s", rawURL)
	}

	canonicalURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	oembedURL := fmt.Sprintf("https://www.youtube.com/oembed?url=%s&format=json", url.QueryEscape(canonicalURL))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, oembedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating oembed request: %w", err)
	}
	req.Header.Set("User-Agent", "ReccoBot/1.0 (+https://github.com/ShrestHitz/Recco)")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		// Fallback if network fails: still return basic metadata with ID & thumbnail
		return &domain.EnrichedMetadata{
			Type:        domain.ItemTypeYouTube,
			Title:       fmt.Sprintf("YouTube Video (%s)", videoID),
			Creator:     "YouTube Creator",
			CoverURL:    fmt.Sprintf("https://i.ytimg.com/vi/%s/hqdefault.jpg", videoID),
			ExternalURL: canonicalURL,
			SourceID:    fmt.Sprintf("yt:%s", videoID),
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &domain.EnrichedMetadata{
			Type:        domain.ItemTypeYouTube,
			Title:       fmt.Sprintf("YouTube Video (%s)", videoID),
			Creator:     "YouTube Creator",
			CoverURL:    fmt.Sprintf("https://i.ytimg.com/vi/%s/hqdefault.jpg", videoID),
			ExternalURL: canonicalURL,
			SourceID:    fmt.Sprintf("yt:%s", videoID),
		}, nil
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 128*1024))
	if err != nil {
		return nil, fmt.Errorf("failed reading oembed response: %w", err)
	}

	var data ytOEmbedResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed parsing oembed json: %w", err)
	}

	coverURL := data.ThumbnailURL
	if coverURL == "" {
		coverURL = fmt.Sprintf("https://i.ytimg.com/vi/%s/hqdefault.jpg", videoID)
	}

	return &domain.EnrichedMetadata{
		Type:        domain.ItemTypeYouTube,
		Title:       data.Title,
		Creator:     data.AuthorName,
		CoverURL:    coverURL,
		ExternalURL: canonicalURL,
		SourceID:    fmt.Sprintf("yt:%s", videoID),
		Metadata: map[string]interface{}{
			"channel_url": data.AuthorURL,
			"video_id":    videoID,
		},
	}, nil
}

func extractYouTubeID(rawURL string) string {
	matches := ytWatchRegex.FindStringSubmatch(rawURL)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
