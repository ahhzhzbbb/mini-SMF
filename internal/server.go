package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	index         = 0
	serverAmmount = 10
)

type reverseProxy struct {
	serverRegistry []*url.URL
	serverAmmount  int
}

func newReverseProxy(ammount int) *reverseProxy {
	return &reverseProxy{
		serverRegistry: make([]*url.URL, ammount),
		serverAmmount:  ammount,
	}
}

func (re reverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hung gay\n"))
	serverUrl := re.serverRegistry[index]
	index = (index + 1) % serverAmmount
	fmt.Fprintf(w, "your request is handling by server number %d, at address %s", index, serverUrl)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("Request: %s %s, Duration: %v\n", r.Method, r.URL.Path, duration)
	})
}

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func main() {
	proxy := newReverseProxy(10)

	for i := range serverAmmount {
		serverUrl, err := url.Parse(fmt.Sprintf("http://127.0.0.1:800%d", i))
		if err != nil {
			log.Fatal(err)
		}

		proxy.serverRegistry[i] = serverUrl
	}

	greetHandler := Chain(proxy, LoggingMiddleware)

	mux := http.NewServeMux()

	mux.Handle("/send", greetHandler)
	mux.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello gay\n"))
	})
	mux.HandleFunc("/greet/{user}", func(w http.ResponseWriter, r *http.Request) {
		username := r.PathValue("user")
		fmt.Fprintf(w, "hello %s gay", username)
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
