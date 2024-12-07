package model

const (
	CourseStatusInitial  string = "Initial"
	CourseStatusRegister string = "Register"
	CourseStatusOngoing  string = "Ongoing"
	CourseStatusComplete string = "Complete"
)

type Course struct {
	Id             string `db:"id"`
	TeacherId      string `db:"teacher_id"`
	TeacherName    string `db:"teacher_name"`
	SubjectId      string `db:"subject_id"`
	SubjectName    string `db:"subject_name"`
	SemesterNumber int    `db:"semester_number"`
	AcademicYear   string `db:"academic_year"`
	Capacity       int    `db:"capacity"`
	Size           int    `db:"size"`
	Status         string `db:"status"`
}
