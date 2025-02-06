package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func BinaryUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20) //10MB
		defer r.Body.Close()

		file, err := os.OpenFile(filepath.Join("tmp", uuid.NewString()), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Error().Err(err).Msg("Failed to open file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		copiedFile, err := io.Copy(file, r.Body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to copy file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info().
			Int("written size", int(copiedFile)).
			Msgf("File uploaded: %s", file.Name())

		w.WriteHeader(http.StatusOK)
	}
}
