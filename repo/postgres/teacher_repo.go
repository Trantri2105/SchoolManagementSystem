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

type TeacherRepo interface {
	InsertTeacher(ctx context.Context, teacher model.Teacher, tx *sqlx.Tx) error
	GetTeacherById(ctx context.Context, id string, tx *sqlx.Tx) (model.Teacher, error)
	DeleteTeacherById(ctx context.Context, id string, tx *sqlx.Tx) error
	UpdateTeacher(ctx context.Context, teacher model.Teacher, tx *sqlx.Tx) error
}

type teacherRepo struct {
	db *sqlx.DB
}

func (r *teacherRepo) InsertTeacher(ctx context.Context, teacher model.Teacher, tx *sqlx.Tx) error {
	query := `INSERT INTO teachers(id, academic_qualification, department) VALUES (:id, :academic_qualification, :department)`

	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, teacher)
	} else {
		_, err = r.db.NamedExecContext(ctx, query, teacher)
	}

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return &error2.UniqueConstraintErr{Message: pqErr.Constraint}
		}
		log.Println("Teacher repo, save teacher err :", err)
		return err
	}

	return nil
}

func (r *teacherRepo) GetTeacherById(ctx context.Context, id string, tx *sqlx.Tx) (model.Teacher, error) {
	query := `SELECT users.id, name, date_of_birth, gender, email, identity_number, phone_number, address, password, role, academic_qualification, department
			FROM users
			JOIN teachers ON users.id = teachers.id
			WHERE users.id = $1`

	var row *sqlx.Row
	if tx != nil {
		row = tx.QueryRowxContext(ctx, query, id)
	} else {
		row = r.db.QueryRowxContext(ctx, query, id)
	}
	var teacher model.Teacher
	err := row.StructScan(&teacher)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return teacher, &error2.ResourceNotFoundErr{Resource: "Teacher"}
		}
		log.Println("Teacher repo, get teacher err :", err)
		return teacher, err
	}
	return teacher, nil
}

func (r *teacherRepo) DeleteTeacherById(ctx context.Context, id string, tx *sqlx.Tx) error {
	query := `DELETE FROM teachers WHERE id = $1`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, id)
	} else {
		_, err = r.db.ExecContext(ctx, query, id)
	}
	if err != nil {
		log.Println("Teacher repo, delete teacher err :", err)
		return err
	}
	return nil
}

func (r *teacherRepo) UpdateTeacher(ctx context.Context, teacher model.Teacher, tx *sqlx.Tx) error {
	var updateFields []string
	t := reflect.TypeOf(teacher)
	v := reflect.ValueOf(teacher)
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
	query := `UPDATE teachers SET ` + strings.Join(updateFields, ",") + ` WHERE id = :id`
	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, teacher)
	} else {
		_, err = r.db.NamedExecContext(ctx, query, teacher)
	}
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return &error2.UniqueConstraintErr{Message: pqErr.Constraint}
		}
		log.Println("Teacher repo, update teacher error :", err)
		return err
	}
	return nil
}

func NewTeacherRepo(db *sqlx.DB) TeacherRepo {
	return &teacherRepo{db: db}
}
