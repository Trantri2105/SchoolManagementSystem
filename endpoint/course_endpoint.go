package endpoint

import (
	"SchoolManagement/dto"
	request2 "SchoolManagement/dto/request"
	"SchoolManagement/dto/response"
	"SchoolManagement/service"
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-playground/validator/v10"
)

type CourseEndpoint interface {
	CreateCourse() endpoint.Endpoint
	GetCourseById() endpoint.Endpoint
	UpdateCourse() endpoint.Endpoint
	DeleteCourseById() endpoint.Endpoint
	RegisterStudentToCourse() endpoint.Endpoint
	UnregisterStudentFromCourse() endpoint.Endpoint
	AddCourseSchedule() endpoint.Endpoint
	GetCourseSchedulesByCourseId() endpoint.Endpoint
	DeleteCourseScheduleById() endpoint.Endpoint
	GetCoursesByUserId() endpoint.Endpoint
}

type courseEndpoint struct {
	courseService service.CourseService
}

func (c *courseEndpoint) CreateCourse() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(request2.CourseRequest)
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			return nil, err
		}
		err := c.courseService.CreateCourse(ctx, req.ToCourse())
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Course created"}, nil
	}
}

func (c *courseEndpoint) GetCourseById() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		course, err := c.courseService.GetCourseById(ctx, req)
		if err != nil {
			return nil, err
		}
		return response.CourseResponse{
			Id:             course.Id,
			TeacherName:    course.TeacherName,
			SubjectName:    course.SubjectName,
			SemesterNumber: course.SemesterNumber,
			AcademicYear:   course.AcademicYear,
			Capacity:       course.Capacity,
			Size:           course.Size,
			Status:         course.Status,
		}, nil
	}
}

func (c *courseEndpoint) UpdateCourse() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(request2.CourseRequest)
		err := c.courseService.UpdateCourse(ctx, req.ToCourse())
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Course updated"}, nil
	}
}

func (c *courseEndpoint) DeleteCourseById() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		err := c.courseService.DeleteCourseById(ctx, req)
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Course deleted"}, nil
	}
}

func (c *courseEndpoint) RegisterStudentToCourse() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(request2.CourseRegistrationRequest)
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			return nil, err
		}
		err := c.courseService.RegisterStudentToCourse(ctx, req.ToCourseRegistration())
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Registered"}, nil
	}
}

func (c *courseEndpoint) UnregisterStudentFromCourse() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(request2.CourseRegistrationRequest)
		err := c.courseService.UnregisterStudentFromCourse(ctx, req.CourseId, req.StudentId)
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Unregistered"}, nil
	}
}

func (c *courseEndpoint) AddCourseSchedule() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(request2.CourseScheduleRequest)
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			return nil, err
		}
		err := c.courseService.AddCourseSchedule(ctx, req.ToCourseSchedule())
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Add course schedule successfully"}, nil
	}
}

func (c *courseEndpoint) GetCourseSchedulesByCourseId() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		schedules, err := c.courseService.GetCourseSchedulesByCourseId(ctx, req)
		if err != nil {
			return nil, err
		}
		var res []response.CourseScheduleResponse
		for _, schedule := range schedules {
			res = append(res, response.CourseScheduleResponse{
				Id:        schedule.Id,
				CourseId:  schedule.CourseId,
				Room:      schedule.Room,
				StartTime: schedule.StartTime,
				EndTime:   schedule.EndTime,
			})
		}
		return res, nil
	}
}

func (c *courseEndpoint) DeleteCourseScheduleById() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		err := c.courseService.DeleteCourseScheduleById(ctx, req)
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Course schedule deleted"}, nil
	}
}

func (c *courseEndpoint) GetCoursesByUserId() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.GetCoursesParams)
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			return nil, err
		}
		courses, err := c.courseService.GetCoursesByUserId(ctx, req.UserId, req.Semester, req.AcademicYear)
		if err != nil {
			return nil, err
		}
		var res []response.CourseResponse
		for _, course := range courses {
			res = append(res, response.CourseResponse{
				Id:             course.Id,
				TeacherName:    course.TeacherName,
				SubjectName:    course.SubjectName,
				SemesterNumber: course.SemesterNumber,
				AcademicYear:   course.AcademicYear,
				Capacity:       course.Capacity,
				Size:           course.Size,
				Status:         course.Status,
			})
		}
		return res, nil
	}
}

func NewCourseEndpoint(courseService service.CourseService) CourseEndpoint {
	return &courseEndpoint{
		courseService: courseService,
	}
}
