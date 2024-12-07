package response

type CourseResponse struct {
	Id             string `json:"id"`
	TeacherName    string `json:"teacher_name"`
	SubjectName    string `json:"subject_name"`
	SemesterNumber int    `json:"semester_number"`
	AcademicYear   string `json:"academic_year"`
	Capacity       int    `json:"capacity"`
	Size           int    `json:"size"`
	Status         string `json:"status"`
}
