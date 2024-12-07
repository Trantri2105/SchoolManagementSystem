package dto

import "SchoolManagement/model"

type SubjectDto struct {
	Id             string `json:"id" validate:"required"`
	Name           string `json:"name" validate:"required"`
	NumberOfCredit int    `json:"number_of_credit" validate:"required"`
	Major          string `json:"major" validate:"required"`
}

func (s *SubjectDto) ToSubject() model.Subject {
	return model.Subject{
		Id:             s.Id,
		Name:           s.Name,
		NumberOfCredit: s.NumberOfCredit,
		Major:          s.Major,
	}
}
