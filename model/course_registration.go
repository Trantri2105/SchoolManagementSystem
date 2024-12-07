package model

type CourseRegistration struct {
	Id        int    `db:"id"`
	CourseId  string `db:"course_id"`
	StudentId string `db:"student_id"`
}
