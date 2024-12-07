package endpoint

import (
	"SchoolManagement/dto"
	"SchoolManagement/dto/response"
	"SchoolManagement/service"
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-playground/validator/v10"
)

type SubjectEndpoint interface {
	CreateSubjectEndpoint() endpoint.Endpoint
	UpdateSubjectEndpoint() endpoint.Endpoint
	DeleteSubjectByIdEndpoint() endpoint.Endpoint
	GetSubjectByIdEndpoint() endpoint.Endpoint
	GetSubjectListEndpoint() endpoint.Endpoint
}

type subjectEndpoint struct {
	subjectService service.SubjectService
}

func (s *subjectEndpoint) CreateSubjectEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SubjectDto)
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			return nil, err
		}
		err := s.subjectService.CreateSubject(ctx, req.ToSubject())
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Subject created successfully"}, nil
	}
}

func (s *subjectEndpoint) UpdateSubjectEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SubjectDto)
		validate := validator.New()
		if err := validate.StructPartial(req, "Id"); err != nil {
			return nil, err
		}
		err := s.subjectService.UpdateSubject(ctx, req.ToSubject())
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Subject updated successfully"}, nil
	}
}

func (s *subjectEndpoint) DeleteSubjectByIdEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		err := s.subjectService.DeleteSubjectById(ctx, req)
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Subject deleted successfully"}, nil
	}
}

func (s *subjectEndpoint) GetSubjectByIdEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		subject, err := s.subjectService.GetSubjectById(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.SubjectDto{
			Id:             subject.Id,
			Name:           subject.Name,
			NumberOfCredit: subject.NumberOfCredit,
			Major:          subject.Major,
		}, nil
	}
}

func (s *subjectEndpoint) GetSubjectListEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.GetSubjectsParamDTO)
		if req.Limit <= 0 || req.Limit > 20 {
			req.Limit = 20
		}
		if req.Offset < 0 {
			req.Offset = 0
		}
		subjects, err := s.subjectService.GetSubjectList(ctx, req)
		if err != nil {
			return nil, err
		}
		var subjectList []dto.SubjectDto
		for _, subject := range subjects {
			subjectList = append(subjectList, dto.SubjectDto{
				Id:             subject.Id,
				Name:           subject.Name,
				NumberOfCredit: subject.NumberOfCredit,
				Major:          subject.Major,
			})
		}
		return subjectList, nil
	}
}

func NewSubjectEndpoint(subjectService service.SubjectService) SubjectEndpoint {
	return &subjectEndpoint{
		subjectService: subjectService,
	}
}
