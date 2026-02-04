package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")

	db, err := database.NewClient(os.Getenv("DB_PATH"))
	if err != nil {
		log.Fatalf("Couldn't connect to database: %v", err)
	}

	cfg := apiConfig{
		db:               db,
		jwtSecret:        os.Getenv("JWT_SECRET"),
		platform:         os.Getenv("PLATFORM"),
		filepathRoot:     os.Getenv("FILEPATH_ROOT"),
		assetsRoot:       os.Getenv("ASSETS_ROOT"),
		s3Bucket:         os.Getenv("S3_BUCKET"),
		s3Region:         os.Getenv("S3_REGION"),
		s3CfDistribution: os.Getenv("S3_CF_DISTRO"),
		port:             os.Getenv("PORT"),
	}

	mux := http.NewServeMux()

	// Static file servers
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(cfg.filepathRoot))))
    
	// Assets handler updated to use the noCacheMiddleware
	assetsHandler := http.StripPrefix("/assets", http.FileServer(http.Dir(cfg.assetsRoot)))
	mux.Handle("/assets/", noCacheMiddleware(assetsHandler))

	// Routes
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/videos", cfg.handlerVideoMetaCreate)
	mux.HandleFunc("GET /api/videos", cfg.handlerVideosRetrieve)
	mux.HandleFunc("GET /api/videos/{videoID}", cfg.handlerVideoGet)
	mux.HandleFunc("POST /api/thumbnail_upload/{videoID}", cfg.handlerUploadThumbnail)

	log.Printf("Serving on: http://localhost:%s/app/\n", cfg.port)
	log.Fatal(http.ListenAndServe(":"+cfg.port, mux))
}