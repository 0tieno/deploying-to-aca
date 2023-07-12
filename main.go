package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var name string = os.Getenv("VISITOR_NAME")

func main() {
	fmt.Println(name)
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", homeHandler)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/html/home.html")
	if err != nil {
		log.Fatal(err)
	}

	if name == "" {
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	data := struct {
		Name string
	}{
		Name: name,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}
