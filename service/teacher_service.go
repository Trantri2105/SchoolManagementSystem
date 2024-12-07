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

type TeacherService interface {
	GetTeacherById(ctx context.Context, id string) (model.Teacher, error)
	UpdateTeacher(ctx context.Context, teacher model.Teacher) error
	CreateTeacher(ctx context.Context, teacher model.Teacher) error
	DeleteTeacherById(ctx context.Context, id string) error
}

type teacherService struct {
	userRepo           postgres.UserRepo
	teacherRepo        postgres.TeacherRepo
	transactionManager repo.TransactionManager
	teacherCache       redis.TeacherCache
	authMiddleware     middleware.AuthMiddleware
}

func (t *teacherService) GetTeacherById(ctx context.Context, id string) (model.Teacher, error) {
	err := t.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin, model.RoleTeacher)
	if err != nil {
		return model.Teacher{}, &error2.UnauthorizedErr{
			Message: "Required admin or teacher role to get teacher info",
		}
	}
	claims := ctx.Value(middleware.JWTClaimsContextKey).(jwt.MapClaims)
	if claims["role"].(string) == model.RoleTeacher {
		if claims["userId"].(string) != id {
			return model.Teacher{}, &error2.UnauthorizedErr{Message: "Unauthorized"}
		}
	}
	teacher, err := t.teacherCache.GetTeacherInfoById(ctx, id)
	if err == nil {
		return teacher, nil
	}
	teacher, err = t.teacherRepo.GetTeacherById(ctx, id, nil)
	if err != nil {
		return model.Teacher{}, err
	}
	t.teacherCache.SaveTeacherInfo(ctx, teacher)
	return teacher, nil
}

func (t *teacherService) UpdateTeacher(ctx context.Context, teacher model.Teacher) error {
	if !reflect.ValueOf(teacher.Role).IsZero() {
		if teacher.Role != model.RoleTeacher && teacher.Role != model.RoleAdmin {
			return &error2.InvalidInputErr{Message: "Role must be teacher or admin"}
		}
		if teacher.Role == model.RoleAdmin {
			err := t.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin)
			if err != nil {
				return &error2.UnauthorizedErr{
					Message: "Required admin role to update this teacher to admin",
				}
			}
		}
	} else {
		err := t.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin, model.RoleTeacher)
		if err != nil {
			return &error2.UnauthorizedErr{
				Message: "Required admin or teacher role to update",
			}
		}
	}
	claims := ctx.Value(middleware.JWTClaimsContextKey).(jwt.MapClaims)
	if claims["role"].(string) == model.RoleTeacher {
		if claims["userId"].(string) != teacher.Id {
			return &error2.UnauthorizedErr{Message: "Unauthorized"}
		}
	}
	if !reflect.ValueOf(teacher.Password).IsZero() {
		hash, err := bcrypt.GenerateFromPassword([]byte(teacher.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Teacher service, update teacher err :", err)
			return err
		}
		teacher.Password = string(hash)
	}

	err := t.transactionManager.ExecTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		e := t.userRepo.UpdateUser(ctx, teacher.User, tx)
		if e != nil {
			return e
		}
		e = t.teacherRepo.UpdateTeacher(ctx, teacher, tx)
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		return err
	}
	t.teacherCache.DeleteTeacherInfoById(ctx, teacher.Id)
	return nil
}

func (t *teacherService) CreateTeacher(ctx context.Context, teacher model.Teacher) error {
	err := t.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin)
	if err != nil {
		return &error2.UnauthorizedErr{
			Message: "Required admin role to register",
		}
	}
	if teacher.Role != model.RoleTeacher && teacher.Role != model.RoleAdmin {
		return &error2.InvalidInputErr{Message: "Role must be teacher or admin"}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(teacher.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Teacher service, create teacher err :", err)
		return err
	}
	teacher.Password = string(hash)

	err = t.transactionManager.ExecTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		e := t.userRepo.InsertUser(ctx, teacher.User, tx)
		if e != nil {
			return e
		}
		e = t.teacherRepo.InsertTeacher(ctx, teacher, tx)
		if e != nil {
			return e
		}
		return nil
	})

	return err
}

func (t *teacherService) DeleteTeacherById(ctx context.Context, id string) error {
	err := t.authMiddleware.CheckUserAuthorities(ctx, model.RoleAdmin)
	if err != nil {
		return &error2.UnauthorizedErr{
			Message: "Required admin role to delete",
		}
	}
	err = t.transactionManager.ExecTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		e := t.teacherRepo.DeleteTeacherById(ctx, id, tx)
		if e != nil {
			return e
		}
		e = t.userRepo.DeleteUserById(ctx, id, tx)
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		return err
	}
	t.teacherCache.DeleteTeacherInfoById(ctx, id)
	return nil
}

func NewTeacherService(userRepo postgres.UserRepo, teacherRepo postgres.TeacherRepo, transactionManager repo.TransactionManager, teacherCache redis.TeacherCache, authMiddleware middleware.AuthMiddleware) TeacherService {
	return &teacherService{
		userRepo:           userRepo,
		teacherRepo:        teacherRepo,
		transactionManager: transactionManager,
		teacherCache:       teacherCache,
		authMiddleware:     authMiddleware,
	}
}
