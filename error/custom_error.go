package error

import (
	"errors"
	"fmt"
)

var WrongPasswordErr = errors.New("wrong password")
var CourseLimitExceededErr = errors.New("course is full")
var CourseRegisterTimoutErr = errors.New("course is not open for register or unregister")

type ResourceNotFoundErr struct {
	Resource string
}

func (e *ResourceNotFoundErr) Error() string {
	return fmt.Sprintf("%s not found", e.Resource)
}

type UniqueConstraintErr struct {
	Message string
}

func (e *UniqueConstraintErr) Error() string {
	return fmt.Sprintf("Unique constraint violated: %s", e.Message)
}

type UnauthorizedErr struct {
	Message string
}

func (e *UnauthorizedErr) Error() string {
	return fmt.Sprintf("Unauthorized: %s", e.Message)
}

type InvalidInputErr struct {
	Message string
}

func (e *InvalidInputErr) Error() string {
	return fmt.Sprintf("Invalid input: %s", e.Message)
}
