package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	go http.ListenAndServe(":8080", &Handler{})
	fmt.Println("HTTP server started on port 8080")

	WaitForCtrlC()
}

func WaitForCtrlC() {
	var end_waiter sync.WaitGroup
	end_waiter.Add(1)
	var signal_channel chan os.Signal
	signal_channel = make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)
	go func() {
		<-signal_channel
		end_waiter.Done()
	}()
	end_waiter.Wait()

	fmt.Println("MAIN", "...bye bye")
}

type Handler struct {
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().Format(time.RFC822), "[origin server] received request to", r.URL)

	w.Header().Add("unheader", "ciao")
	w.Header().Add("unheader", "ciao2")
	fmt.Fprint(w, "origin server response")
}
