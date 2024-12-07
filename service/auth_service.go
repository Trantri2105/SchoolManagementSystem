package service

import (
	"SchoolManagement/dto/response"
	error2 "SchoolManagement/error"
	"SchoolManagement/repo/postgres"
	"SchoolManagement/utils"
	"context"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, id string, password string) (response.LoginResponse, error)
}

type authService struct {
	userRepo postgres.UserRepo
	jwtUtils utils.JwtUtils
}

func (a *authService) Login(ctx context.Context, id string, password string) (response.LoginResponse, error) {
	user, err := a.userRepo.GetUserById(ctx, id, nil)
	if err != nil {
		return response.LoginResponse{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return response.LoginResponse{}, error2.WrongPasswordErr
	}
	token, expireTime, err := a.jwtUtils.CreateToken(id, user.Role)
	if err != nil {
		return response.LoginResponse{}, err
	}
	return response.LoginResponse{
		AccessToken: token,
		ExpiresIn:   expireTime,
		Role:        user.Role,
	}, nil
}

func NewAuthService(userRepo postgres.UserRepo, jwtUtils utils.JwtUtils) AuthService {
	return &authService{userRepo: userRepo, jwtUtils: jwtUtils}
}
