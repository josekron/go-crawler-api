package main

import (
	"database/sql"
	"encoding/json"

	"log"

	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(dbname string) {
	var err error
	a.DB, err = sql.Open("postgres", dbname)
	if err != nil {
		log.Fatal(err)
	}

	// create the table if it doesn't exist
	_, err = a.DB.Exec("CREATE TABLE IF NOT EXISTS jobs (id SERIAL PRIMARY KEY, title TEXT, company TEXT, url TEXT, web TEXT)")

	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/jobs", a.getJobs).Methods("GET")
	a.Router.HandleFunc("/crawl", a.crawlURL).Methods("POST")

}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getJobs(w http.ResponseWriter, r *http.Request) {
	users, err := getJobs(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (a *App) crawlURL(w http.ResponseWriter, r *http.Request) {
	webURL := struct {
		URL string `json:"url"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&webURL)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var jobs []Job = crawlURL(webURL.URL)
	var jobsCreated []Job
	for _, job := range jobs {
		//TODO: check job URL doesn't exist in the DB
		createJob(a.DB, &job)
		jobsCreated = append(jobsCreated, job)
	}

	respondWithJSON(w, http.StatusOK, jobsCreated)
}
