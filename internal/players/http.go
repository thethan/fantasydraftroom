package players

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/thethan/fantasydraftroom/internal/middleware"
	"github.com/thethan/fantasydraftroom/internal/users"
	"io/ioutil"
	"net/http"
	"strconv"
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
	//
	router.Methods("GET").Path("/players/drafts/{draftID}").Handler(httptransport.NewServer(
		endpoints.PlayersOrder,
		decodeHTTPGetPlayerOrderList,
		encodeHTTPGenericResponse,
		options...,
	))

	router.Methods("POST").Path("/players/drafts/{draftID}/preference").Handler(httptransport.NewServer(
		endpoints.PlayerPreference,
		decodeHTTPPostPlayerPreferenceList,
		encodeHTTPGenericResponse,
		options...,
	))

	router.Methods("GET", "POST", "PUT").Path("/yahoo/callback").Handler(httptransport.NewServer(
		endpoints.LoginEndpoint,
		decodeHTTPYahoo,
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

func decodeHTTPGetPlayerOrderList(ctx context.Context, r *http.Request) (interface{}, error) {
	var req DraftPlayerRankingsRequest
	vars := mux.Vars(r)

	draftID, ok := vars[PARAM_DRAFTID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("could not get a %s from request", PARAM_DRAFTID))
	}

	req.DraftID, _ = strconv.Atoi(draftID)

	return req, nil
}

func decodeHTTPPostPlayerPreferenceList(ctx context.Context, r *http.Request) (interface{}, error) {
	var req UserPlayerPreferenceRequest
	vars := mux.Vars(r)

	draftID, ok := vars[PARAM_DRAFTID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("could not get a %s from request", PARAM_DRAFTID))
	}

	req.DraftID, _ = strconv.Atoi(draftID)

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	//var intIDs []int


	err = json.Unmarshal(b, &req)
	if err != nil {
		return nil, err
	}
	//
	//req.Body.PlayerIDs = make([]PlayerID, len(intIDs))
	//for idx := range intIDs {
	//	req.PlayerIDs = append(req.PlayerIDs, req.PlayerIDs[idx])
	//}

	fmt.Printf("Request Length %d \n", len(req.Body.PlayerIDs))
	fmt.Printf("%+v \n", req)

	return req, nil
}



func decodeHTTPYahoo(ctx context.Context, r *http.Request) (interface{}, error) {
	var req UserPlayerPreferenceRequest

	//var intIDs []int

	//
	//req.Body.PlayerIDs = make([]PlayerID, len(intIDs))
	//for idx := range intIDs {
	//	req.PlayerIDs = append(req.PlayerIDs, req.PlayerIDs[idx])
	//}

	return req, nil
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
