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

	pathToDB := os.Getenv("DB_PATH")
	if pathToDB == "" {
		log.Fatal("DB_PATH must be set")
	}

	db, err := database.NewClient(pathToDB)
	if err != nil {
		log.Fatalf("Couldn't connect to database: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM environment variable is not set")
	}

	filepathRoot := os.Getenv("FILEPATH_ROOT")
	if filepathRoot == "" {
		log.Fatal("FILEPATH_ROOT environment variable is not set")
	}

	assetsRoot := os.Getenv("ASSETS_ROOT")
	if assetsRoot == "" {
		log.Fatal("ASSETS_ROOT environment variable is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	cfg := apiConfig{
		db:           db,
		jwtSecret:    jwtSecret,
		platform:     platform,
		filepathRoot: filepathRoot,
		assetsRoot:   assetsRoot,
		port:         port,
	}

	err = cfg.ensureAssetsDir()
	if err != nil {
		log.Fatalf("Couldn't create assets directory: %v", err)
	}

	mux := http.NewServeMux()
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", appHandler)

	assetsHandler := http.StripPrefix("/assets", http.FileServer(http.Dir(assetsRoot)))
	mux.Handle("/assets/", cacheMiddleware(assetsHandler))

	// Register API routes (use correct patterns in handler funcs)
	mux.HandleFunc("/api/login", cfg.handlerLogin)
	mux.HandleFunc("/api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("/api/revoke", cfg.handlerRevoke)

	mux.HandleFunc("/api/users", cfg.handlerUsersCreate)

	mux.HandleFunc("/api/videos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			cfg.handlerVideoMetaCreate(w, r)
		case http.MethodGet:
			cfg.handlerVideosRetrieve(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/api/videos/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cfg.handlerVideoGet(w, r)
		case http.MethodDelete:
			cfg.handlerVideoMetaDelete(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/api/thumbnails/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			cfg.handlerUploadThumbnail(w, r)
		case http.MethodGet:
			cfg.handlerThumbnailGet(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/api/video_upload/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			cfg.handlerUploadVideo(w, r)
			return
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("/admin/reset", cfg.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on: http://localhost:%s/app/\n", port)
	log.Fatal(srv.ListenAndServe())
}
