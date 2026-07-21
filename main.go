package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type PageData struct {
	Title string
	Time  string
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := PageData{
			Title: "Lucas Oliveira",
			Time:  time.Now().Format(time.RFC822),
		}

		tmpl.Execute(w, data)
	})
	port := ":8080"
	fmt.Println("sv up @ port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
