package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// what := "Go is an open source programming language"
// where = append(where, "http://golang.org")

type Query struct {
	What  string `json:"what" xml:"what"`
	Where string `json:"where" xml:"where"`
}

func getBody(addr string) string {
	resp, err := http.Get(addr)
	if err != nil {
		fmt.Println(err)

	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	expBody := string(body)
	return expBody
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/", firstHandler)
	router.HandleFunc("/search", searchHandler)
	router.HandleFunc("/setcookie", setCookieHandler)
	router.HandleFunc("/takecookie", takeCookieHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func firstHandler(wr http.ResponseWriter, req *http.Request) {
	wr.Write([]byte("There is a search at /search by json or xml: \"what\" and \"where\" fields"))
}

func searchHandler(wr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		wr.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	contentTypeHeader := req.Header.Get("Content-Type")

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		wr.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	request := &Query{}

	switch contentTypeHeader {
	case "application/xml":
		if err = xml.Unmarshal(data, request); err != nil {
			log.Println(err)
			wr.WriteHeader(http.StatusInternalServerError)
			return
		}
	default: //json
		if err = json.Unmarshal(data, request); err != nil {
			log.Println(err)
			wr.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
	log.Printf("What: %s\nWhere: %s\n", request.What, request.Where)

	if strings.Contains(getBody(request.Where), request.What) {
		obj, err := json.Marshal(request.Where)
		if err != nil {
			log.Println(err)
		}
		wr.Header().Set("Content-Type", "application/json")
		wr.Write(obj)
	}
	wr.WriteHeader(http.StatusOK)
}

func setCookieHandler(wr http.ResponseWriter, req *http.Request) {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	http.SetCookie(wr, &http.Cookie{
		Name:       "cookie",
		Value:      "setting value",
		Path:       "/",
		Domain:     "",
		Expires:    expiration,
		RawExpires: "",
		MaxAge:     0,
		Secure:     false,
		HttpOnly:   true,
		SameSite:   0,
		Raw:        "",
		Unparsed:   []string{},
	})
	fmt.Println(wr, "Put Cookie")
}

func takeCookieHandler(wr http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("cookie")
	if err != nil {
		http.Error(wr, http.StatusText(400), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(wr, "Name: ", c.Name, "\n", "Value: ", c.Value, "\n", "Expires: ", c.Expires)
}
