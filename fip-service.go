package main

import (
	"fmt"
	"io"
	"net/http"
	"github.com/tidwall/gjson"
)

type FipSong struct {
	Interpreters []string
	Name   string
}

type FipService interface {
	GetCurrentSong(from string) (FipSong, error)
}

type defaultFipService struct{
	FIPApiURL string
}

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

	interpretersParsed := bodyParsed.Get("now.song.interpreters").Array()

	if len(interpretersParsed) == 0 {
		return FipSong{}, fmt.Errorf("couldnt parse the interpreters of current song: %v", bodyParsed.Raw)
	}

	interpreters := make([]string, len(interpretersParsed))
	for i, authorParsed := range interpretersParsed {
		interpreters[i] = authorParsed.String()
	}
	return FipSong{
		Interpreters: interpreters,
		Name:   bodyParsed.Get("now.song.release.title").String(),
	}, nil

}

func (f *defaultFipService) GetCurrentSong(from string) (FipSong, error) {
	req := buildRequest(from, f.FIPApiURL)
	res, err := http.Get(req)

	if err != nil {
		return FipSong{}, fmt.Errorf("error while fetching %s: %w", req, err)
	} else if res.StatusCode != http.StatusOK {
		return FipSong{}, fmt.Errorf("FIP API return error status %d: %s", res.StatusCode, req)
	}
	defer res.Body.Close()

	currentSong, err := buildSongFromFipAPI(res.Body)
	if err != nil {
		return FipSong{}, err // errors are already explicit
	}

	return currentSong, nil
}

func NewFipService(fipApiURL string) FipService {
	return &defaultFipService{
		FIPApiURL: fipApiURL,
	}
}

// https://www.radiofrance.fr/fip/api/live?webradio=fip_pop
// https://www.radiofrance.fr/fip/api?webradio=fip_pop
