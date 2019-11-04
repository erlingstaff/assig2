package main

import (
	"log"
	"net/http"
	"time"
)

var seconds int //global value to store starttime of the server in UNIX time

func commitHandler(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	lmt := ""
	vrf := ""
	keys, ok := r.URL.Query()["limit"]
	if ok {
		lmt = keys[0]
	}
	keys, ok = r.URL.Query()["auth"]
	if ok {
		vrf = keys[0]
	}
	getCommits(w, lmt, vrf) //parameters = url path, ResponseWriter
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	diagnosTics(w, seconds) //parameters = ResponseWriter and UNIX time of server start
}

func languagesHandler(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json") //makes content-type application/json
	auth := ""
	lmt := ""
	keys, ok := r.URL.Query()["auth"]
	if ok {
		auth = keys[0]
	}
	keys, ok = r.URL.Query()["limit"]
	if ok {
		lmt = keys[0]
	}
	getLanguages(w, auth, lmt, r) //parameters = url path, ResponseWriter and limit query key

}

func main() {
	port := "5632"
	serverstart := int(time.Now().Unix()) //logging unix time of server start as
	//a global variable, used as parameter above
	seconds = serverstart
	http.HandleFunc("/repocheck/v1/commits/", commitHandler) //3 different handlers
	http.HandleFunc("/repocheck/v1/issues/", languagesHandler)
	http.HandleFunc("/repocheck/v1/status/", statusHandler)
	http.HandleFunc("/repocheck/v1/webhooks/", webhookHandler)
	//http.HandleFunc("/repocheck/v1/service/", serviceHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil)) //server on port 5632
}
