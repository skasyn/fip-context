package main

import (
	"fmt"

	"github.com/skasyn/fip-context/fipcontextrepo"
	"github.com/skasyn/fip-context/fipservice"
	"github.com/skasyn/fip-context/wikiservice"
)

type FipContextService interface {
	Current(from string) (fipservice.FipSong, error)
}

type defaultFipContextService struct {
	WikiService wikiservice.WikiService
	FipService  fipservice.FipService
	Repo        *fipcontextrepo.FipSongRepository
}

func (f *defaultFipContextService) Current(from string) (fipservice.FipSong, error) {
	fipSong, err := f.FipService.GetCurrentSong(from)

	if err != nil {
		return fipservice.FipSong{}, fmt.Errorf("couldnt get current song from FipService: %w", err)
	}

	artistPageTitles, err := f.WikiService.GetArtistPageTitlesByNames(fipSong.Interpreters)
	if err != nil {
		return fipservice.FipSong{}, fmt.Errorf("couldnt get artistPageTitle from WikiService: %w", err)
	}

	genres, err := f.WikiService.GetGenresFromArtists(artistPageTitles)
	if err != nil {
		return fipservice.FipSong{}, fmt.Errorf("couldnt get genres from %v: %w", artistPageTitles, err)
	}
	fipSong.Genres = genres
	fmt.Printf("new fipsong: %v\n", fipSong)
	fipSongModel := f.Repo.FipSongToFipSongModel(&fipSong)
	err = f.Repo.Create(fipSongModel)

	if err != nil {
		return fipservice.FipSong{}, fmt.Errorf("couldnt create fipsong %v in database:  %w", fipSong, err)
	}
	return fipSong, nil
}

func NewFipContextService(cfg *Config, r *fipcontextrepo.FipSongRepository) (FipContextService, error) {
	wikiService, err := wikiservice.NewWikiService(cfg.WikiApiURL, cfg.DbpediaURL)

	if err != nil {
		return nil, fmt.Errorf("failed to create wikiService with wikiApiURL %s and dbpediaURL %s: %w", cfg.WikiApiURL, cfg.FIPApiURL, err)
	}
	return &defaultFipContextService{
		WikiService: wikiService,
		FipService:  fipservice.NewFipService(cfg.FIPApiURL),
		Repo:        r,
	}, nil
}
