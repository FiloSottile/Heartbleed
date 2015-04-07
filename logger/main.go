package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

func main() {
	f, err := os.OpenFile(os.Args[1], os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	lock := new(sync.Mutex)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		lock.Lock()
		defer lock.Unlock()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
			return
		}
		_, err = f.Write(body)
		if err != nil {
			log.Print(err)
			return
		}
		_, err = f.WriteString("\n")
		if err != nil {
			log.Print(err)
			return
		}
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		f.Close()
		os.Exit(0)
	}()

	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
