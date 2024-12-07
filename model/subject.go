package model

type Subject struct {
	Id             string `db:"id"`
	Name           string `db:"name"`
	NumberOfCredit int    `db:"number_of_credit"`
	Major          string `db:"major"`
}
