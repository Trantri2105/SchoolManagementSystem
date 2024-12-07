package service

import (
	"SchoolManagement/dto"
	error2 "SchoolManagement/error"
	"SchoolManagement/middleware"
	"SchoolManagement/model"
	"SchoolManagement/repo/postgres"
	"context"
)

type SubjectService interface {
	CreateSubject(ctx context.Context, subject model.Subject) error
	UpdateSubject(ctx context.Context, subject model.Subject) error
	DeleteSubjectById(ctx context.Context, id string) error
	GetSubjectById(ctx context.Context, id string) (model.Subject, error)
	GetSubjectList(ctx context.Context, params dto.GetSubjectsParamDTO) ([]model.Subject, error)
}

type subjectService struct {
	subjectRepo    postgres.SubjectRepo
	authMiddleware middleware.AuthMiddleware
}

func (s *subjectService) CreateSubject(ctx context.Context, subject model.Subject) error {
	err := s.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin)
	if err != nil {
		return &error2.UnauthorizedErr{Message: "Require admin role to create subject"}
	}
	return s.subjectRepo.InsertSubject(ctx, subject, nil)
}

func (s *subjectService) UpdateSubject(ctx context.Context, subject model.Subject) error {
	err := s.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin)
	if err != nil {
		return &error2.UnauthorizedErr{Message: "Require admin role to update subject"}
	}
	return s.subjectRepo.UpdateSubject(ctx, subject, nil)
}

func (s *subjectService) DeleteSubjectById(ctx context.Context, id string) error {
	err := s.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin)
	if err != nil {
		return &error2.UnauthorizedErr{Message: "Require admin role to delete subject"}
	}
	return s.subjectRepo.DeleteSubjectById(ctx, id, nil)
}

func (s *subjectService) GetSubjectById(ctx context.Context, id string) (model.Subject, error) {
	return s.subjectRepo.GetSubjectById(ctx, id, nil)
}

func (s *subjectService) GetSubjectList(ctx context.Context, params dto.GetSubjectsParamDTO) ([]model.Subject, error) {
	return s.subjectRepo.GetSubjectList(ctx, params, nil)
}

func NewSubjectService(subjectRepo postgres.SubjectRepo, authMiddleware middleware.AuthMiddleware) SubjectService {
	return &subjectService{subjectRepo: subjectRepo, authMiddleware: authMiddleware}
}
