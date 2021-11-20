package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)


func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	//logger for writting information messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	//logger for writting error messages
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	infoLog.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	errorLog.Fatal(err)
}