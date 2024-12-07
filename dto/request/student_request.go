package request

import (
	"SchoolManagement/model"
)

type StudentRequest struct {
	Id             string `json:"id" validate:"required"`
	Name           string `json:"name" validate:"required"`
	DateOfBirth    string `json:"date_of_birth" validate:"required"`
	Gender         string `json:"gender" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	IdentityNumber string `json:"identity_number" validate:"required"`
	PhoneNumber    string `json:"phone_number" validate:"required,number"`
	Address        string `json:"address" validate:"required"`
	Password       string `json:"password" validate:"required"`
	SchoolYear     string `json:"school_year" validate:"required"`
	Major          string `json:"major" validate:"required"`
}

func (s *StudentRequest) ToStudent() model.Student {
	return model.Student{
		User: model.User{
			Id:             s.Id,
			Name:           s.Name,
			DateOfBirth:    s.DateOfBirth,
			Gender:         s.Gender,
			Email:          s.Email,
			IdentityNumber: s.IdentityNumber,
			PhoneNumber:    s.PhoneNumber,
			Address:        s.Address,
			Password:       s.Password,
			Role:           string(model.RoleStudent),
		},
		SchoolYear: s.SchoolYear,
		Major:      s.Major,
	}
}
