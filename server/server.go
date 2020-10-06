// Package server provides the HTTP server to render the "vote" web page.
package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Server is the Vote server.
type Server struct {
	Router *mux.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.HandleFunc("/", s.handleView()).Methods(http.MethodGet)
	s.Router.HandleFunc("/", s.handleSave()).Methods(http.MethodPost)

	s.Router.ServeHTTP(w, r)
}

func (s *Server) handleView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		voterID, err := getVoterID(r)
		if err != nil {
			http.Error(w, "get voter id", http.StatusInternalServerError)
			return
		}
		vote, err := getVote(voterID)
		if err != nil {
			http.Error(w, "get vote", http.StatusInternalServerError)
			return
		}
		renderTemplate(w, "index", vote)
	}
}

func (s *Server) handleSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		voterID, err := getVoterID(r)
		if err != nil {
			http.Error(w, "get voter id", http.StatusInternalServerError)
			return
		}
		vote := r.FormValue("vote")
		if err := saveVote(voterID, vote); err != nil {
			http.Error(w, "save vote", http.StatusInternalServerError)
			return
		}
		renderTemplate(w, "index", vote)
	}
}

// getVoterID reads the voter_id cookie.
// If the cookie doesn't exist then generated a random UUID as the voter ID.
func getVoterID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("voter_id")
	switch err {
	case nil:
		return cookie.Value, nil
	case http.ErrNoCookie:
		rand, err := uuid.NewRandom()
		if err != nil {
			return "", fmt.Errorf("generate random UUID for the voter: %v", err)
		}
		return rand.String(), nil
	default:
		log.Printf("ERROR: server: get voter ID: %v\n", err)
		return "", fmt.Errorf(`get "voter_id" cookie: %v`, err)
	}
}

// getVote queries the API microservice to retrieve the vote.
// If the vote doesn't exist, then returns an empty string "".
func getVote(voterID string) (string, error) {
	endpoint := fmt.Sprintf("http://api.%s:8080/votes/%s", os.Getenv("COPILOT_SERVICE_DISCOVERY_ENDPOINT"), voterID)
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Printf("WARN: server: coudln't get vote for voter id %s: %v\n", voterID, err)
		return "", nil
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("WARN: server: get vote response status: %d\n", resp.StatusCode)
		return "", nil
	}


	defer resp.Body.Close()
	data := struct {
		Result string `json:"vote"`
	}{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		log.Printf("ERROR: server: decode vote data: %v\n", err)
		return "", fmt.Errorf("server: decode vote: %v",err)
	}
	log.Printf("INFO: server: received vote %s for voter id %s\n", data.Result, voterID)
	return data.Result, nil
}

func saveVote(voterID, vote string) error {
	dat, err := json.Marshal(&struct{
		VoterID string `json:"voter_id"`
		Vote string `json:"vote"`
	}{
		VoterID: voterID,
		Vote:    vote,
	})
	if err != nil {
		log.Printf("ERROR: server encode save vote data: %v\n", err)
		return fmt.Errorf("server: encode save vote data: %v", err)
	}

	endpoint := fmt.Sprintf("http://api.%s:8080/votes", os.Getenv("COPILOT_SERVICE_DISCOVERY_ENDPOINT"))
	_, err = http.Post(endpoint, "application/json", bytes.NewBuffer(dat))
	if err != nil {
		log.Printf("ERROR: server: save vote %s for voter id %s: %v\n", vote, voterID, err)
		return fmt.Errorf("server: save vote: %v", err)
	}
	log.Printf("INFO: server: saved vote %s for voter id %s\n", vote, voterID)
	return nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, vote string) {
	t, err := template.ParseFiles(filepath.Join("templates", tmpl + ".html"))
	if err != nil {
		log.Fatalf("parse file: %v\n", err)
	}
	t.Execute(w, struct {
		Vote string
	} {
		Vote: vote,
	})
}