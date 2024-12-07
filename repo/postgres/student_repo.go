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

type StudentRepo interface {
	GetStudentById(ctx context.Context, id string, tx *sqlx.Tx) (model.Student, error)
	UpdateStudent(ctx context.Context, student model.Student, tx *sqlx.Tx) error
	DeleteStudentById(ctx context.Context, id string, tx *sqlx.Tx) error
	InsertStudent(ctx context.Context, student model.Student, tx *sqlx.Tx) error
}

type studentRepo struct {
	db *sqlx.DB
}

func (s *studentRepo) GetStudentById(ctx context.Context, id string, tx *sqlx.Tx) (model.Student, error) {
	query := `SELECT users.id, name, date_of_birth, gender, email, identity_number, phone_number, address, password, role, school_year, major
			FROM users
			JOIN students s ON users.id = s.id
			WHERE users.id = $1`

	var row *sqlx.Row
	if tx != nil {
		row = tx.QueryRowxContext(ctx, query, id)
	} else {
		row = s.db.QueryRowxContext(ctx, query, id)
	}
	var student model.Student
	err := row.StructScan(&student)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return student, &error2.ResourceNotFoundErr{Resource: "Student"}
		}
		log.Println("Student repo, get student err:", err)
		return student, err
	}
	return student, nil
}

func (s *studentRepo) UpdateStudent(ctx context.Context, student model.Student, tx *sqlx.Tx) error {
	var updateFields []string
	t := reflect.TypeOf(student)
	v := reflect.ValueOf(student)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if field.Anonymous {
			continue
		}
		if !value.IsZero() {
			updateFields = append(updateFields, field.Tag.Get("db")+" = :"+field.Tag.Get("db"))
		}
	}
	if len(updateFields) == 0 {
		return nil
	}
	query := `UPDATE students SET ` + strings.Join(updateFields, ",") + ` WHERE id = :id`

	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, student)
	} else {
		_, err = s.db.NamedExecContext(ctx, query, student)
	}
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return &error2.UniqueConstraintErr{Message: pqErr.Constraint}
		}
		log.Println("Student repo, update teacher error :", err)
		return err
	}
	return nil
}

func (s *studentRepo) DeleteStudentById(ctx context.Context, id string, tx *sqlx.Tx) error {
	query := `DELETE FROM students WHERE id = $1`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, id)
	} else {
		_, err = s.db.ExecContext(ctx, query, id)
	}
	if err != nil {
		log.Println("Student repo, delete student error :", err)
		return err
	}
	return nil
}

func (s *studentRepo) InsertStudent(ctx context.Context, student model.Student, tx *sqlx.Tx) error {
	query := `INSERT INTO students(id, school_year, major) VALUES (:id, :school_year, :major)`
	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, student)
	} else {
		_, err = s.db.NamedExecContext(ctx, query, student)
	}
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return &error2.UniqueConstraintErr{Message: pqErr.Constraint}
		}
		log.Println("Student repo, save student err :", err)
		return err
	}
	return nil
}

func NewStudentRepo(db *sqlx.DB) StudentRepo {
	return &studentRepo{db: db}
}
