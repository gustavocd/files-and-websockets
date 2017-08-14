package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gustavocd/files-and-websockets/handlers"
	"github.com/sirupsen/logrus"
)

func main() {

	r := mux.NewRouter().StrictSlash(false)
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file, _ := ioutil.ReadFile("./public/index.html")
		w.Write(file)
	})

	r.HandleFunc("/ws/upload", handlers.Upload)

	go handlers.HandleFile()

	logrus.Println("Running localhost:8081")
	http.ListenAndServe(":8081", r)
}
