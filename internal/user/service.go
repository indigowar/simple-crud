package user

import (
	"context"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
)

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
type Repository interface {
	Create(context.Context, User) error
	GetByID(context.Context, int64) (interface{}, error)
	GetAll(context.Context) (interface{}, error)
	Update(context.Context, User) (string, error)
	Delete(context.Context, int64) (string, error)
}

type userservice struct {
	repository Repository
	logger     log.Logger
}

type UserService interface {
	Create(context.Context, User) (string, error)
	GetByID(context.Context, int64) (interface{}, error)
	GetAll(context.Context) (interface{}, error)
	Update(context.Context, User) (string, error)
	Delete(context.Context, int64) (string, error)
}

func NewService(rep Repository, logger log.Logger) UserService {
	return &userservice{
		repository: rep,
		logger:     logger,
	}
}

func (s userservice) Create(ctx context.Context, u User) (string, error) {
	logger := log.With(s.logger, "method", "Create")

	if err := s.repository.Create(ctx, u); err != nil {
		level.Error(logger).Log("err from repo is ", err)
		return "", err
	}

	return "success", nil
}

func (s userservice) GetByID(ctx context.Context, id int64) (interface{}, error) {
	logger := log.With(s.logger, "method", "GetByID")

	u, err := s.repository.GetByID(ctx, id)

	if err != nil {
		level.Error(logger).Log("err ", err)
		return nil, err
	}
	return u, nil
}

func (s userservice) GetAll(ctx context.Context) (interface{}, error) {
	logger := log.With(s.logger, "method", "GetAll")

	customer, err := s.repository.GetAll(ctx)

	if err != nil {
		level.Error(logger).Log("err ", err)
		return nil, err
	}

	return customer, nil
}

func (s userservice) Update(ctx context.Context, u User) (string, error) {
	logger := log.With(s.logger, "method", "Update")

	msg, err := s.repository.Update(ctx, u)

	if err != nil {
		level.Error(logger).Log("err from repo is ", err)
		return "", err
	}

	return msg, nil
}

func (s userservice) Delete(ctx context.Context, id int64) (string, error) {
	logger := log.With(s.logger, "method", "Delete")

	msg, err := s.repository.Delete(ctx, id)

	if err != nil {
		level.Error(logger).Log("err ", err)
		return "", err
	}

	return msg, nil
}
