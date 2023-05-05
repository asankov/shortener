package shortener

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/asankov/shortener/pkg/links"
	"github.com/gorilla/mux"
)

func (s *Shortener) routes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/{id}", s.handleGetLink).Methods(http.MethodGet)

	// TODO: auth
	apiRoutes := router.PathPrefix("/api/v1").Subrouter()
	apiRoutes.HandleFunc("/{id}", s.handleCreateLink).Methods(http.MethodPost)
	apiRoutes.HandleFunc("/{id}", s.handleDeleteLink).Methods(http.MethodDelete)
	// TODO
	// apiRoutes.HandleFunc("/{id}", s.handleGetLinkMetrics).Methods(http.MethodGet)

	return router
}

func (s *Shortener) handleGetLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link, err := s.db.GetByID(id)
	if err != nil {
		if errors.Is(err, links.ErrLinkNotFound) {
			w.WriteHeader(http.StatusNotFound)
		}

		s.logger.Warnf("unknown error while getting link by id [%s]: %v", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: increment the number of clicks

	http.Redirect(w, r, link.URL, http.StatusFound)
}

type linkRequest struct {
	URL string `json:"url"`
}

func (s *Shortener) handleCreateLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var link linkRequest
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		s.logger.Warnf("failed to decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.db.Create(id, link.URL); err != nil {
		s.logger.Warnf("failed to create link: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Shortener) handleDeleteLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.db.Delete(id); err != nil {
		s.logger.Warnf("failed to delete link: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
