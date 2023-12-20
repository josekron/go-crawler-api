package main

import (
	"database/sql"
	"log"
)

type Job struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Company string `json:"company"`
	URL     string `json:"url"`
	Web     string `json:"web"`
}

func getJobs(db *sql.DB) ([]Job, error) {
	rows, err := db.Query(
		"SELECT * FROM jobs")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []Job{}

	for rows.Next() {
		var j Job
		if err := rows.Scan(&j.ID, &j.Title, &j.Company, &j.URL, &j.Web); err != nil {
			return nil, err
		}
		users = append(users, j)
	}

	return users, nil
}

func getJobByID(db *sql.DB, id int) (Job, error) {
	var job Job

	err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&job.ID, &job.Title, &job.Company, &job.URL, &job.Web)
	if err != nil {
		log.Fatal(err)
	}

	return job, nil
}

func createJob(db *sql.DB, job *Job) {

	err := db.QueryRow("INSERT INTO jobs (title, company, url, web) VALUES ($1, $2, $3, $4) RETURNING id", job.Title, job.Company, job.URL, job.Web).Scan(&job.ID)
	if err != nil {
		log.Fatal(err)
	}
}
