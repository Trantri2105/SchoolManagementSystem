package response

type GetStudentResponse struct {
	Id             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	DateOfBirth    string `json:"date_of_birth,omitempty"`
	Gender         string `json:"gender,omitempty"`
	Email          string `json:"email,omitempty"`
	IdentityNumber string `json:"identity_number,omitempty"`
	PhoneNumber    string `json:"phone_number,omitempty"`
	Address        string `json:"address,omitempty"`
	SchoolYear     string `json:"school_year,omitempty"`
	Major          string `json:"major,omitempty"`
}
