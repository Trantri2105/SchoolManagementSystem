package service

import (
	error2 "SchoolManagement/error"
	"SchoolManagement/middleware"
	"SchoolManagement/model"
	"SchoolManagement/repo"
	"SchoolManagement/repo/postgres"
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
)

type CourseService interface {
	CreateCourse(ctx context.Context, course model.Course) error
	GetCourseById(ctx context.Context, id string) (model.Course, error)
	UpdateCourse(ctx context.Context, course model.Course) error
	DeleteCourseById(ctx context.Context, id string) error
	RegisterStudentToCourse(ctx context.Context, courseRegistration model.CourseRegistration) error
	UnregisterStudentFromCourse(ctx context.Context, courseId string, studentId string) error
	AddCourseSchedule(ctx context.Context, schedule model.CourseSchedule) error
	GetCourseSchedulesByCourseId(ctx context.Context, courseId string) ([]model.CourseSchedule, error)
	DeleteCourseScheduleById(ctx context.Context, id string) error
	GetCoursesByUserId(ctx context.Context, userId string, semester int, academicYear string) ([]model.Course, error)
}

type courseService struct {
	courseRepo         postgres.CourseRepo
	transactionManager repo.TransactionManager
	authMiddleware     middleware.AuthMiddleware
	userRepo           postgres.UserRepo
}

func (c *courseService) CreateCourse(ctx context.Context, course model.Course) error {
	err := c.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin)
	if err != nil {
		return &error2.UnauthorizedErr{Message: "Required admin role to create course"}
	}
	course.Status = model.CourseStatusInitial
	return c.courseRepo.CreateCourse(ctx, course, nil)
}

func (c *courseService) GetCourseById(ctx context.Context, id string) (model.Course, error) {
	return c.courseRepo.GetCourseById(ctx, id, nil)
}

func (c *courseService) UpdateCourse(ctx context.Context, course model.Course) error {
	err := c.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin)
	if err != nil {
		return &error2.UnauthorizedErr{Message: "Required admin role to update course"}
	}
	if course.Status != "" && course.Status != model.CourseStatusInitial && course.Status != model.CourseStatusRegister && course.Status != model.CourseStatusOngoing && course.Status != model.CourseStatusComplete {
		return &error2.InvalidInputErr{Message: "Course status must be Initial, Register, Ongoing or Complete"}
	}
	return c.courseRepo.UpdateCourse(ctx, course, nil)
}

func (c *courseService) DeleteCourseById(ctx context.Context, id string) error {
	err := c.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin)
	if err != nil {
		return &error2.UnauthorizedErr{Message: "Required admin role to delete course"}
	}
	return c.courseRepo.DeleteCourseById(ctx, id, nil)
}

func (c *courseService) RegisterStudentToCourse(ctx context.Context, courseRegistration model.CourseRegistration) error {
	return c.transactionManager.ExecTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		err := c.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin, model.RoleStudent)
		if err != nil {
			return &error2.UnauthorizedErr{Message: "Required admin role or student role to register student"}
		}
		claims := ctx.Value(middleware.JWTClaimsContextKey).(jwt.MapClaims)
		if claims["role"].(string) == model.RoleStudent {
			if claims["userId"].(string) != courseRegistration.StudentId {
				return &error2.UnauthorizedErr{Message: "Unauthorized"}
			}
		}
		course, err := c.courseRepo.GetCourseForUpdate(ctx, courseRegistration.CourseId, tx)
		if err != nil {
			return err
		}
		if course.Status != model.CourseStatusRegister {
			return error2.CourseRegisterTimoutErr
		}
		if course.Capacity-course.Size == 0 {
			return error2.CourseLimitExceededErr
		}
		err = c.courseRepo.InsertCourseRegistration(ctx, courseRegistration, tx)
		if err != nil {
			return err
		}
		course.Size += 1
		err = c.courseRepo.UpdateCourse(ctx, course, tx)
		if err != nil {
			return err
		}
		return nil
	})
}

func (c *courseService) UnregisterStudentFromCourse(ctx context.Context, courseId string, studentId string) error {
	err := c.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin, model.RoleStudent)
	if err != nil {
		return &error2.UnauthorizedErr{Message: "Required admin or student role to delete student from course"}
	}
	claims := ctx.Value(middleware.JWTClaimsContextKey).(jwt.MapClaims)
	if claims["role"].(string) == model.RoleStudent {
		if claims["userId"].(string) != studentId {
			return &error2.UnauthorizedErr{Message: "Unauthorized"}
		}
	}
	return c.transactionManager.ExecTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		course, err := c.courseRepo.GetCourseById(ctx, courseId, tx)
		if err != nil {
			return err
		}
		if course.Status != model.CourseStatusRegister {
			return error2.CourseRegisterTimoutErr
		}
		err = c.courseRepo.DeleteCourseRegistration(ctx, courseId, studentId, tx)
		if err != nil {
			return err
		}
		err = c.courseRepo.DecreaseCourseSize(ctx, courseId, 1, tx)
		return err
	})
}

func (c *courseService) AddCourseSchedule(ctx context.Context, schedule model.CourseSchedule) error {
	err := c.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin)
	if err != nil {
		return &error2.UnauthorizedErr{Message: "Required admin role to add course schedule"}
	}
	return c.courseRepo.AddCourseSchedule(ctx, schedule, nil)
}

func (c *courseService) GetCourseSchedulesByCourseId(ctx context.Context, courseId string) ([]model.CourseSchedule, error) {
	return c.courseRepo.GetCourseSchedulesByCourseId(ctx, courseId, nil)
}

func (c *courseService) DeleteCourseScheduleById(ctx context.Context, id string) error {
	err := c.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin)
	if err != nil {
		return &error2.UnauthorizedErr{Message: "Required admin role to delete course schedule"}
	}
	return c.courseRepo.DeleteCourseScheduleById(ctx, id, nil)
}

func (c *courseService) GetCoursesByUserId(ctx context.Context, userId string, semester int, academicYear string) ([]model.Course, error) {
	claims := ctx.Value(middleware.JWTClaimsContextKey).(jwt.MapClaims)
	role := claims["role"].(string)
	if role != model.RoleAdmin {
		if claims["userId"].(string) != userId {
			return nil, &error2.UnauthorizedErr{Message: "Unauthorized"}
		}
	}
	userInfo, err := c.userRepo.GetUserById(ctx, userId, nil)
	if err != nil {
		return nil, err
	}
	return c.courseRepo.GetCoursesByUserId(ctx, userId, userInfo.Role, semester, academicYear, nil)
}

func NewCourseService(courseRepo postgres.CourseRepo, transactionManager repo.TransactionManager, authMiddleware middleware.AuthMiddleware, userRepo postgres.UserRepo) CourseService {
	return &courseService{
		courseRepo:         courseRepo,
		transactionManager: transactionManager,
		authMiddleware:     authMiddleware,
		userRepo:           userRepo,
	}
}
