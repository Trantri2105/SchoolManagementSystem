package request

import "SchoolManagement/model"

type CourseRequest struct {
	Id             string `json:"id" validate:"required"`
	TeacherId      string `json:"teacher_id" validate:"required"`
	SubjectId      string `json:"subject_id" validate:"required"`
	SemesterNumber int    `json:"semester_number" validate:"required"`
	AcademicYear   string `json:"academic_year" validate:"required"`
	Capacity       int    `json:"capacity" validate:"required"`
	Status         string `json:"status"`
}

func (req *CourseRequest) ToCourse() model.Course {
	return model.Course{
		Id:             req.Id,
		TeacherId:      req.TeacherId,
		SubjectId:      req.SubjectId,
		SemesterNumber: req.SemesterNumber,
		AcademicYear:   req.AcademicYear,
		Capacity:       req.Capacity,
		Status:         req.Status,
		Size:           0,
	}
}
