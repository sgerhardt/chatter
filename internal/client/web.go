package client

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
)

// FromWebsite reads and parses text from a website
func (c *ElevenLabs) FromWebsite(url string) ([]string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch website: %w", err)
	}
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			log.Println("failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch website: %s", resp.Status)
	}

	text, err := extractTextFromHTML(resp.Body)
	if err != nil {
		return nil, err
	}

	return batchText(text, c.Config.CharacterRequestLimit), nil
}

// extractTextFromHTML extracts text from HTML document
func extractTextFromHTML(r io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	var sb strings.Builder
	// Select relevant tags and extract text
	count := 0
	doc.Find("title, h1, h2, h3, h4, h5, h6, p").Each(func(_ int, s *goquery.Selection) {
		count += utf8.RuneCountInString(s.Text())
		sb.WriteString(s.Text())
		sb.WriteString("\n")
	})

	sb.Len()

	return sb.String(), nil
}

// batchText splits the text into chunks of specified size
func batchText(text string, size int) []string {
	var batches []string
	runes := []rune(text)
	for len(runes) > size {
		batch := runes[:size]
		runes = runes[size:]
		batches = append(batches, string(batch))
	}
	batches = append(batches, string(runes))
	return batches
}
