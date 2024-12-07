package model

type Student struct {
	User
	SchoolYear string `db:"school_year"`
	Major      string `db:"major"`
}
