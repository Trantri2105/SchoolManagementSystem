package endpoint

import (
	request2 "SchoolManagement/dto/request"
	"SchoolManagement/dto/response"
	"SchoolManagement/service"
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type StudentEndpoint interface {
	RegisterStudentEndpoint() endpoint.Endpoint
	UpdateStudentEndpoint() endpoint.Endpoint
	DeleteStudentByIdEndpoint() endpoint.Endpoint
	GetStudentByIdEndpoint() endpoint.Endpoint
}

type studentEndpoint struct {
	studentService service.StudentService
}

func (s *studentEndpoint) RegisterStudentEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(request2.StudentRequest)
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			return nil, err
		}
		err := s.studentService.CreateStudent(ctx, req.ToStudent())
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Student register successfully"}, nil
	}
}

func (s *studentEndpoint) UpdateStudentEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(request2.StudentRequest)
		fields := []string{"Id"}
		if !reflect.ValueOf(req.Email).IsZero() {
			fields = append(fields, "Email")
		}
		validate := validator.New()
		if err := validate.StructPartial(req, fields...); err != nil {
			return nil, err
		}
		err := s.studentService.UpdateStudent(ctx, req.ToStudent())
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Student updated successfully"}, nil
	}
}

func (s *studentEndpoint) DeleteStudentByIdEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		err := s.studentService.DeleteStudentById(ctx, req)
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Student deleted successfully"}, nil
	}
}

func (s *studentEndpoint) GetStudentByIdEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		student, err := s.studentService.GetStudentById(ctx, req)
		if err != nil {
			return nil, err
		}
		return response.GetStudentResponse{
			Id:             student.Id,
			Name:           student.Name,
			DateOfBirth:    student.DateOfBirth,
			Gender:         student.Gender,
			Email:          student.Email,
			IdentityNumber: student.IdentityNumber,
			PhoneNumber:    student.PhoneNumber,
			Address:        student.Address,
			SchoolYear:     student.SchoolYear,
			Major:          student.Major,
		}, nil
	}
}

func NewStudentEndpoint(studentService service.StudentService) StudentEndpoint {
	return &studentEndpoint{
		studentService: studentService,
	}
}
