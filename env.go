package main

import (
	"errors"
	"log"
	"os"
)

type Config struct {
	FIPApiURL string
	WikiApiURL string
	DbpediaURL string
	Env string
}

func LoadConfig() (*Config, error) {
	fipAPIURL, ok := os.LookupEnv("FIP_API")
	if !ok {
		return &Config{}, errors.New("FIP_API is not set in env")
	} else if fipAPIURL == "" {
		return &Config{}, errors.New("FIP_API is empty")
	}

	wikiApiURL, ok := os.LookupEnv("WIKI_API")
	if !ok {
		return &Config{}, errors.New("WIKI_API is not set in env")
	} else if wikiApiURL == "" {
		return &Config{}, errors.New("WIKI_API is empty")
	}

	dbpediaURl, ok := os.LookupEnv("DBPEDIA_SPARQL")
	if !ok {
		return &Config{}, errors.New("DBPEDIA_SPARQL is not set in env")
	} else if dbpediaURl == "" {
		return &Config{}, errors.New("DBPEDIA_SPARQL is empty")
	}

	env := os.Getenv("FIP_CONTEXT_ENV")

	if env == "" {
		env = "development"
		log.Println("warning: FIP_CONTEXT_ENV is not set, \"development\" will be used")
	}
	
	return &Config{
		FIPApiURL: fipAPIURL,
		WikiApiURL: wikiApiURL,
		DbpediaURL: dbpediaURl,
		Env: env,
	}, nil
}