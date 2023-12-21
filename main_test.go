package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	// TODO: change the database url to a test database url
	a.Initialize("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec("CREATE TABLE IF NOT EXISTS jobs (id SERIAL PRIMARY KEY, title TEXT, company TEXT, url TEXT, web TEXT)"); err != nil {
		log.Fatal(err)
	}
}
func clearTable() {
	a.DB.Exec("DELETE FROM jobs")
	a.DB.Exec("ALTER TABLE jobs AUTO_INCREMENT = 1")
}

func TestEmptyJobs(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/jobs", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCrawlJobs(t *testing.T) {
	clearTable()
	// TODO: change the body url to a fixed url that won't change
	payload := []byte(`{"url":"https://boards.greenhouse.io/monzo"}`)
	req, _ := http.NewRequest("POST", "/crawl", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m []Job
	json.Unmarshal(response.Body.Bytes(), &m)
	if len(m) == 0 {
		t.Errorf("Expected a non empty array. Got %+v", m)
	}
	if m[0].Company != "monzo" {
		t.Errorf("Expected Company monzo. Got %s", m[0].Company)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
