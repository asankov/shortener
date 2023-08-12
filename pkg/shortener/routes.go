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
	// TODO: auth
	router.HandleFunc("/admin", s.handleAdmin).Methods(http.MethodGet)
	router.HandleFunc("/{id}", s.handleGetLink).Methods(http.MethodGet)

	router.HandleFunc("/admin/login", s.handleAdminLoginPage).Methods(http.MethodGet)
	router.HandleFunc("/admin/login", s.handleAdminLogin).Methods(http.MethodPost)

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

		s.logger.Warn("unknown error while getting link by id", "link_id", id, "error", err)
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
		s.logger.Error("Error while decoding request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if link.ID == "" {
		id, err := s.idGenerator.GenerateID()
		if err != nil {
			s.logger.Error("Error while generating ID", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		link.ID = id
	}

	if err := s.db.Create(link.ID, link.URL); err != nil {
		s.logger.Error("Error while creating link", "error", err)
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
		s.logger.Error("Error while deleting link", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

var (
	tmpl               = template.Must(template.ParseFiles("./internal/ui/template/admin-page.html"))
	adminLoginPageTmpl = template.Must(template.ParseFiles("./internal/ui/template/admin-login.html"))
)

func (s *Shortener) handleAdmin(w http.ResponseWriter, r *http.Request) {
	type pageData struct {
		Links []*links.Link
	}

	links, err := s.db.GetAll()
	if err != nil {
		s.logger.Error("error while getting all links", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	rr := pageData{
		Links: links,
	}
	if err := tmpl.Execute(w, rr); err != nil {
		s.logger.Error("Error while executing template", "error", err)
	}
}

func (s *Shortener) handleAdminLoginPage(w http.ResponseWriter, r *http.Request) {
	if err := adminLoginPageTmpl.Execute(w, r); err != nil {
		s.logger.Error("Error while executing template", "error", err)
	}
}

func (s *Shortener) handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Error("error while decoding request body", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	user, err := s.userService.Get(req.Email, req.Password)
	if err != nil {
		s.logger.Error("error while getting user", "error", err, "email", req.Email)
		w.WriteHeader(http.StatusInternalServerError)
	}

	token, err := s.authenticator.NewTokenForUser(user)
	if err != nil {
		s.logger.Error("error while generating token for user", "error", err, "email", req.Email)
		w.WriteHeader(http.StatusInternalServerError)
	}

	type response struct {
		Token string `json:"token"`
	}

	if err := json.NewEncoder(w).Encode(response{Token: token}); err != nil {
		s.logger.Error("error while encoding response", "error", err, "email", req.Email)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
