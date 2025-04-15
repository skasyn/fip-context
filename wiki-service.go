package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/knakk/sparql"
	"github.com/tidwall/gjson"
)

type WikiService interface {
	GetArtistPageTitleByName(name string) (string, error)
	getGenresByArtistPageTitle(name string) ([]string, error)
	GetGenresFromArtists(artists []string) ([]string, error)
}

type defaultWikiService struct {
	WikiApiURL  string
	DbpediaRepo *sparql.Repo
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

// http://en.wikipedia.org/w/api.php?action=query&prop=revisions&rvprop=content&format=xmlfm&titles=Scary%20Monsters%20and%20Nice%20Sprites&rvsection=0
func (w defaultWikiService) GetArtistPageTitleByName(name string) (string, error) {
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

// sanitize name injected in SPARQL query
func safeSPARQLRessourceName(name string) string {
	uriCompliantName := strings.ReplaceAll(name, " ", "_")

	// handle accented letters
	reg, _ := regexp.Compile(`[^A-Za-zÀ-ÖØ-öø-ÿ0-9_(),.-]`)
	sanitized := reg.ReplaceAllString(uriCompliantName, "")

	// Ensure there are no SPARQL-specific characters that could change query meaning
	sanitized = strings.ReplaceAll(sanitized, "'", "")
	sanitized = strings.ReplaceAll(sanitized, "\"", "")
	sanitized = strings.ReplaceAll(sanitized, "\\", "")
	sanitized = strings.ReplaceAll(sanitized, "`", "")

	return sanitized
}

func (w defaultWikiService) getGenresByArtistPageTitle(name string) ([]string, error) {
	safeArtistName := safeSPARQLRessourceName(name)
	sparqlQuery := fmt.Sprintf(`
	SELECT ?property ?value WHERE {
		dbr:%s ?property ?value .
		FILTER (?property = dbp:genre)
	}
	LIMIT 10
`, safeArtistName)

	res, err := w.DbpediaRepo.Query(sparqlQuery)
	if err != nil {
		return []string{}, fmt.Errorf("failed to query dbpedia with %s: %w", safeArtistName, err)
	}
	fmt.Printf("name: %s | res: %v\n", safeArtistName, res.Solutions())

	genres := make([]string, len(res.Solutions()))
	for i, binding := range res.Solutions() {
		genreURI := binding["value"].String()
		genreSplited := strings.Split(genreURI, `/`)
		if len(genreSplited) != 0 {
			genres[i] = genreSplited[len(genreSplited)-1]
		} else {
			log.Printf("Couldnt split following URI: %s", genreURI)
		}
	}
	return slices.DeleteFunc(genres, func(g string) bool {
		return g == ""
	}), nil
}

// Look for the genres of music played by each artist of a given list
// Genres will only appear once
// artists[] need to be a list of artist page title on wikipedia (ex: The Beatles -> The_Beatles)
func (w defaultWikiService) GetGenresFromArtists(artists []string) ([]string, error) {
	genresMap := make(map[string]bool)
	var mutex sync.Mutex
	var wg sync.WaitGroup

	errorChan := make(chan error, len(artists))

	wg.Add(len(artists))
	for _, artist := range artists {
		go func(artist string) {
			defer wg.Done()

			genres, err := w.getGenresByArtistPageTitle(artist)

			if err != nil {
				errorChan <- fmt.Errorf("error getting genres for %s: %w", artist, err)
				return
			}

			// only keep one occurence of each genres
			mutex.Lock()
			for _, s := range genres {
				genresMap[s] = true
			}
			mutex.Unlock()
		}(artist)
	}

	wg.Wait()
	close(errorChan)

	genres := make([]string, 0, len(genresMap))
	for g := range genresMap {
		genres = append(genres, g)
	}

	var errorsArray []error
	for err := range errorChan {
		errorsArray = append(errorsArray, err)
	}
	if len(errorsArray) > 0 {
		return genres, errors.Join(errorsArray...)
	}
	return genres, nil
}

func NewWikiService(wikiApiURL string, dbpediaSPARQLEndpoint string) (WikiService, error) {
	sparqlRepo, err := sparql.NewRepo(dbpediaSPARQLEndpoint,
		sparql.Timeout(time.Second*30),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create SPARQL repository: %w", err)
	}
	return defaultWikiService{
		WikiApiURL:  wikiApiURL,
		DbpediaRepo: sparqlRepo,
	}, nil
}
