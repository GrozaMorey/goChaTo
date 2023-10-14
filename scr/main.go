package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8000", "http server address")

func statusOk(w http.ResponseWriter, r *http.Request) {
	log.Println("got request", r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Methot not allowed", http.StatusMethodNotAllowed)
		return
	}

}

func main() {
	flag.Parse()

	http.HandleFunc("/", statusOk)

	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
