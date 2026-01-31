package main

import (
    "encoding/base64"
    "encoding/json"
    "net/http"
    "os"
    "path/filepath"
    "fmt"
    "strings"

    "github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
    "github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
    "github.com/google/uuid"
)

func (cfg *apiConfig) handlerVideoMetaCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		database.CreateVideoParams
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	params.UserID = userID

	video, err := cfg.db.CreateVideo(params.CreateVideoParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create video", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, video)
}

func (cfg *apiConfig) handlerVideoMetaDelete(w http.ResponseWriter, r *http.Request) {
    // parse videoID from URL path
    videoIDString := strings.TrimPrefix(r.URL.Path, "/api/videos/")
    if videoIDString == "" {
        respondWithError(w, http.StatusBadRequest, "Missing video ID in path", nil)
        return
    }
    videoID, err := uuid.Parse(videoIDString)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
        return
    }

    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
        return
    }
    _, err = auth.ValidateJWT(token, cfg.jwtSecret)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
        return
    }

    err = cfg.db.DeleteVideo(videoID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not delete video", err)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerVideoGet(w http.ResponseWriter, r *http.Request) {
    // parse videoID from URL path
    videoIDString := strings.TrimPrefix(r.URL.Path, "/api/videos/")
    if videoIDString == "" {
        respondWithError(w, http.StatusBadRequest, "Missing video ID in path", nil)
        return
    }
    videoID, err := uuid.Parse(videoIDString)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
        return
    }

    video, err := cfg.db.GetVideo(videoID)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "Failed to get video", err)
        return
    }

    respondWithJSON(w, http.StatusOK, video)
}

func (cfg *apiConfig) handlerVideosRetrieve(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	_, err = auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	videos, err := cfg.db.GetVideos()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get videos", err)
		return
	}

	// If a video has no thumbnail in the DB, try to load the default asset,
	// encode it as a data URL, save it back to the DB and include it in the response.
	for i := range videos {
		if videos[i].ThumbnailURL == nil {
			assetPath := filepath.Join(cfg.assetsRoot, "boots-image-horizontal.png")
			if _, err := os.Stat(assetPath); err == nil {
				if b, err := os.ReadFile(assetPath); err == nil {
					enc := base64.StdEncoding.EncodeToString(b)
					dataURL := fmt.Sprintf("data:image/png;base64,%s", enc)
					videos[i].ThumbnailURL = &dataURL
					_ = cfg.db.UpdateVideo(videos[i]) // best-effort update; ignore errors here
				}
			}
		}
	}

	respondWithJSON(w, http.StatusOK, videos)
}
