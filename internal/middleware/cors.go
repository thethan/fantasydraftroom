package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	HdrContentType            = "Content-Type"
	HdrContentLen             = "Content-Length"
	HdrAcceptEncoding         = "Accept-Encoding"
	HdrContentEncoding        = "Content-Encoding"
	HdrVary                   = "Vary"
	HdrAllowOriginHeader      = "Access-Control-Allow-Origin"
	HdrOptionMethod           = "OPTIONS"
	HdrExposeHeadersHeader    = "Access-Control-Expose-Headers"
	HdrMaxAgeHeader           = "Access-Control-Max-Age"
	HdrAllowMethodsHeader     = "Access-Control-Allow-Methods"
	HdrAllowHeadersHeader     = "Access-Control-Allow-Headers"
	HdrAllowCredentialsHeader = "Access-Control-Allow-Credentials"
	HdrRequestMethodHeader    = "Access-Control-Request-Method"
	HdrRequestHeadersHeader   = "Access-Control-Request-Headers"
	HdrOriginHeader           = "Origin"
	MatchAll                  = "*"
	Gzip                      = "gzip"
	Deflate                   = "deflate"
	HdrAccept                 = "Accept"
	HdrAcceptLanguage         = "Accept-Language"
	HdrContentLanguage        = "Content-Language"
	HdrOrigin                 = "Origin"
)
// CORSOption represents a functional option for configuring the CORS middleware.
type CORSOption func(*cors) error

type cors struct {
	h                      http.Handler
	allowedHeaders         []string
	allowedMethods         []string
	allowedOrigins         []string
	allowedOriginValidator OriginValidator
	exposedHeaders         []string
	maxAge                 int
	ignoreOptions          bool
	allowCredentials       bool
	optionStatusCode       int
}

// OriginValidator takes an origin string and returns whether or not that origin is allowed.
type OriginValidator func(string) bool

var (
	defaultCorsOptionStatusCode = 200
	defaultCorsMethods          = []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPatch}
	defaultCorsHeaders          = []string{HdrAccept, HdrAcceptLanguage, HdrContentLanguage, HdrOrigin}
	// (WebKit/Safari v9 sends the Origin header by default in AJAX requests)
)

func (ch *cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get(HdrOriginHeader)
	if !ch.isOriginAllowed(origin) {
		if r.Method != HdrOptionMethod || ch.ignoreOptions {
			ch.h.ServeHTTP(w, r)
		}

		return
	}

	if r.Method == HdrOptionMethod {
		if ch.ignoreOptions {
			ch.h.ServeHTTP(w, r)
			return
		}

		if _, ok := r.Header[HdrRequestMethodHeader]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		method := r.Header.Get(HdrRequestMethodHeader)
		if !ch.isMatch(method, ch.allowedMethods) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		requestHeaders := strings.Split(r.Header.Get(HdrRequestHeadersHeader), ",")
		var allowedHeaders []string
		for _, v := range requestHeaders {
			canonicalHeader := http.CanonicalHeaderKey(strings.TrimSpace(v))
			if canonicalHeader == "" || ch.isMatch(canonicalHeader, defaultCorsHeaders) {
				continue
			}

			if !ch.isMatch(canonicalHeader, ch.allowedHeaders) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			allowedHeaders = append(allowedHeaders, canonicalHeader)
		}

		if len(allowedHeaders) > 0 {
			w.Header().Set(HdrAllowHeadersHeader, strings.Join(allowedHeaders, ","))
		}

		if ch.maxAge > 0 {
			w.Header().Set(HdrMaxAgeHeader, strconv.Itoa(ch.maxAge))
		}

		if !ch.isMatch(method, defaultCorsMethods) {
			w.Header().Set(HdrAllowMethodsHeader, method)
		}
	} else {
		if len(ch.exposedHeaders) > 0 {
			w.Header().Set(HdrExposeHeadersHeader, strings.Join(ch.exposedHeaders, ","))
		}
	}

	if ch.allowCredentials {
		w.Header().Set(HdrAllowCredentialsHeader, "true")
	}

	if len(ch.allowedOrigins) > 1 {
		w.Header().Set(HdrVary, HdrOriginHeader)
	}

	returnOrigin := origin
	if ch.allowedOriginValidator == nil && len(ch.allowedOrigins) == 0 {
		returnOrigin = "*"
	} else {
		for _, o := range ch.allowedOrigins {
			// A configuration of * is different than explicitly setting an allowed
			// origin. Returning arbitrary origin headers in an access control allow
			// origin header is unsafe and is not required by any use case.
			if o == MatchAll {
				returnOrigin = "*"
				break
			}
		}
	}
	w.Header().Set(HdrAllowOriginHeader, returnOrigin)

	if r.Method == HdrOptionMethod {
		w.WriteHeader(ch.optionStatusCode)
		return
	}
	ch.h.ServeHTTP(w, r)
}

// CORS provides Cross-Origin Resource Sharing middleware.
func CORS(opts ...CORSOption) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		ch := parseCORSOptions(opts...)
		ch.h = h
		return ch
	}
}

