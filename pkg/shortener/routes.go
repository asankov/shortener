package shortener

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"

	"github.com/asankov/shortener/pkg/links"
	"github.com/gorilla/mux"
)

func (s *Shortener) routes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/admin", s.handleAdmin).Methods(http.MethodGet)
	router.HandleFunc("/{id}", s.handleGetLink).Methods(http.MethodGet)

	// TODO: auth
	apiRoutes := router.PathPrefix("/api/v1").Subrouter()
	apiRoutes.HandleFunc("/links", s.handleCreateLink).Methods(http.MethodPost)
	apiRoutes.HandleFunc("/links/{id}", s.handleDeleteLink).Methods(http.MethodDelete)

	// TODO
	// apiRoutes.HandleFunc("/{id}", s.handleGetLinkMetrics).Methods(http.MethodGet)

	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fs))

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

type createLinkRequest struct {
	ID  string `json:"id,omitempty"`
	URL string `json:"url,omitempty"`
}

func (s *Shortener) handleCreateLink(w http.ResponseWriter, r *http.Request) {
	var link createLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		s.logger.WithError(err).Error("Error while decoding request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if link.ID == "" {
		id, err := s.idGenerator.GenerateID()
		if err != nil {
			s.logger.WithError(err).Error("Error while generating ID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		link.ID = id
	}

	if err := s.db.Create(link.ID, link.URL); err != nil {
		s.logger.WithError(err).Error("Error while creating link")
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
		s.logger.WithError(err).Error("Error while deleting link")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

var (
	tmpl = template.Must(template.ParseFiles("./internal/ui/template/admin-page.html"))
)

func (s *Shortener) handleAdmin(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, nil); err != nil {
		s.logger.Println("[ERROR] error while executing template: %v", err)
	}
}
