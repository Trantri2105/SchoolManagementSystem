package postgres

import (
	error2 "SchoolManagement/error"
	"SchoolManagement/model"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"reflect"
	"strings"
)

type UserRepo interface {
	InsertUser(ctx context.Context, user model.User, tx *sqlx.Tx) error
	GetUserById(ctx context.Context, id string, tx *sqlx.Tx) (model.User, error)
	UpdateUser(ctx context.Context, user model.User, tx *sqlx.Tx) error
	DeleteUserById(ctx context.Context, id string, tx *sqlx.Tx) error
}

type userRepo struct {
	db *sqlx.DB
}

func (u *userRepo) InsertUser(ctx context.Context, user model.User, tx *sqlx.Tx) error {
	query := `INSERT INTO users(id, name, date_of_birth, gender, email, identity_number, phone_number, address, password, role) 
			VALUES (:id, :name, :date_of_birth, :gender, :email, :identity_number, :phone_number, :address, :password, :role)`

	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, user)
	} else {
		_, err = u.db.ExecContext(ctx, query, user)
	}

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return &error2.UniqueConstraintErr{Message: pqErr.Constraint}
		}
		log.Println("User repo, save user err :", err)
		return err
	}

	return nil
}

func (u *userRepo) GetUserById(ctx context.Context, id string, tx *sqlx.Tx) (model.User, error) {
	query := `SELECT id, name, date_of_birth, gender, email, identity_number, phone_number, address, password, role FROM users WHERE id=$1`

	var row *sqlx.Row
	if tx != nil {
		row = tx.QueryRowxContext(ctx, query, id)
	} else {
		row = u.db.QueryRowxContext(ctx, query, id)
	}

	var user model.User
	err := row.StructScan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, &error2.ResourceNotFoundErr{Resource: "User"}
		}
		log.Println("User repo, get user err :", err)
		return model.User{}, err
	}

	return user, nil
}

func (u *userRepo) UpdateUser(ctx context.Context, user model.User, tx *sqlx.Tx) error {
	var updateFields []string
	t := reflect.TypeOf(user)
	v := reflect.ValueOf(user)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if !value.IsZero() {
			updateFields = append(updateFields, field.Tag.Get("db")+" = :"+field.Tag.Get("db"))
		}
	}
	if len(updateFields) == 0 {
		return nil
	}
	query := `UPDATE users SET ` + strings.Join(updateFields, ",") + ` WHERE id = :id`

	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, user)
	} else {
		_, err = u.db.ExecContext(ctx, query, user)
	}
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return &error2.UniqueConstraintErr{Message: pqErr.Constraint}
		}
		log.Println("User repo, update user err :", err)
		return err
	}
	return nil
}

func (u *userRepo) DeleteUserById(ctx context.Context, id string, tx *sqlx.Tx) error {
	query := `DELETE FROM users WHERE id = $1`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, id)
	} else {
		_, err = u.db.ExecContext(ctx, query, id)
	}
	if err != nil {
		log.Println("User repo, delete user err :", err)
		return err
	}
	return nil
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepo{db: db}
}
