package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type FipSong struct {
	Author string
	Name   string
}

type FipService interface {
	GetCurrentSong(from string) (FipSong, error)
}

type DefaultFipService struct{}

func buildRequest(from string, fipUrl string) string {
	path := fipUrl + "/live"

	if from != "" {
		path = fmt.Sprintf("%s?webradio=%s", fipUrl, from)
	}
	return path
}

func (f *DefaultFipService) GetCurrentSong(from string) (FipSong, error) {
	var currentSong FipSong

	fipAPI, ok := os.LookupEnv("FIP_API")
	if !ok {
		return currentSong, errors.New("FIP_API is not set in env")
	} else if fipAPI == "" {
		return currentSong, errors.New("FIP_API is empty")
	}

	req := buildRequest(from, fipAPI)
	resp, err := http.Get(req)

	if err != nil {
		return currentSong, fmt.Errorf("error while fetching %s: %w", req, err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)

	return currentSong, nil
}
