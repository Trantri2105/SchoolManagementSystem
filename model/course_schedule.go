package model

type CourseSchedule struct {
	Id        int    `db:"id"`
	CourseId  string `db:"course_id"`
	Room      string `db:"room"`
	StartTime int64  `db:"start_time"`
	EndTime   int64  `db:"end_time"`
}
