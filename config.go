package main

import "github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"

type apiConfig struct {
	db               database.Client
	jwtSecret        string
	platform         string
	filepathRoot     string
	assetsRoot       string
	s3Bucket         string
	s3Region         string
	s3CfDistribution string
	port             string
}

// Global videoThumbnails map and thumbnail struct removed