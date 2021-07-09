package main

import (
	"log"
	"net/http"
	"os"
	"oss/dataServer/heartbeat"
	"oss/dataServer/locate"
	"oss/dataServer/objects"
	"oss/dataServer/temp"
)

func main() {
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
