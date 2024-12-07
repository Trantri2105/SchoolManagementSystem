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

type CourseRepo interface {
	CreateCourse(ctx context.Context, course model.Course, tx *sqlx.Tx) error
	GetCourseById(ctx context.Context, id string, tx *sqlx.Tx) (model.Course, error)
	GetCourseForUpdate(ctx context.Context, id string, tx *sqlx.Tx) (model.Course, error)
	UpdateCourse(ctx context.Context, course model.Course, tx *sqlx.Tx) error
	DeleteCourseById(ctx context.Context, id string, tx *sqlx.Tx) error
	InsertCourseRegistration(ctx context.Context, courseRegistration model.CourseRegistration, tx *sqlx.Tx) error
	DeleteCourseRegistration(ctx context.Context, courseId string, studentId string, tx *sqlx.Tx) error
	AddCourseSchedule(ctx context.Context, schedule model.CourseSchedule, tx *sqlx.Tx) error
	GetCourseSchedulesByCourseId(ctx context.Context, courseId string, tx *sqlx.Tx) ([]model.CourseSchedule, error)
	DeleteCourseScheduleById(ctx context.Context, id string, tx *sqlx.Tx) error
	GetCoursesByUserId(ctx context.Context, userId string, role string, semester int, academicYear string, tx *sqlx.Tx) ([]model.Course, error)
	DecreaseCourseSize(ctx context.Context, courseId string, quantity int, tx *sqlx.Tx) error
}

type courseRepo struct {
	db *sqlx.DB
}

func (c *courseRepo) DecreaseCourseSize(ctx context.Context, courseId string, quantity int, tx *sqlx.Tx) error {
	query := `UPDATE courses SET size = size - $1 WHERE id = $2`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, quantity, courseId)
	} else {
		_, err = c.db.ExecContext(ctx, query, quantity, courseId)
	}
	if err != nil {
		log.Println("Course repo, increase course size err: ", err)
		return err
	}
	return nil
}

func (c *courseRepo) GetCourseForUpdate(ctx context.Context, id string, tx *sqlx.Tx) (model.Course, error) {
	query := `SELECT id, teacher_id, subject_id, semester_number, academic_year, capacity, size, status FROM courses WHERE id = $1 FOR UPDATE`

	var row *sqlx.Row
	if tx != nil {
		row = tx.QueryRowxContext(ctx, query, id)
	} else {
		row = c.db.QueryRowxContext(ctx, query, id)
	}
	course := model.Course{}
	err := row.StructScan(&course)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return course, &error2.ResourceNotFoundErr{Resource: "course"}
		}
		log.Println("Course repo, get course for update err: ", err)
		return course, err
	}
	return course, nil
}

func (c *courseRepo) InsertCourseRegistration(ctx context.Context, courseRegistration model.CourseRegistration, tx *sqlx.Tx) error {
	query := `INSERT INTO course_registrations(course_id, student_id) VALUES (:course_id,:student_id)`

	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, courseRegistration)
	} else {
		_, err = c.db.NamedExecContext(ctx, query, courseRegistration)
	}
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return &error2.UniqueConstraintErr{Message: "student already registered"}
		}
		log.Println("Course repo, insert course registration err:", err)
		return err
	}
	return nil
}

func (c *courseRepo) UpdateCourse(ctx context.Context, course model.Course, tx *sqlx.Tx) error {
	var updateFields []string
	t := reflect.TypeOf(course)
	v := reflect.ValueOf(course)
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
	query := `UPDATE courses SET ` + strings.Join(updateFields, ",") + ` WHERE id = :id`
	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, course)
	} else {
		_, err = c.db.NamedExecContext(ctx, query, course)
	}

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return &error2.UniqueConstraintErr{Message: "student already registered"}
		}
		log.Println("Course repo, insert course registration err:", err)
		return err
	}
	return nil
}

