# Go Crawler API

API to web scraping job portals

Right now only supports greenhouse.io

### Technologies

- Golang + Gorilla Mux
- PostgreSQL
- Docker

### Example

`POST localhost:8000/crawl`

```
{
    "url": "https://boards.greenhouse.io/monzo"
}
```

`GET localhost:8000/jobs`

### How to run

- Docker-Compose: `docker-compose up --build`
- PostgreSQL on local: change `DATABASE_URL` in `main.go`

### Debug with VS-Code

- Add `Launch Package` to create a launch.json
- If you have issues with DAP server and DLV:
    - Ctrl + Shift + P
    - Type `Go: Install/Update Tools`