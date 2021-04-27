package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

func getBody(addr string) string {
	resp, err := http.Get(addr)
	if err != nil {
		fmt.Println(err)

	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()
	// fmt.Println(body)
	return body
}

func search(query string, massive []string) []string {
	var result []string
	// fmt.Println(query)
	// fmt.Println(massive)
	for i := 0; i < len(massive); i++ {
		if strings.Contains(getBody(massive[i]), query) {

			result = append(result, massive[i])
		}
	}
	return result
}

func main() {
	var where []string
	what := "Go is an open source programming language"
	where = append(where, "http://golang.org")

	fmt.Println(search(what, where))

}