func (c *courseRepo) GetCoursesByUserId(ctx context.Context, userId string, role string, semester int, academicYear string, tx *sqlx.Tx) ([]model.Course, error) {
	var query string
	if role == model.RoleStudent {
		query = `SELECT courses.id, users.name as teacher_name, subjects.name as subject_name, courses.semester_number, courses.academic_year, courses.capacity, courses.size, courses.status
 				FROM course_registrations
 				JOIN courses ON course_registrations.course_id = courses.id
 				JOIN users ON courses.teacher_id = users.id
 				JOIN subjects ON courses.subject_id = subjects.id
 				WHERE course_registrations.student_id = $1 AND courses.semester_number = $2 AND courses.academic_year = $3`
	} else {
		query = `SELECT courses.id, users.name as teacher_name, subjects.name as subject_name, courses.semester_number, courses.academic_year, courses.capacity, courses.size, courses.status
 				FROM courses
 				JOIN users ON courses.teacher_id = users.id
 				JOIN subjects ON courses.subject_id = subjects.id
 				WHERE courses.teacher_id = $1 AND courses.semester_number = $2 AND courses.academic_year = $3`
	}

	var err error
	var rows *sqlx.Rows
	if tx != nil {
		rows, err = tx.QueryxContext(ctx, query, userId, semester, academicYear)
	} else {
		rows, err = c.db.QueryxContext(ctx, query, userId, semester, academicYear)
	}
	if err != nil {
		log.Println("Course repo, ger course by user id err: ", err)
		return nil, err
	}
	var courses []model.Course
	for rows.Next() {
		var course model.Course
		err = rows.StructScan(&course)
		if err != nil {
			log.Println("Course repo, ger course by user id err: ", err)
			return nil, err
		}
		courses = append(courses, course)
	}
	if len(courses) == 0 {
		return nil, &error2.ResourceNotFoundErr{Resource: "courses"}
	}
	return courses, nil
}

func (c *courseRepo) GetCourseSchedulesByCourseId(ctx context.Context, courseId string, tx *sqlx.Tx) ([]model.CourseSchedule, error) {
	query := `SELECT id, course_id, room, start_time, end_time FROM course_schedules WHERE course_id = $1`

	var err error
	var rows *sqlx.Rows
	if tx != nil {
		rows, err = tx.QueryxContext(ctx, query, courseId)
	} else {
		rows, err = c.db.QueryxContext(ctx, query, courseId)
	}
	if err != nil {
		log.Println("Course repo, get course schedule err: ", err)
		return nil, err
	}
	var schedules []model.CourseSchedule
	for rows.Next() {
		var schedule model.CourseSchedule
		err = rows.StructScan(&schedule)
		if err != nil {
			log.Println("Course repo, get course schedule err: ", err)
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	if len(schedules) == 0 {
		return nil, &error2.ResourceNotFoundErr{Resource: "course schedules"}
	}
	return schedules, nil
}

func (c *courseRepo) DeleteCourseScheduleById(ctx context.Context, id string, tx *sqlx.Tx) error {
	query := `DELETE FROM course_schedules WHERE id = $1`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, id)
	} else {
		_, err = c.db.ExecContext(ctx, query, id)
	}
	if err != nil {
		log.Println("Course repo, delete course schedule err: ", err)
		return err
	}
	return nil
}

func (c *courseRepo) CreateCourse(ctx context.Context, course model.Course, tx *sqlx.Tx) error {
	query := `INSERT INTO courses(id, teacher_id, subject_id, semester_number, academic_year,capacity,size, status) 
			VALUES (:id, :teacher_id, :subject_id, :semester_number, :academic_year, :capacity, :size, :status)`

	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, course)
	} else {
		_, err = c.db.NamedExecContext(ctx, query, course)
	}
	if err != nil {
		log.Println("Course repo, create course err: ", err)
		return err
	}
	return nil
}

