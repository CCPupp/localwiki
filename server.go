// websockets.go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/bmizerany/pq"
	"github.com/kabukky/httpscerts"
)

func main() {

	var navProjectsTemplate = "navProjectsTemplate.html"
	var footerTemplate = "footerTemplate.html"
	var headerTemplate = "headerTemplate.html"
	var navTemplate = "navTemplate.html"

	// Check if the cert files are available.
	err := httpscerts.Check("cert.pem", "key.pem")
	// If they are not available, generate new ones.
	if err != nil {
		err = httpscerts.Generate("cert.pem", "key.pem", "127.0.0.1:8080")
		if err != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}
	}
	// Handler points to available directories
	http.Handle("/home/", http.StripPrefix("/home/", http.FileServer(http.Dir("home"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("scripts"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, readFile(headerTemplate))
		if strings.Contains(r.URL.Path[1:], "projects") {
			fmt.Fprintf(w, readFile(navProjectsTemplate))
		} else {
			fmt.Fprintf(w, readFile(navTemplate))
		}
		if r.URL.Path[1:] == "" {
			http.ServeFile(w, r, "home/index.html")
		} else {
			http.ServeFile(w, r, "home/"+r.URL.Path[1:]+".html")
		}
		fmt.Fprintf(w, readFile(footerTemplate))
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
