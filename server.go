// websockets.go
package main

import (
	"elocalc"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	http.Handle("/home/about/", http.StripPrefix("/home/about/", http.FileServer(http.Dir("home/about"))))
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

	// Serve /calc with a text response.
	http.HandleFunc("/calc", func(w http.ResponseWriter, r *http.Request) {
		// Parses Form
		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
		}
		// Extracts information passed from AJAX statement on examplecalc.html
		p1elo := r.FormValue("P1")
		p1eloint, _ := strconv.Atoi(p1elo)
		p2elo := r.FormValue("P2")
		p2eloint, _ := strconv.Atoi(p2elo)
		winnerElo, loserElo := elocalc.CalcK(1, p1eloint, p2eloint, 0, 0, "Player 1", "Player 2")
		// Display all calc through the console
		println(p1eloint, p2eloint)
		println(winnerElo, loserElo)
		fmt.Fprintf(w, "<h1>Player 1: "+strconv.Itoa(winnerElo)+"</h1>")
		fmt.Fprintf(w, "<h1>Player 2: "+strconv.Itoa(loserElo)+"</h1>")

	})

	// Clears the output
	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
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
