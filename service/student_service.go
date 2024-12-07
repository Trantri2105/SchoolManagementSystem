package service

import (
	error2 "SchoolManagement/error"
	"SchoolManagement/middleware"
	"SchoolManagement/model"
	"SchoolManagement/repo"
	"SchoolManagement/repo/postgres"
	"SchoolManagement/repo/redis"
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"reflect"
)

type StudentService interface {
	GetStudentById(ctx context.Context, id string) (model.Student, error)
	UpdateStudent(ctx context.Context, student model.Student) error
	CreateStudent(ctx context.Context, student model.Student) error
	DeleteStudentById(ctx context.Context, id string) error
}

type studentService struct {
	studentRepo        postgres.StudentRepo
	userRepo           postgres.UserRepo
	transactionManager repo.TransactionManager
	studentCache       redis.StudentCache
	authMiddleware     middleware.AuthMiddleware
}

func (s *studentService) GetStudentById(ctx context.Context, id string) (model.Student, error) {
	err := s.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin, model.RoleStudent)
	if err != nil {
		return model.Student{}, &error2.UnauthorizedErr{
			Message: "Required admin role or student to get student info",
		}
	}
	claims := ctx.Value(middleware.JWTClaimsContextKey).(jwt.MapClaims)
	if claims["role"].(string) == model.RoleStudent {
		if claims["userId"].(string) != id {
			return model.Student{}, &error2.UnauthorizedErr{Message: "Unauthorized"}
		}
	}
	student, err := s.studentCache.GetStudentById(ctx, id)
	if err == nil {
		return student, nil
	}
	student, err = s.studentRepo.GetStudentById(ctx, id, nil)
	if err != nil {
		return model.Student{}, err
	}
	s.studentCache.SaveStudent(ctx, student)
	return student, nil
}

func (s *studentService) UpdateStudent(ctx context.Context, student model.Student) error {
	err := s.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin, model.RoleStudent)
	if err != nil {
		return &error2.UnauthorizedErr{
			Message: "Required admin role or student to update",
		}
	}
	claims := ctx.Value(middleware.JWTClaimsContextKey).(jwt.MapClaims)
	if claims["role"].(string) == model.RoleStudent {
		if claims["userId"].(string) != student.Id {
			return &error2.UnauthorizedErr{Message: "Unauthorized"}
		}
	}
	if !reflect.ValueOf(student.Password).IsZero() {
		hash, e := bcrypt.GenerateFromPassword([]byte(student.Password), bcrypt.DefaultCost)
		if e != nil {
			log.Println("Student service, update student err :", err)
			return e
		}
		student.Password = string(hash)
	}

	err = s.transactionManager.ExecTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		e := s.userRepo.UpdateUser(ctx, student.User, tx)
		if e != nil {
			return e
		}
		e = s.studentRepo.UpdateStudent(ctx, student, tx)
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		return err
	}
	s.studentCache.DeleteStudentById(ctx, student.Id)
	return nil
}

func (s *studentService) CreateStudent(ctx context.Context, student model.Student) error {
	err := s.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin)
	if err != nil {
		return &error2.UnauthorizedErr{
			Message: "Required admin role to register",
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(student.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Student service, create student err :", err)
		return err
	}
	student.Password = string(hash)

	err = s.transactionManager.ExecTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		e := s.userRepo.InsertUser(ctx, student.User, tx)
		if e != nil {
			return e
		}
		e = s.studentRepo.InsertStudent(ctx, student, tx)
		if e != nil {
			return e
		}
		return nil
	})
	return err
}

func (s *studentService) DeleteStudentById(ctx context.Context, id string) error {
	err := s.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin)
	if err != nil {
		return &error2.UnauthorizedErr{
			Message: "Required admin role to delete",
		}
	}
	err = s.transactionManager.ExecTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		e := s.studentRepo.DeleteStudentById(ctx, id, tx)
		if e != nil {
			return e
		}
		e = s.userRepo.DeleteUserById(ctx, id, tx)
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		return err
	}
	s.studentCache.DeleteStudentById(ctx, id)
	return nil
}

func NewStudentService(studentRepo postgres.StudentRepo, userRepo postgres.UserRepo, transactionManager repo.TransactionManager, studentCache redis.StudentCache, authMiddleware middleware.AuthMiddleware) StudentService {
	return &studentService{
		studentRepo:        studentRepo,
		userRepo:           userRepo,
		transactionManager: transactionManager,
		studentCache:       studentCache,
		authMiddleware:     authMiddleware,
	}
}
