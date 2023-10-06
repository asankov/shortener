package shortener

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/asankov/shortener/internal/apis"
	"github.com/asankov/shortener/internal/links"
)

func (s *Shortener) routes() http.Handler {
	return apis.HandlerWithOptions(s.handler, apis.GorillaServerOptions{
		Middlewares: []apis.MiddlewareFunc{
			s.handler.authenticated,
		},
	})
}

func (h *handler) GetLinkById(w http.ResponseWriter, r *http.Request, linkId string) {
	link, err := h.db.GetByID(linkId)
	if err != nil {
		if errors.Is(err, links.ErrLinkNotFound) {
			w.WriteHeader(http.StatusNotFound)
		}

		h.logger.Warn("unknown error while getting link by id", "link_id", linkId, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: increment the number of clicks

	http.Redirect(w, r, link.URL, http.StatusFound)
}

func (h *handler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	var req apis.AdminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("error while decoding request body", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := h.userService.GetUser(req.Username, req.Password)
	if err != nil {
		h.logger.Error("error while getting user", "error", err, "username", req.Username)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := h.authenticator.NewTokenForUser(user)
	if err != nil {
		h.logger.Error("error while generating token for user", "error", err, "username", req.Username)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	if err := json.NewEncoder(w).Encode(apis.AdminLoginResponse{Token: token}); err != nil {
		h.logger.Error("error while encoding response", "error", err, "email", req.Username)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) CreateNewLink(w http.ResponseWriter, r *http.Request) {
	var link apis.CreateShortLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		h.logger.Error("Error while decoding request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if link.Id == nil {
		id, err := h.idGenerator.GenerateID()
		if err != nil {
			h.logger.Error("Error while generating ID", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		link.Id = &id
	}

	if err := h.db.Create(*link.Id, link.Url); err != nil {
		if errors.Is(err, links.ErrLinkAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		h.logger.Error("Error while creating link", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handler) GetLinkMetrics(w http.ResponseWriter, r *http.Request, linkId string) {
	panic("not implemented")
}

func (h *handler) DeleteShortLink(w http.ResponseWriter, r *http.Request, linkID string) {
	if err := h.db.Delete(linkID); err != nil {
		h.logger.Error("Error while deleting link", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
