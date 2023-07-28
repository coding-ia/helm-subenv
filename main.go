package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	argsWithProg := os.Args

	s := strings.TrimPrefix(argsWithProg[4], "subenv://")

	content := ""
	uri, err := url.ParseRequestURI(s)

	if err != nil {
		content = parseFile(s)
	}

	_ = uri

	fmt.Print(content)
}

func parseHttpContent(uri string) string {
	resp, err := http.Get(uri)
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
