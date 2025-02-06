package main

import (
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

func main() {
	file, err := os.Open("test.pdf")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open file")
	}
	defer file.Close()

	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/binary", file)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create request")
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to send request")
	}

	defer resp.Body.Close()

	log.Debug().Msgf("response: %d", resp.StatusCode)
}
