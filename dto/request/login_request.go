package request

type LoginRequest struct {
	Id       string `json:"id" validate:"required"`
	Password string `json:"password" validate:"required"`
}
