package request

import "SchoolManagement/model"

type CourseRegistrationRequest struct {
	CourseId  string `json:"course_id" validate:"required"`
	StudentId string `json:"student_id" validate:"required"`
}

func (req *CourseRegistrationRequest) ToCourseRegistration() model.CourseRegistration {
	return model.CourseRegistration{
		CourseId:  req.CourseId,
		StudentId: req.StudentId,
	}
}
