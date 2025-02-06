package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func ChecksumUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		checksum := r.Header.Get("X-Checksum")

		if checksum == "" {
			log.Error().Msg("checksum is required")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hash := md5.New()
		if _, err := io.Copy(hash, r.Body); err != nil {
			log.Error().Err(err).Msg("failed to copy checksum from header")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		file, err := os.OpenFile(filepath.Join("tmp", uuid.NewString()), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Error().Err(err).Msg("Failed to open file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		tee := io.TeeReader(r.Body, hash)
		if _, err := io.Copy(file, tee); err != nil {
			log.Error().Err(err).Msg("failed to copy body to response writer")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		computedChecksum := hex.EncodeToString(hash.Sum(nil))
		if computedChecksum != checksum {
			log.Error().Msg("checksum mismatch")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
