package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

type FipSong struct {
	Author []string
	Name   string
}

type FipService interface {
	GetCurrentSong(from string) (FipSong, error)
}

type DefaultFipService struct{}

func buildRequest(from string, fipUrl string) string {
	path := fipUrl + "/live"

	if from != "" {
		path = fmt.Sprintf("%s?webradio=%s", path, from)
	}
	return path
}

func buildSongFromFipAPI(res io.ReadCloser) (FipSong, error) {
	body, err := io.ReadAll(res)

	if err != nil {
		return FipSong{}, fmt.Errorf("failed to read response body: %w", err)
	}
	bodyParsed := gjson.ParseBytes(body)

	authorsParsed := bodyParsed.Get("now.song.interpreters").Array()

	if len(authorsParsed) == 0 {
		return FipSong{}, fmt.Errorf("couldnt parse the interpreters of current song: %v", bodyParsed.Raw)
	}

	authors := make([]string, len(authorsParsed))
	for i, authorParsed := range authorsParsed {
		authors[i] = authorParsed.String()
	}
	return FipSong{
		Author: authors,
		Name:   bodyParsed.Get("now.song.release.title").String(),
	}, nil

}

func (f *DefaultFipService) GetCurrentSong(from string) (FipSong, error) {

	fipAPI, ok := os.LookupEnv("FIP_API")
	if !ok {
		return FipSong{}, errors.New("FIP_API is not set in env")
	} else if fipAPI == "" {
		return FipSong{}, errors.New("FIP_API is empty")
	}

	req := buildRequest(from, fipAPI)
	res, err := http.Get(req)

	if err != nil {
		return FipSong{}, fmt.Errorf("error while fetching %s: %w", req, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return FipSong{}, fmt.Errorf("FIP API return error status %d: %s", res.StatusCode, req)
	}

	currentSong, err := buildSongFromFipAPI(res.Body)
	if err != nil {
		return FipSong{}, err // errors are already explicit
	}

	return currentSong, nil
}

// https://www.radiofrance.fr/fip/api/live?webradio=fip_pop
// https://www.radiofrance.fr/fip/api?webradio=fip_pop