func parseCORSOptions(opts ...CORSOption) *cors {
	ch := &cors{
		allowedMethods:   defaultCorsMethods,
		allowedHeaders:   defaultCorsHeaders,
		allowedOrigins:   []string{"*"},
		optionStatusCode: defaultCorsOptionStatusCode,
	}

	for _, option := range opts {
		_ = option(ch)
	}

	return ch
}

//
// Functional options for configuring CORS.
//

// AllowedHeaders adds the provided headers to the list of allowed headers in a
// CORS request.
// This is an append operation so the headers Accept, Accept-Language,
// and Content-Language are always allowed.
// Content-Type must be explicitly declared if accepting Content-Types other than
// application/x-www-form-urlencoded, multipart/form-data, or text/plain.
func AllowedHeaders(headers []string) CORSOption {
	return func(ch *cors) error {
		for _, v := range headers {
			normalizedHeader := http.CanonicalHeaderKey(strings.TrimSpace(v))
			if normalizedHeader == "" {
				continue
			}

			if !ch.isMatch(normalizedHeader, ch.allowedHeaders) {
				ch.allowedHeaders = append(ch.allowedHeaders, normalizedHeader)
			}
		}

		return nil
	}
}

// AllowedMethods can be used to explicitly allow methods in the
// Access-Control-Allow-Methods header.
// This is a replacement operation so you must also
// pass GET, HEAD, and POST if you wish to support those methods.
func AllowedMethods(methods []string) CORSOption {
	return func(ch *cors) error {
		ch.allowedMethods = []string{}
		for _, v := range methods {
			normalizedMethod := strings.ToUpper(strings.TrimSpace(v))
			if normalizedMethod == "" {
				continue
			}

			if !ch.isMatch(normalizedMethod, ch.allowedMethods) {
				ch.allowedMethods = append(ch.allowedMethods, normalizedMethod)
			}
		}

		return nil
	}
}

// AllowedOrigins sets the allowed origins for CORS requests, as used in the
// 'Allow-Access-Control-Origin' HTTP header.
// Note: Passing in a []string{"*"} will allow any domain.
func AllowedOrigins(origins []string) CORSOption {
	return func(ch *cors) error {
		for _, v := range origins {
			if v == MatchAll {
				ch.allowedOrigins = []string{MatchAll}
				return nil
			}
		}

		ch.allowedOrigins = origins
		return nil
	}
}

// AllowedOriginValidator sets a function for evaluating allowed origins in CORS requests, represented by the
// 'Allow-Access-Control-Origin' HTTP header.
func AllowedOriginValidator(fn OriginValidator) CORSOption {
	return func(ch *cors) error {
		ch.allowedOriginValidator = fn
		return nil
	}
}

// OptionStatusCode sets a custom status code on the OPTIONS requests.
// Default behaviour sets it to 200 to reflect best practices. This is option is not mandatory
// and can be used if you need a custom status code (i.e 204).
//
// More informations on the spec:
// https://fetch.spec.whatwg.org/#cors-preflight-fetch
func OptionStatusCode(code int) CORSOption {
	return func(ch *cors) error {
		ch.optionStatusCode = code
		return nil
	}
}

// ExposedHeaders can be used to specify headers that are available
// and will not be stripped out by the user-agent.
func ExposedHeaders(headers []string) CORSOption {
	return func(ch *cors) error {
		ch.exposedHeaders = []string{}
		for _, v := range headers {
			normalizedHeader := http.CanonicalHeaderKey(strings.TrimSpace(v))
			if normalizedHeader == "" {
				continue
			}

			if !ch.isMatch(normalizedHeader, ch.exposedHeaders) {
				ch.exposedHeaders = append(ch.exposedHeaders, normalizedHeader)
			}
		}

		return nil
	}
}

// MaxAge determines the maximum age (in seconds) between preflight requests. A
// maximum of 10 minutes is allowed. An age above this value will default to 10
// minutes.
func MaxAge(age int) CORSOption {
	return func(ch *cors) error {
		// Maximum of 10 minutes.
		if age > 600 {
			age = 600
		}

		ch.maxAge = age
		return nil
	}
}

// IgnoreOptions causes the CORS middleware to ignore OPTIONS requests, instead
// passing them through to the next handler. This is useful when your application
// or framework has a pre-existing mechanism for responding to OPTIONS requests.
func IgnoreOptions() CORSOption {
	return func(ch *cors) error {
		ch.ignoreOptions = true
		return nil
	}
}

// AllowCredentials can be used to specify that the user agent may pass
// authentication details along with the request.
func AllowCredentials() CORSOption {
	return func(ch *cors) error {
		ch.allowCredentials = true
		return nil
	}
}

func (ch *cors) isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	if ch.allowedOriginValidator != nil {
		return ch.allowedOriginValidator(origin)
	}

	if len(ch.allowedOrigins) == 0 {
		return true
	}

	for _, allowedOrigin := range ch.allowedOrigins {
		if allowedOrigin == origin || allowedOrigin == MatchAll {
			return true
		}
	}

	return false
}

func (ch *cors) isMatch(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}
