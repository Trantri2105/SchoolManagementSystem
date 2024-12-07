package dto

type GetCoursesParams struct {
	UserId       string `json:"user_id" validate:"required"`
	Semester     int    `json:"semester" validate:"required"`
	AcademicYear string `json:"academic_year" validate:"required"`
}
