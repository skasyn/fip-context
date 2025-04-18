package main

import (
	"fmt"
)

type FipContextService interface {
	Current(from string) (FipSong, error)
}

type defaultFipContextService struct {
	WikiService WikiService
	FipService  FipService
}

func (f *defaultFipContextService) Current(from string) (FipSong, error) {
	fipSong, err := f.FipService.GetCurrentSong(from)

	if err != nil {
		return FipSong{}, fmt.Errorf("couldnt get current song from FipService: %w", err)
	}

	artistPageTitles, err := f.WikiService.GetArtistPageTitlesByNames(fipSong.Interpreters)
	if err != nil {
		return FipSong{}, fmt.Errorf("couldnt get artistPageTitle from WikiService: %w", err)
	}

	genres, err := f.WikiService.GetGenresFromArtists(artistPageTitles)
	if err != nil {
		return FipSong{}, fmt.Errorf("couldnt get genres from %v: %w", artistPageTitles, err)
	}
	fipSong.Genres = genres
	return fipSong, nil
}

func NewFipContextService(fipAPiURL string, wikiApiURL string, dbpediaURL string) (FipContextService, error) {
	wikiService, err := NewWikiService(wikiApiURL, dbpediaURL)

	if err != nil {
		return nil, fmt.Errorf("failed to create wikiService with wikiApiURL %s and dbpediaURL %s: %w", wikiApiURL, dbpediaURL, err)
	}
	return &defaultFipContextService{
		WikiService: wikiService,
		FipService:  NewFipService(fipAPiURL),
	}, nil
}
