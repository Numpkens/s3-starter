package main

import (
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

type apiConfig struct {
	db           database.Client
	jwtSecret    string
	platform     string
	filepathRoot string
	assetsRoot   string
	port         string
}