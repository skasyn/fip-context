package main

import (
	"fmt"
	"log"
)

type FipContextService interface {
	Current(from string) (FipSong, error)
}

type defaultFipContextService struct {
	WikiService WikiService
	FipService FipService
}

func (f *defaultFipContextService) Current(from string) (FipSong, error) {
	fipSong, err := f.FipService.GetCurrentSong(from)

	if err != nil {
		return FipSong{}, fmt.Errorf("couldnt get current song from FipService: %w", err)
	}

	log.Printf("fipSong %v", fipSong)
	artistPageTitle, err := f.WikiService.GetArtistByName(fipSong.Interpreters[0])
	if err != nil {
		return FipSong{}, fmt.Errorf("couldnt get artistPageTitle from WikiService: %w", err)
	}
	log.Printf("artistPageTitle: %s", artistPageTitle)

	return fipSong, nil
}

func NewFipContextService(fipAPiURL string, wikiApiURL string) FipContextService {
	return &defaultFipContextService{
		WikiService: NewWikiService(wikiApiURL),
		FipService: NewFipService(fipAPiURL),
	}
}