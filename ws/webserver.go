package ws

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

// Starts a web server to expose health-check endpoints
// the assumption is that as long as the server works, the sync works too
// We may add more checks to make sure the sync executed correctly
func RunWS(httpAddr string) {
	log.Println("Starting HTTP on port", httpAddr)
	http.HandleFunc("/health-check", handler)
	http.ListenAndServe(httpAddr, nil)
}
