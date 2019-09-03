package leagues

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/thethan/fantasydraftroom/internal/middleware"
	"github.com/thethan/fantasydraftroom/internal/users"
	"io/ioutil"
	"net/http"
)

const APIToken = "api_token"
const USER = "user"
const PARAM_DRAFTID = "draftID"


// NewHTTPHandler returns an HTTP handler that makes a set of endpoints
// available on predefined paths.
func NewHTTPHandler(router *mux.Router, endpoints Set, logger log.Logger) http.Handler {
	// Zipkin HTTP Server Trace can either be instantiated per endpoint with a
	// provided operation name or a global tracing service can be instantiated
	// without an operation name and fed to each Go kit endpoint as ServerOption.
	// In the latter case, the operation name will be the endpoint's http method.
	// We demonstrate a global tracing service here.
	//zipkinServer := zipkin.HTTPServerTrace(zipkinTracer)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
		httptransport.ServerBefore(users.GetBearerTokenFromHeaderToContext),
	}
	router.Methods("GET").Path("/leagues/{id}").Handler(httptransport.NewServer(
		endpoints.LeagueEndpoint,
		decodeLeagueID,
		encodeHTTPGenericResponse,
		options...,
	))


	return router
}


func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {

	w.Header().Set(middleware.HdrAllowOriginHeader, "*")
	w.Header().Set(middleware.HdrRequestMethodHeader, "POST, GET, OPTIONS, HEAD, PATCH, DELETE")
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func err2code(err error) int {

	return http.StatusBadRequest
}

func errorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

type errorWrapper struct {
	Error string `json:"error"`
}


func decodeLeagueID(ctx context.Context, r *http.Request) (interface{}, error) {
	var leageRequest LeagueRequest
	mapOfVars := mux.Vars(r)
	leagueID, ok := mapOfVars["id"]
	if !ok {
		return nil, errors.New("Could not get leagueID from uri params")
	}
	leageRequest.LeagueID = leagueID
	return leageRequest, nil
}

// encodeHTTPGenericRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// encodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}

	w.Header().Set(middleware.HdrAllowOriginHeader, "*")
	w.Header().Set(middleware.HdrRequestMethodHeader, "POST, GET, OPTIONS, HEAD, PATCH, DELETE")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
