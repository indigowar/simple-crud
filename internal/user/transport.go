package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
)

func MakeCreateEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateRequest)
		msg, err := s.Create(ctx, req.u)

		return CreateResponse{Message: msg, Err: err}, nil
	}
}

func MakeGetByIDEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetByIDRequest)
		u, err := s.GetByID(ctx, req.ID)
		if err != nil {
			return GetByIDResponse{U: nil, Err: err}, nil
		}
		return GetByIDResponse{U: u, Err: nil}, nil
	}
}

func MakeGetAllEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		users, err := s.GetAll(ctx)
		if err != nil {
			return GetAllResponse{User: nil, Err: errors.New("nothing was found")}, nil
		}
		return GetAllResponse{User: users, Err: nil}, nil
	}
}

func MakeDeleteEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteRequest)

		msg, err := s.Delete(ctx, req.ID)

		return DeleteResponse{Msg: msg, Err: err}, nil
	}
}

func MakeUpdateEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateRequest)

		msg, err := s.Update(ctx, req.u)

		return UpdateResponse{Msg: msg, Err: err}, nil
	}
}

func DecodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req.u); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeGetByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strId := vars["id"]
	intId, _ := strconv.Atoi(strId)
	req := GetByIDRequest{
		ID: int64(intId),
	}
	return req, nil
}

func DecodeGetAllRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GetAllRequest
	return req, nil
}

func DecodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strId := vars["customerid"]
	intId, _ := strconv.Atoi(strId)
	req := DeleteRequest{
		ID: int64(intId),
	}
	return req, nil
}

func DecodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req.u); err != nil {
		return nil, err
	}
	return req, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type (
	CreateRequest struct {
		u User
	}

	CreateResponse struct {
		Message string `json:"message"`
		Err     error  `json:"error,omitempty"`
	}

	GetByIDRequest struct {
		ID int64 `json:"id"`
	}

	GetByIDResponse struct {
		U   interface{} `json:"user"`
		Err error       `json:"error,omitempty"`
	}

	GetAllRequest struct{}

	GetAllResponse struct {
		User interface{} `json:"user,omitempty"`
		Err  error       `json:"error,omitempty"`
	}

	DeleteRequest struct {
		ID int64 `json:"id"`
	}

	DeleteResponse struct {
		Msg string `json:"response"`
		Err error  `json:"error,omitempty"`
	}

	UpdateRequest struct {
		u User
	}

	UpdateResponse struct {
		Msg string `json:"status,omitempty"`
		Err error  `json:"error,omitempty"`
	}
)
