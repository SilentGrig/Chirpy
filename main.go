package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	FileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		FileserverHits: 0,
	}

	router := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerHealth)
	apiRouter.HandleFunc("/reset", apiCfg.handlerReset)
	router.Mount("/api/", apiRouter)

	template := template.Must(template.ParseFiles("metrics.html"))

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerMetrics(template))
	router.Mount("/admin", adminRouter)

	corsMux := middlewareCors(router)

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
