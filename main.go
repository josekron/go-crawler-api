package main

import (
	"os"
)

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("DATABASE_URL"))
	// a.Initialize(
	// 	"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	a.Run(":8000")
}
