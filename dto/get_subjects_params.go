package dto

type GetSubjectsParamDTO struct {
	PaginationParams
	Major string `json:"major"`
}
