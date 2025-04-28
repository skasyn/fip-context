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
	psqlConnStr string
	Env string
}

func LoadConfig() (*Config, error) {
	fipAPIURL, ok := os.LookupEnv("FIP_API")
	if !ok {
		return nil, errors.New("FIP_API is not set in env")
	} else if fipAPIURL == "" {
		return nil, errors.New("FIP_API is empty")
	}

	wikiApiURL, ok := os.LookupEnv("WIKI_API")
	if !ok {
		return nil, errors.New("WIKI_API is not set in env")
	} else if wikiApiURL == "" {
		return nil, errors.New("WIKI_API is empty")
	}

	dbpediaURl, ok := os.LookupEnv("DBPEDIA_SPARQL")
	if !ok {
		return nil, errors.New("DBPEDIA_SPARQL is not set in env")
	} else if dbpediaURl == "" {
		return nil, errors.New("DBPEDIA_SPARQL is empty")
	}

	psqlConnStr, ok := os.LookupEnv("PSQL_CONN_STR")
	if !ok {
		return nil, errors.New("PSQL_CONN_STR is not set in env")
	} else if psqlConnStr == "" {
		return nil, errors.New("PSQL_CONN_STR is empty")
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
		psqlConnStr: psqlConnStr,
		Env: env,
	}, nil
}