package ws

import (
	"log"
	"net/http"
	"encoding/json"
//	"time"
	"github.com/adobe-apiplatform/api-gateway-config-supervisor/sync"
)

var status = sync.GetStatusInstance()

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	status.Status = "OK"
	json.NewEncoder(w).Encode(status)
}

// Starts a web server to expose health-check endpoints
// the assumption is that as long as the server works, the sync works too
// We may add more checks to make sure the sync executed correctly
func RunWS(httpAddr string) {
	log.Println("Starting HTTP on port", httpAddr)
	http.HandleFunc("/health-check", healthCheckHandler)
	http.ListenAndServe(httpAddr, nil)
}
