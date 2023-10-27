// Package apis provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package apis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oapi-codegen/runtime"
)

const (
	JWTScopes = "JWT.Scopes"
)

// AdminLoginRequest defines model for AdminLoginRequest.
type AdminLoginRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// AdminLoginResponse defines model for AdminLoginResponse.
type AdminLoginResponse struct {
	Token string `json:"token"`
}

// CreateShortLinkRequest defines model for CreateShortLinkRequest.
type CreateShortLinkRequest struct {
	ID  *string `json:"id,omitempty"`
	URL string  `json:"url"`
}

// CreateShortLinkResponse defines model for CreateShortLinkResponse.
type CreateShortLinkResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// GetLinkMetricsResponse defines model for GetLinkMetricsResponse.
type GetLinkMetricsResponse struct {
	ID      string      `json:"id"`
	Metrics LinkMetrics `json:"metrics"`
	URL     string      `json:"url"`
}

// GetLinksResponse defines model for GetLinksResponse.
type GetLinksResponse struct {
	Links *[]Link `json:"links,omitempty"`
}

// Link defines model for Link.
type Link struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// LinkMetrics defines model for LinkMetrics.
type LinkMetrics struct {
	Clicks int `json:"clicks"`
}

// LoginAdminJSONRequestBody defines body for LoginAdmin for application/json ContentType.
type LoginAdminJSONRequestBody = AdminLoginRequest

// CreateNewLinkJSONRequestBody defines body for CreateNewLink for application/json ContentType.
type CreateNewLinkJSONRequestBody = CreateShortLinkRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /api/v1/admin/login)
	LoginAdmin(w http.ResponseWriter, r *http.Request)

	// (GET /api/v1/links)
	GetAllLinks(w http.ResponseWriter, r *http.Request)

	// (POST /api/v1/links)
	CreateNewLink(w http.ResponseWriter, r *http.Request)
	// Delete link
	// (DELETE /api/v1/links/{linkId})
	DeleteShortLink(w http.ResponseWriter, r *http.Request, linkID string)
	// Get Link Metrics
	// (GET /api/v1/links/{linkId})
	GetLinkMetrics(w http.ResponseWriter, r *http.Request, linkID string)
	// Redirect to link
	// (GET /{linkId})
	GetLinkById(w http.ResponseWriter, r *http.Request, linkID string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// LoginAdmin operation middleware
func (siw *ServerInterfaceWrapper) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.LoginAdmin(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetAllLinks operation middleware
func (siw *ServerInterfaceWrapper) GetAllLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAllLinks(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateNewLink operation middleware
func (siw *ServerInterfaceWrapper) CreateNewLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, JWTScopes, []string{"admin"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateNewLink(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// DeleteShortLink operation middleware
func (siw *ServerInterfaceWrapper) DeleteShortLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "linkId" -------------
	var linkID string

	err = runtime.BindStyledParameter("simple", false, "linkId", mux.Vars(r)["linkId"], &linkID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "linkId", Err: err})
		return
	}

	ctx = context.WithValue(ctx, JWTScopes, []string{"admin"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteShortLink(w, r, linkID)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetLinkMetrics operation middleware
func (siw *ServerInterfaceWrapper) GetLinkMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "linkId" -------------
	var linkID string

	err = runtime.BindStyledParameter("simple", false, "linkId", mux.Vars(r)["linkId"], &linkID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "linkId", Err: err})
		return
	}

	ctx = context.WithValue(ctx, JWTScopes, []string{"admin"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetLinkMetrics(w, r, linkID)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetLinkById operation middleware
func (siw *ServerInterfaceWrapper) GetLinkById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "linkId" -------------
	var linkID string

	err = runtime.BindStyledParameter("simple", false, "linkId", mux.Vars(r)["linkId"], &linkID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "linkId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetLinkById(w, r, linkID)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{})
}

type GorillaServerOptions struct {
	BaseURL          string
	BaseRouter       *mux.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r *mux.Router) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r *mux.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options GorillaServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = mux.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.HandleFunc(options.BaseURL+"/api/v1/admin/login", wrapper.LoginAdmin).Methods("POST")

	r.HandleFunc(options.BaseURL+"/api/v1/links", wrapper.GetAllLinks).Methods("GET")

	r.HandleFunc(options.BaseURL+"/api/v1/links", wrapper.CreateNewLink).Methods("POST")

	r.HandleFunc(options.BaseURL+"/api/v1/links/{linkId}", wrapper.DeleteShortLink).Methods("DELETE")

	r.HandleFunc(options.BaseURL+"/api/v1/links/{linkId}", wrapper.GetLinkMetrics).Methods("GET")

	r.HandleFunc(options.BaseURL+"/{linkId}", wrapper.GetLinkById).Methods("GET")

	return r
}
