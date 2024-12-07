package request

import "SchoolManagement/model"

type TeacherRequest struct {
	Id                    string `json:"id" validate:"required"`
	Name                  string `json:"name" validate:"required"`
	DateOfBirth           string `json:"date_of_birth" validate:"required"`
	Gender                string `json:"gender" validate:"required"`
	Email                 string `json:"email" validate:"required,email"`
	IdentityNumber        string `json:"identity_number" validate:"required"`
	PhoneNumber           string `json:"phone_number" validate:"required,number"`
	Address               string `json:"address" validate:"required"`
	Password              string `json:"password" validate:"required"`
	Role                  string `json:"role" validate:"required"`
	AcademicQualification string `json:"academic_qualification" validate:"required"`
	Department            string `json:"department" validate:"required"`
}

func (t *TeacherRequest) ToTeacher() model.Teacher {
	return model.Teacher{
		User: model.User{
			Id:             t.Id,
			Name:           t.Name,
			DateOfBirth:    t.DateOfBirth,
			Gender:         t.Gender,
			Email:          t.Email,
			IdentityNumber: t.IdentityNumber,
			PhoneNumber:    t.PhoneNumber,
			Address:        t.Address,
			Password:       t.Password,
			Role:           t.Role,
		},
		AcademicQualification: t.AcademicQualification,
		Department:            t.Department,
	}
}
