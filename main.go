package main

import (
	"fmt"
	"github.com/getsops/sops/v3/decrypt"
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

	r, _ := regexp.Compile("^(.*?)://")
	prefix := r.FindString(argsWithProg[4])

	content := ""

	switch prefix {
	case "subenv://":
		content = getContent(argsWithProg[4])
	case "subenv+sops://":
		content = getContentSops(argsWithProg[4])
	}

	parsed := os.ExpandEnv(string(content))

	fmt.Print(parsed)
}

func getContent(path string) string {
	s := strings.TrimPrefix(path, "subenv://")

	if isPathUri(s) {
		return string(downloadContent(s))
	} else {
		return fileContent(s)
	}
}

func getContentSops(path string) string {
	s := strings.TrimPrefix(path, "subenv+sops://")

	if isPathUri(s) {
		content := downloadContent(s)
		data, err := decrypt.Data(content, "yaml")

		if err != nil {
			os.Exit(0)
		}

		return string(data)
	} else {
		data, err := decrypt.File(s, "yaml")

		if err != nil {
			os.Exit(0)
		}

		return string(data)
	}
}

func isPathUri(path string) bool {
	isUri, _ := regexp.MatchString("^https?://", path)
	return isUri
}

func downloadContent(path string) []byte {
	expanded := os.ExpandEnv(path)
	uri, err := url.ParseRequestURI(expanded)

	if err != nil {
		os.Exit(0)
	}

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

	return content
}

func fileContent(path string) string {
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

	return string(content)
}
