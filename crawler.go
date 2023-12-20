package main

import (
	"strings"

	"github.com/gocolly/colly"
)

// Different Portals will implements the interface because
// they will need different approaches to scrap the information
type JobPortal interface {
	crawlURL(url string) []Job
}

type GreenHousePortal struct{}

func (p GreenHousePortal) crawlURL(url string) []Job {
	var jobs []Job
	var urlSplit []string = strings.Split(url, "/")
	var company string = urlSplit[len(urlSplit)-1]

	c := colly.NewCollector()

	c.OnHTML("div.opening", func(e *colly.HTMLElement) {
		job := Job{}

		job.Title = strings.TrimSpace(strings.ReplaceAll(e.Text, "\n", ""))
		job.URL = e.ChildAttr("a", "href")
		job.Web = url
		job.Company = company

		jobs = append(jobs, job)
	})

	c.Visit(url)

	return jobs
}

func crawlURL(url string) []Job {
	var jobs []Job

	if strings.Contains(url, "boards.greenhouse.io") {
		var jobPortal = GreenHousePortal{}
		jobs = jobPortal.crawlURL(url)
	}

	return jobs
}
