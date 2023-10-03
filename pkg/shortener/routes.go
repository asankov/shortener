package shortener

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"

	"github.com/asankov/shortener/internal/apis"
	"github.com/asankov/shortener/pkg/links"
	"github.com/gorilla/mux"
)

func (s *Shortener) routes() http.Handler {
	router := mux.NewRouter()
	// // TODO: auth
	router.HandleFunc("/admin", s.handleAdmin).Methods(http.MethodGet)

	router.HandleFunc("/admin/login", s.handleAdminLoginPage).Methods(http.MethodGet)

	// TODO: auth
	apiRoutes := router.PathPrefix("/api/v1").Subrouter()
	apiRoutes.HandleFunc("/links/{id}", s.handleDeleteLink).Methods(http.MethodDelete)

	// TODO
	// apiRoutes.HandleFunc("/{id}", s.handleGetLinkMetrics).Methods(http.MethodGet)

	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fs))

	return apis.HandlerFromMux(s, router)

	// return router
}

func (s *Shortener) GetLinkById(w http.ResponseWriter, r *http.Request, linkId string) {
	link, err := s.db.GetByID(linkId)
	if err != nil {
		if errors.Is(err, links.ErrLinkNotFound) {
			w.WriteHeader(http.StatusNotFound)
		}

		s.logger.Warn("unknown error while getting link by id", "link_id", linkId, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: increment the number of clicks

	http.Redirect(w, r, link.URL, http.StatusFound)
}

func (s *Shortener) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	var req apis.AdminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Error("error while decoding request body", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	user, err := s.userService.Get(req.Username, req.Password)
	if err != nil {
		s.logger.Error("error while getting user", "error", err, "username", req.Username)
		w.WriteHeader(http.StatusInternalServerError)
	}

	token, err := s.authenticator.NewTokenForUser(user)
	if err != nil {
		s.logger.Error("error while generating token for user", "error", err, "username", req.Username)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(apis.AdminLoginResponse{Token: token}); err != nil {
		s.logger.Error("error while encoding response", "error", err, "email", req.Username)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Shortener) CreateNewLink(w http.ResponseWriter, r *http.Request) {
	var link apis.CreateShortLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		s.logger.Error("Error while decoding request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if link.Id == nil {
		id, err := s.idGenerator.GenerateID()
		if err != nil {
			s.logger.Error("Error while generating ID", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		link.Id = &id
	}

	if err := s.db.Create(*link.Id, link.Url); err != nil {
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
