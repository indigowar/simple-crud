package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-kit/log"
)

var (
	ErrRepo                = errors.New("unable to handle Repo Request")
	ErrIDNotFound          = errors.New("id is not found")
	ErrPhoneNumberNotFound = errors.New("phone number is not found")
)

type repo struct {
	db     *sql.DB
	logger log.Logger
}

func NewRepo(db *sql.DB, logger log.Logger) (Repository, error) {
	return &repo{
		db:     db,
		logger: logger,
	}, nil
}

func (r repo) Create(ctx context.Context, u User) error {
	stmt := `DO $$
	BEGIN
	CREATE TABLE IF NOT EXISTS Users (
	id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 124567895123 CACHE 1 ),
	email character varying COLLATE pg_catalog."default",
	phone character varying COLLATE pg_catalog."default",
	CONSTRAINT users_pkey PRIMARY KEY (id)
	);
	END;
	$$;`
	_, err := r.db.Exec(stmt)
	fmt.Println(err)

	_, err = r.db.ExecContext(ctx, "INSERT INTO Users(email, phone) VALUES ($1, $2)", u.Email, u.Phone)
	fmt.Println(err)
	if err != nil {
		fmt.Println("Error occured inside Create in repo")
		return err
	} else {
		fmt.Println("User Created:", u.Email)
	}
	return nil
}

func (r repo) GetByID(ctx context.Context, id int64) (interface{}, error) {
	u := User{}

	err := r.db.QueryRowContext(ctx, "SELECT u.id,u.email,u.phone FROM Users as u where u.id = $1", id).Scan(&u.ID, &u.Email, &u.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, ErrIDNotFound
		}
		return u, err
	}
	return u, nil
}

func (r repo) GetAll(ctx context.Context) (interface{}, error) {
	user := User{}
	var res []interface{}
	rows, err := r.db.QueryContext(ctx, "SELECT u.id,u.email,u.phone FROM Users as u ")
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrIDNotFound
		}
		return user, err
	}

	defer rows.Close()
	for rows.Next() {
		_ = rows.Scan(&user.ID, &user.Email, &user.Phone)
		res = append([]interface{}{user}, res...)
	}
	return res, nil
}

func (r repo) Update(ctx context.Context, u User) (string, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE Users SET email=$1 , phone = $2 WHERE id = $3", u.Email, u.Phone, u.ID)
	if err != nil {
		return "", err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return "", err
	}
	if rowCnt == 0 {
		return "", ErrIDNotFound
	}

	return "successfully updated", err
}

func (r repo) Delete(ctx context.Context, id int64) (string, error) {
	res, err := r.db.ExecContext(ctx, "DELETE FROM Users WHERE id = $1 ", id)
	if err != nil {
		return "", err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return "", err
	} else if rowCnt == 0 {
		return "", ErrIDNotFound
	}
	return "Successfully deleted ", nil
}
