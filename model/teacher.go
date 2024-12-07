package model

type Teacher struct {
	User
	AcademicQualification string `db:"academic_qualification"`
	Department            string `db:"department"`
}
