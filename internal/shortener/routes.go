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

	if err := h.db.IncrementClicks(linkId); err != nil {
		h.logger.Warn("error while incrementing number of clicks", "link_id", linkId, "error", err)
	}

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

	if link.ID == nil {
		id, err := h.idGenerator.GenerateID()
		if err != nil {
			h.logger.Error("Error while generating ID", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		link.ID = &id
	}

	if err := h.db.Create(*link.ID, link.URL); err != nil {
		if errors.Is(err, links.ErrLinkAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		h.logger.Error("Error while creating link", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(apis.CreateShortLinkResponse{
		ID:  *link.ID,
		URL: link.URL,
	}); err != nil {
		h.logger.Error("error while encoding response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) GetLinkMetrics(w http.ResponseWriter, r *http.Request, linkID string) {
	link, err := h.db.GetByID(linkID)
	if err != nil {
		if errors.Is(err, links.ErrLinkNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		h.logger.Warn("unknown error while getting link by id", "link_id", linkID, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(apis.GetLinkMetricsResponse{
		ID:  link.ID,
		URL: link.URL,
		Metrics: apis.LinkMetrics{
			Clicks: link.Metrics.Clicks,
		},
	}); err != nil {
		h.logger.Error("error while encoding response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) DeleteShortLink(w http.ResponseWriter, r *http.Request, linkID string) {
	if err := h.db.Delete(linkID); err != nil {
		h.logger.Error("Error while deleting link", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) GetAllLinks(w http.ResponseWriter, r *http.Request) {

	links, err := h.db.GetAll()
	if err != nil {
		h.logger.Error("Error while getting links", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := make([]apis.Link, 0, len(links))
	for _, link := range links {
		res = append(res, apis.Link{
			ID:  link.ID,
			URL: link.URL,
		})
	}

	if err := json.NewEncoder(w).Encode(apis.GetLinksResponse{Links: &res}); err != nil {
		h.logger.Error("error while encoding response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
