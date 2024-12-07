package response

type GetTeacherResponse struct {
	Id                    string `json:"id,omitempty"`
	Name                  string `json:"name,omitempty"`
	DateOfBirth           string `json:"date_of_birth,omitempty"`
	Gender                string `json:"gender,omitempty"`
	Email                 string `json:"email,omitempty"`
	IdentityNumber        string `json:"identity_number,omitempty"`
	PhoneNumber           string `json:"phone_number,omitempty"`
	Address               string `json:"address,omitempty"`
	Role                  string `json:"role,omitempty"`
	AcademicQualification string `json:"academic_qualification,omitempty"`
	Department            string `json:"department,omitempty"`
}
