package main

import (
	"log"
	"net/http"
	"os"
	"oss/apiServer/heartbeat"
	"oss/apiServer/locate"
	"oss/apiServer/objects"

	"./temp"
	"./versions"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
