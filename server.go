// websockets.go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/bmizerany/pq"
)

func main() {

	var navProjectsTemplate = "navProjectsTemplate.html"
	var footerTemplate = "footerTemplate.html"
	var headerTemplate = "headerTemplate.html"
	var navTemplate = "navTemplate.html"

	// Handler points to available directories
	http.Handle("/home/", http.StripPrefix("/home/", http.FileServer(http.Dir("home"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("scripts"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, readFile(headerTemplate))
		if strings.Contains(r.URL.Path[1:], "projects") {
			fmt.Fprint(w, readFile(navProjectsTemplate))
		} else {
			fmt.Fprint(w, readFile(navTemplate))
		}
		if r.URL.Path[1:] == "" {
			http.ServeFile(w, r, "home/index.html")
		} else {
			http.ServeFile(w, r, "home/"+r.URL.Path[1:]+".html")
		}
		fmt.Fprint(w, readFile(footerTemplate))
	})

	//Serves local webpage for testing
	if true {
		errhttp := http.ListenAndServe(":8080", nil)
		if errhttp != nil {
			log.Fatal("Web server (HTTPS): ", errhttp)
		}
	} else {
		//Serves the webpage
		errhttps := http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
		if errhttps != nil {
			log.Fatal("Web server (HTTPS): ", errhttps)
		}
	}
}

func readFile(fileName string) string {
	data, err := ioutil.ReadFile("partials/" + fileName)
	if err != nil {
		return "Error parsing file."
	}
	return string(data)
}
