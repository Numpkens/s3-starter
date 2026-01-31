package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerThumbnailGet(w http.ResponseWriter, r *http.Request) {
	// thumbnails are now stored as data URLs in the video's thumbnail_url column.
	// This endpoint is deprecated.
	respondWithError(w, http.StatusNotFound, "Thumbnail endpoint removed; thumbnails are embedded in video metadata", nil)
}
