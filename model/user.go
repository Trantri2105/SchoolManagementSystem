package model

const (
	RoleAdmin   string = "Admin"
	RoleTeacher string = "Teacher"
	RoleStudent string = "Student"
)

type User struct {
	Id             string `db:"id"`
	Name           string `db:"name"`
	DateOfBirth    string `db:"date_of_birth"`
	Gender         string `db:"gender"`
	Email          string `db:"email"`
	IdentityNumber string `db:"identity_number"`
	PhoneNumber    string `db:"phone_number"`
	Address        string `db:"address"`
	Password       string `db:"password"`
	Role           string `db:"role"`
}
