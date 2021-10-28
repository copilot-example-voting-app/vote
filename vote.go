// Package vote starts the vote server.
package vote

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/copilot-example-voting-app/vote/server"
)

// Run starts the server.
func Run() error {
	addr := flag.String("addr", ":8080", "port to listen on")
	flag.Parse()
	handler, err := server.NewServer()
	if err != nil {
		return err
	}

	s := http.Server{
		Addr:         *addr,
		Handler:      handler,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	log.Printf("INFO: vote: listen on port %s\n", *addr)
	return s.ListenAndServe()
}
