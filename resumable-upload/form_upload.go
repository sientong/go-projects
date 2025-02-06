package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func FormUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(5 << 20); err != nil {
			log.Error().Err(err).Msg("failed to parse multipart form")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.MultipartForm.RemoveAll()

		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Error().Err(err).Msg("failed to get file from form")
		}
		defer file.Close()

		target, err := os.OpenFile(filepath.Join("tmp", handler.Filename), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Error().Err(err).Msg("failed to create file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		copiedFile, err := io.Copy(target, file)
		if err != nil {
			log.Error().Err(err).Msg("failed to copy file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info().
			Int("written size", int(copiedFile)).
			Msgf("File uploaded: %s", target.Name())

		w.WriteHeader(http.StatusOK)
	}
}
