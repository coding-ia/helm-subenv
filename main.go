package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func main() {
	argsWithProg := os.Args

	if len(argsWithProg) != 5 {
		os.Exit(0)
	}

	s := strings.TrimPrefix(argsWithProg[4], "subenv://")
	content := ""
	isUri, _ := regexp.MatchString("^https?://", s)

	if isUri {
		expanded := os.ExpandEnv(s)
		uri, err := url.ParseRequestURI(expanded)

		if err != nil {
			os.Exit(0)
		}

		content = parseHttpContent(uri)
	} else {
		content = parseFile(s)
	}

	fmt.Print(content)
}

func parseHttpContent(uri *url.URL) string {
	resp, err := http.Get(uri.String())
	if err != nil {
		os.Exit(0)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.Exit(0)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		os.Exit(0)
	}

	parsed := os.ExpandEnv(string(content))
	return parsed
}

func parseFile(path string) string {
	if _, err := os.Stat(path); err != nil {
		os.Exit(0)
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	content, err := io.ReadAll(file)
	parsed := os.ExpandEnv(string(content))

	return parsed
}
