package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	go http.ListenAndServe(":8081", &Handler{})
	fmt.Println("HTTP server started on port 8081")

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
	fmt.Println(time.Now().Format(time.RFC822), "[proxy] received request to", r.URL)

	// define origin server URL
	serverURL, err := url.Parse("http://127.0.0.1:8080")
	if err != nil {
		log.Fatal("invalid origin server URL")
	}
	// set req Host, URL and Request URI to forward a request to the origin server
	r.Host = serverURL.Host
	r.URL.Host = serverURL.Host
	r.URL.Scheme = serverURL.Scheme
	r.RequestURI = ""

	// save the response from the origin server
	serverResponse, err := http.DefaultClient.Do(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err)
		return
	}

	// copy headers
	for k, values := range serverResponse.Header {
		for _, v := range values {
			w.Header().Add(k, v)
		}
	}

	// return response to the client
	io.Copy(w, serverResponse.Body)

}
