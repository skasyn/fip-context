package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"github.com/tidwall/gjson"
)

type WikiService interface {
	GetArtistByName(name string) (string, error)
	GetGenresByPageTitle(pageTitle string) ([]string, error)
}

type defaultWikiService struct {
	WikiApiURL string
}

// https://en.wikipedia.org/w/api.php?action=opensearch
//
//	&search=test 			# What to look for
//	&limit=1					# Only the first result
//	&namespace=0			# Only articles
//	&format=json
func buildFindNameRequest(name string, wikiUrl string) string {
	nameWithoutWhiteSpace := strings.ReplaceAll(name, " ", "%20") // Needed to search in wikipedia
	path := fmt.Sprintf("%s?action=opensearch&search=%s&limit=1&namespace=0&format=json",
		wikiUrl, nameWithoutWhiteSpace)

	return path
}

func parseArtistPageTitleFromWikiAPI(res io.ReadCloser) (string, error) {
	body, err := io.ReadAll(res)

	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	bodyParsed := gjson.GetBytes(body, "1.0").String() // Get the first element of the second element

	if bodyParsed == "" {
		return "", fmt.Errorf("couldnt get artist page title from Wiki API: %s", gjson.GetBytes(body, "*").String())
	}
	return bodyParsed, nil
}

func (w defaultWikiService) GetArtistByName(name string) (string, error) {
	req := buildFindNameRequest(name, w.WikiApiURL)
	res, err := http.Get(req)

	if err != nil {
		return "", fmt.Errorf("error while fetching %s: %w", req, err)
	} else if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Wiki API returned error status %d: %s", res.StatusCode, req)
	}
	defer res.Body.Close()

	artistTitlePage, err := parseArtistPageTitleFromWikiAPI(res.Body)
	if err != nil {
		return "", fmt.Errorf("couldnt parse artist page title from wiki API response: %w", err)
	}
	return artistTitlePage, nil
}

func (w defaultWikiService) GetGenresByPageTitle(name string) ([]string, error) {
	var res []string
	return res, nil
}

func NewWikiService(wikiApiURL string) WikiService {
	return defaultWikiService{
		WikiApiURL: wikiApiURL,
	}
}
