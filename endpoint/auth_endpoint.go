package endpoint

import (
	request2 "SchoolManagement/dto/request"
	"SchoolManagement/service"
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-playground/validator/v10"
)

type AuthEndpoint interface {
	Login() endpoint.Endpoint
}

type authEndpoint struct {
	authService service.AuthService
}

func (a *authEndpoint) Login() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(request2.LoginRequest)
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			return nil, err
		}
		res, err := a.authService.Login(ctx, req.Id, req.Password)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func NewAuthEndpoint(authService service.AuthService) AuthEndpoint {
	return &authEndpoint{authService: authService}
}
