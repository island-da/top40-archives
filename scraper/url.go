package scraper

import (
	"log"
	"net/url"
	"path"
	"strconv"
	"strings"
)

func ParseYear(targetUrl string) int {
	parsedURL, err := url.Parse(targetUrl)
	if err != nil {
		log.Fatalf("Failed to parse URL %s: %v", targetUrl, err)
	}

	var parsedYear int
	segments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(segments) >= 2 {
		parsedYear, _ = strconv.Atoi(segments[len(segments)-2])
	} else {
		log.Println("Not enough segments in path:", parsedURL.Path)
	}
	return parsedYear
}

func ParseDateArchives(targetUrl string) int {
	parsedURL, err := url.Parse(targetUrl)
	if err != nil {
		log.Fatalf("Failed to parse URL %s: %v", targetUrl, err)
	}

	var parsedDate int
	segments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(segments) >= 2 {
		dateWithExt := segments[len(segments)-1]
		parsedDate, _ = strconv.Atoi(strings.TrimSuffix(dateWithExt, path.Ext(dateWithExt))[:2])
	} else {
		log.Println("Not enough segments in path:", parsedURL.Path)
	}
	return parsedDate
}