func (c *courseRepo) GetCourseById(ctx context.Context, id string, tx *sqlx.Tx) (model.Course, error) {
	query := `SELECT courses.id, users.name as teacher_name, subjects.name as subject_name, courses.semester_number, courses.academic_year, courses.capacity, courses.size, courses.status
			FROM courses
			JOIN users ON courses.teacher_id = users.id
			JOIN subjects ON courses.subject_id = subjects.id
			WHERE courses.id = $1`

	var row *sqlx.Row
	if tx != nil {
		row = tx.QueryRowxContext(ctx, query, id)
	} else {
		row = c.db.QueryRowxContext(ctx, query, id)
	}
	course := model.Course{}
	err := row.StructScan(&course)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return course, &error2.ResourceNotFoundErr{Resource: "Course"}
		}
		log.Println("Course repo, get course err: ", err)
		return course, err
	}
	return course, nil
}

func (c *courseRepo) DeleteCourseById(ctx context.Context, id string, tx *sqlx.Tx) error {
	query := `DELETE FROM courses WHERE id = $1`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, id)
	} else {
		_, err = c.db.ExecContext(ctx, query, id)
	}
	if err != nil {
		log.Println("Course repo, delete course err: ", err)
		return err
	}
	return nil
}

func (c *courseRepo) RegisterStudentToCourse(ctx context.Context, courseRegistration model.CourseRegistration) error {
	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println("Course repo, begin tx err: ", err)
		return err
	}
	defer func() {
		if err != nil {
			if !errors.Is(err, error2.CourseLimitExceededErr) {
				log.Println("Course repo, register student to cause error:", err)
			}
			if err = tx.Rollback(); err != nil {
				log.Println("Course repo, rollback transaction error:", err)
			}
		}
	}()
	var course model.Course
	query := `SELECT capacity, size FROM courses WHERE id = $1 FOR UPDATE`
	row := c.db.QueryRowxContext(ctx, query, courseRegistration.CourseId)
	err = row.Scan(&course.Capacity, &course.Size)
	if err != nil {
		log.Println("Course repo, register student to course err:", err)
		return err
	}
	if course.Capacity-course.Size == 0 {
		return error2.CourseLimitExceededErr
	}
	query = `INSERT INTO course_registrations(course_id, student_id) VALUES (:student_id,:course_id)`
	_, err = tx.NamedExecContext(ctx, query, courseRegistration)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return &error2.UniqueConstraintErr{Message: "student already registered"}
		}
		log.Println("Course repo, register student to course err:", err)
		return err
	}
	query = `UPDATE courses SET size = size + 1 WHERE id = $1`
	_, err = tx.ExecContext(ctx, query, courseRegistration.CourseId)
	if err != nil {
		log.Println("Course repo, register student to course err:", err)
		return err
	}
	return tx.Commit()
}

func (c *courseRepo) DeleteCourseRegistration(ctx context.Context, courseId string, studentId string, tx *sqlx.Tx) error {
	query := `DELETE FROM course_registrations WHERE course_id = $1 AND student_id = $2`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, courseId, studentId)
	} else {
		_, err = c.db.ExecContext(ctx, query, courseId, studentId)
	}
	if err != nil {
		log.Println("Course repo, delete student from course err:", err)
		return err
	}
	return nil
}

func (c *courseRepo) AddCourseSchedule(ctx context.Context, schedule model.CourseSchedule, tx *sqlx.Tx) error {
	query := `INSERT INTO course_schedules(course_id, room, start_time, end_time) VALUES (:course_id, :room, :start_time, :end_time)`
	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, schedule)
	} else {
		_, err = c.db.NamedExecContext(ctx, query, schedule)
	}
	if err != nil {
		log.Println("Course repo, add course schedule err:", err)
		return err
	}
	return nil
}

func NewCourseRepo(db *sqlx.DB) CourseRepo {
	return &courseRepo{
		db: db,
	}
}
