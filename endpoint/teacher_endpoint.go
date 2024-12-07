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

type TeacherEndpoint interface {
	RegisterTeacherEndpoint() endpoint.Endpoint
	UpdateTeacherEndpoint() endpoint.Endpoint
	GetTeacherByIdEndpoint() endpoint.Endpoint
	DeleteTeacherByIdEndpoint() endpoint.Endpoint
}

type teacherEndpoint struct {
	teacherService service.TeacherService
}

func (t *teacherEndpoint) RegisterTeacherEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(request2.TeacherRequest)
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			return nil, err
		}
		err := t.teacherService.CreateTeacher(ctx, req.ToTeacher())
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Teacher register successfully"}, nil
	}
}

func (t *teacherEndpoint) UpdateTeacherEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(request2.TeacherRequest)
		fields := []string{"Id"}
		if !reflect.ValueOf(req.Email).IsZero() {
			fields = append(fields, "Email")
		}
		validate := validator.New()
		if err := validate.StructPartial(req, fields...); err != nil {
			return nil, err
		}
		err := t.teacherService.UpdateTeacher(ctx, req.ToTeacher())
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Teacher update successfully"}, nil
	}
}

func (t *teacherEndpoint) GetTeacherByIdEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		teacher, err := t.teacherService.GetTeacherById(ctx, req)
		if err != nil {
			return nil, err
		}
		return response.GetTeacherResponse{
			Id:                    teacher.Id,
			Name:                  teacher.Name,
			DateOfBirth:           teacher.DateOfBirth,
			Gender:                teacher.Gender,
			Email:                 teacher.Email,
			IdentityNumber:        teacher.IdentityNumber,
			PhoneNumber:           teacher.PhoneNumber,
			Address:               teacher.Address,
			Role:                  teacher.Role,
			AcademicQualification: teacher.AcademicQualification,
			Department:            teacher.Department,
		}, nil
	}
}

func (t *teacherEndpoint) DeleteTeacherByIdEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		err := t.teacherService.DeleteTeacherById(ctx, req)
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Teacher delete successfully"}, nil
	}
}

func NewTeacherEndpoint(teacherService service.TeacherService) TeacherEndpoint {
	return &teacherEndpoint{
		teacherService: teacherService,
	}
}
