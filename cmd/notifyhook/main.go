package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		if msg, err := ioutil.ReadAll(r.Body); err != nil {
			log.Printf("can get notification body: %s", err.Error())
		} else {
			log.Printf("got notification: %s", string(msg))
		}
	})
	err := http.ListenAndServe(":8080", http.DefaultServeMux)
	panic(err)
}
