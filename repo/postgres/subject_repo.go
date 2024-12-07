package postgres

import (
	"SchoolManagement/dto"
	error2 "SchoolManagement/error"
	"SchoolManagement/model"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
	"reflect"
	"strings"
)

type SubjectRepo interface {
	InsertSubject(ctx context.Context, subject model.Subject, tx *sqlx.Tx) error
	UpdateSubject(ctx context.Context, subject model.Subject, tx *sqlx.Tx) error
	DeleteSubjectById(ctx context.Context, id string, tx *sqlx.Tx) error
	GetSubjectById(ctx context.Context, id string, tx *sqlx.Tx) (model.Subject, error)
	GetSubjectList(ctx context.Context, params dto.GetSubjectsParamDTO, tx *sqlx.Tx) ([]model.Subject, error)
}

type subjectRepo struct {
	db *sqlx.DB
}

func (s *subjectRepo) InsertSubject(ctx context.Context, subject model.Subject, tx *sqlx.Tx) error {
	query := `INSERT INTO subjects(id, name, number_of_credit, major) VALUES (:id, :name, :number_of_credit, :major)`

	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, subject)
	} else {
		_, err = s.db.NamedExecContext(ctx, query, subject)
	}

	if err != nil {
		log.Println("Subject repo, insert subject err: ", err)
		return err
	}
	return nil
}

func (s *subjectRepo) UpdateSubject(ctx context.Context, subject model.Subject, tx *sqlx.Tx) error {
	var updateFields []string
	t := reflect.TypeOf(subject)
	v := reflect.ValueOf(subject)
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
	query := `UPDATE subjects SET ` + strings.Join(updateFields, ",") + ` WHERE id = :id`

	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, subject)
	} else {
		_, err = s.db.NamedExecContext(ctx, query, subject)
	}
	if err != nil {
		log.Println("Subject repo, update subject err: ", err)
		return err
	}
	return nil
}

func (s *subjectRepo) DeleteSubjectById(ctx context.Context, id string, tx *sqlx.Tx) error {
	query := `DELETE FROM subjects WHERE id = $1`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, id)
	} else {
		_, err = s.db.ExecContext(ctx, query, id)
	}
	if err != nil {
		log.Println("Subject repo, delete subject err: ", err)
		return err
	}
	return nil
}

func (s *subjectRepo) GetSubjectById(ctx context.Context, id string, tx *sqlx.Tx) (model.Subject, error) {
	query := `SELECT id, name, number_of_credit, major FROM subjects WHERE id = $1`
	var row *sqlx.Row
	if tx != nil {
		row = tx.QueryRowxContext(ctx, query, id)
	} else {
		row = s.db.QueryRowxContext(ctx, query, id)
	}
	var subject model.Subject
	err := row.StructScan(&subject)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return subject, &error2.ResourceNotFoundErr{Resource: "Subject"}
		}
		log.Println("Subject repo, get subject err: ", err)
		return subject, err
	}
	return subject, nil
}

func (s *subjectRepo) GetSubjectList(ctx context.Context, params dto.GetSubjectsParamDTO, tx *sqlx.Tx) ([]model.Subject, error) {
	var rows *sqlx.Rows
	var err error
	if params.Major != "" {
		query := `SELECT id, name, number_of_credit, major FROM subjects WHERE major = $1 ORDER BY id LIMIT $2 OFFSET $3`
		if tx != nil {
			rows, err = tx.QueryxContext(ctx, query, params.Major, params.Limit, params.Offset)
		} else {
			rows, err = s.db.QueryxContext(ctx, query, params.Major, params.Limit, params.Offset)
		}
	} else {
		query := `SELECT id, name, number_of_credit, major FROM subjects ORDER BY id LIMIT $1 OFFSET $2`
		if tx != nil {
			rows, err = tx.QueryxContext(ctx, query, params.Limit, params.Offset)
		} else {
			rows, err = s.db.QueryxContext(ctx, query, params.Limit, params.Offset)
		}
	}
	if err != nil {
		log.Println("Subject repo, get subjects err: ", err)
		return nil, err
	}
	var subjects []model.Subject
	defer rows.Close()
	for rows.Next() {
		var subject model.Subject
		err = rows.StructScan(&subject)
		if err != nil {
			log.Println("Subject repo, get subjects err: ", err)
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	if len(subjects) == 0 {
		return nil, &error2.ResourceNotFoundErr{Resource: "Subject"}
	}
	return subjects, nil
}

func NewSubjectRepo(db *sqlx.DB) SubjectRepo {
	return &subjectRepo{db: db}
}
