package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/ric-ram/go-chirpy/internal/database"
)

type apiConfig struct {
	jwtSecret      string
	fileserverHits int
	DB             *database.DB
}

var debugMode = flag.Bool("debug", false, "Enable debug mode")

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	flag.Parse()

	if *debugMode {
		err := os.Remove("database.json")
		if err != nil {
			fmt.Println(err)
		}
	}

	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		jwtSecret:      jwtSecret,
		fileserverHits: 0,
		DB:             db,
	}

	router := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handleReadiness)
	apiRouter.Get("/reset", apiCfg.handlerReset)
	apiRouter.Get("/chirps", apiCfg.handlerChirpsGet)
	apiRouter.Get("/chirps/{id}", apiCfg.handlerChirpsGetById)
	apiRouter.Post("/chirps", apiCfg.handlerChirpsPost)
	apiRouter.Post("/users", apiCfg.handlerUsersPost)
	apiRouter.Post("/login", apiCfg.handlerUserLogin)
	apiRouter.Post("/refresh", apiCfg.handlerTokenRefresh)
	apiRouter.Post("/revoke", apiCfg.handlerTokenRevoke)
	apiRouter.Put("/users", apiCfg.handlerUserUpdate)
	router.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)
	router.Mount("/admin", adminRouter)

	corsMux := middlewareCors(router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
