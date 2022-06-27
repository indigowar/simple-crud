package main

import (
	"net/http"
	"os"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/indigowar/simple-crud/internal/user"
	"github.com/indigowar/simple-crud/pkg/db"
	"github.com/indigowar/simple-crud/pkg/middleware/metrics"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	r := mux.NewRouter()

	metricsMiddleware := metrics.NewMetricsMiddleware()

	db := db.GetDB()

	var userSrvc user.UserService

	{
		repo, err := user.NewRepo(db, logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		userSrvc = user.NewService(repo, logger)
	}

	createUserHandler := httptransport.NewServer(
		user.MakeCreateEndpoint(userSrvc),
		user.DecodeCreateRequest,
		user.EncodeResponse,
	)

	getUserByIDHandler := httptransport.NewServer(
		user.MakeGetByIDEndpoint(userSrvc),
		user.DecodeGetByIDRequest,
		user.EncodeResponse,
	)

	getAllUsersHandler := httptransport.NewServer(
		user.MakeGetAllEndpoint(userSrvc),
		user.DecodeGetAllRequest,
		user.EncodeResponse,
	)

	deleteUserHandler := httptransport.NewServer(
		user.MakeDeleteEndpoint(userSrvc),
		user.DecodeDeleteRequest,
		user.EncodeResponse,
	)

	updateUserHandler := httptransport.NewServer(
		user.MakeUpdateEndpoint(userSrvc),
		user.DecodeUpdateRequest,
		user.EncodeResponse,
	)

	http.Handle("/", r)
	r.Handle("/metrics", promhttp.Handler())

	r.Handle("/user", getAllUsersHandler).Methods(http.MethodGet)
	r.Handle("/user/{id}", getUserByIDHandler).Methods(http.MethodGet)
	r.Handle("/user", createUserHandler).Methods(http.MethodPut)
	r.Handle("/user", updateUserHandler).Methods(http.MethodPost)
	r.Handle("/user/{id}", deleteUserHandler).Methods(http.MethodDelete)
	r.Use(metricsMiddleware.Metrics)

	logger.Log("msg", "HTTP", "addr", ":8000")
	logger.Log("err", http.ListenAndServe(":800", nil))
}
