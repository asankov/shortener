package shortener

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/asankov/shortener/internal/apis"
	"github.com/asankov/shortener/internal/links"
)

func (s *Shortener) routes() http.Handler {
	return apis.Handler(s.handler)
}

func (s *handler) GetLinkById(w http.ResponseWriter, r *http.Request, linkId string) {
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

func (s *handler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
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

func (s *handler) CreateNewLink(w http.ResponseWriter, r *http.Request) {
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

func (s *handler) GetLinkMetrics(w http.ResponseWriter, r *http.Request, linkId string) {
	panic("not implemented")
}

func (s *handler) DeleteShortLink(w http.ResponseWriter, r *http.Request, linkID string) {
	if err := s.db.Delete(linkID); err != nil {
		s.logger.Error("Error while deleting link", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
