package api

// The profilesvc is just over HTTP, so we just have a single transport.go.

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

)

func encodePostProfileRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/profiles/")
	req.Method, req.URL.Path = "POST", "/profiles/"
	return encodeRequest(ctx, req, request)
}

func encodeGetProfileRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/profiles/{id}")
	r := request.(GetProfileRequest)
	profileID := url.QueryEscape(r.ID)
	req.Method, req.URL.Path = "GET", "/profiles/"+profileID
	return encodeRequest(ctx, req, request)
}

func encodePutProfileRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PUT").Path("/profiles/{id}")
	r := request.(PutProfileRequest)
	profileID := url.QueryEscape(r.ID)
	req.Method, req.URL.Path = "PUT", "/profiles/"+profileID
	return encodeRequest(ctx, req, request)
}

func encodePatchProfileRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PATCH").Path("/profiles/{id}")
	r := request.(PatchProfileRequest)
	profileID := url.QueryEscape(r.ID)
	req.Method, req.URL.Path = "PATCH", "/profiles/"+profileID
	return encodeRequest(ctx, req, request)
}

func encodeDeleteProfileRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("DELETE").Path("/profiles/{id}")
	r := request.(DeleteProfileRequest)
	profileID := url.QueryEscape(r.ID)
	req.Method, req.URL.Path = "DELETE", "/profiles/"+profileID
	return encodeRequest(ctx, req, request)
}

func encodeGetAddressesRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/profiles/{id}/addresses/")
	r := request.(GetAddressesRequest)
	profileID := url.QueryEscape(r.ProfileID)
	req.Method, req.URL.Path = "GET", "/profiles/"+profileID+"/addresses/"
	return encodeRequest(ctx, req, request)
}

func encodeGetAddressRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/profiles/{id}/addresses/{addressID}")
	r := request.(GetAddressRequest)
	profileID := url.QueryEscape(r.ProfileID)
	addressID := url.QueryEscape(r.AddressID)
	req.Method, req.URL.Path = "GET", "/profiles/"+profileID+"/addresses/"+addressID
	return encodeRequest(ctx, req, request)
}

func encodePostAddressRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/profiles/{id}/addresses/")
	r := request.(PostAddressRequest)
	profileID := url.QueryEscape(r.ProfileID)
	req.Method, req.URL.Path = "POST", "/profiles/"+profileID+"/addresses/"
	return encodeRequest(ctx, req, request)
}

func encodeDeleteAddressRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("DELETE").Path("/profiles/{id}/addresses/{addressID}")
	r := request.(DeleteAddressRequest)
	profileID := url.QueryEscape(r.ProfileID)
	addressID := url.QueryEscape(r.AddressID)
	req.Method, req.URL.Path = "DELETE", "/profiles/"+profileID+"/addresses/"+addressID
	return encodeRequest(ctx, req, request)
}

func decodePostProfileResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response PostProfileResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetProfileResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response GetProfileResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePutProfileResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response PutProfileResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePatchProfileResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response PatchProfileResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteProfileResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response DeleteProfileResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetAddressesResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response GetAddressesResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetAddressResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response GetAddressResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePostAddressResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response PostAddressResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteAddressResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response DeleteAddressResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// profilesvc endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}




