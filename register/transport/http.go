package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/liubo51617/user/register/endpoint"
	"net/http"
	"os"
)

var ErrBadRequest = errors.New("invalid request parameter")

func MakeHttpHandler(ctx context.Context, endpoints *endpoint.RegisterEndpoints) http.Handler {
	r := mux.NewRouter()
	kitlog := log.NewLogfmtLogger(os.Stderr)
	kitlog = log.With(kitlog, "ts", log.DefaultTimestampUTC)
	kitlog = log.With(kitlog, "caller", log.DefaultCaller)

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(kitlog)),
		kithttp.ServerErrorEncoder(encodeError) ,
	}

	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		EncodeJSONResponse,
		options...,
		))

	r.Methods("GET").Path("/discovery/name").Handler(kithttp.NewServer(
		endpoints.DiscoveryEndpoint,
		decodeDiscoveryRequest,
		EncodeJSONResponse,
		options...,
		))

	return r
}

func decodeDiscoveryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	serviceName := r.URL.Query().Get("serviceName")

	if serviceName == "" {
		return nil, ErrBadRequest
	}
	return &endpoint.DiscoveryRequest{ServiceName:serviceName,}, nil
}

func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpoint.HealthRequest{},nil
}

func EncodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
