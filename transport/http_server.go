package transport

import (
	"SchoolManagement/dto"
	"SchoolManagement/dto/request"
	"SchoolManagement/dto/response"
	"SchoolManagement/endpoint"
	error2 "SchoolManagement/error"
	"SchoolManagement/middleware"
	"SchoolManagement/repo"
	"SchoolManagement/repo/postgres"
	"SchoolManagement/repo/redis"
	"SchoolManagement/service"
	"SchoolManagement/utils"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	redis2 "github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
	"strings"
)

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var notFoundErr *error2.ResourceNotFoundErr
	var uniqueConstraintErr *error2.UniqueConstraintErr
	var unauthorizedErr *error2.UnauthorizedErr
	var invalidInputErr *error2.InvalidInputErr
	var validationError validator.ValidationErrors
	switch {
	case errors.Is(err, error2.WrongPasswordErr):
		w.WriteHeader(http.StatusBadRequest)
	case errors.As(err, &notFoundErr):
		w.WriteHeader(http.StatusNotFound)
	case errors.As(err, &uniqueConstraintErr):
		w.WriteHeader(http.StatusBadRequest)
	case errors.As(err, &unauthorizedErr):
		w.WriteHeader(http.StatusUnauthorized)
	case errors.As(err, &invalidInputErr):
		w.WriteHeader(http.StatusBadRequest)
	case errors.As(err, &validationError):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_ = json.NewEncoder(w).Encode(response.Message{Error: err.Error()})
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req request.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeLoginResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeRegisterStudentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req request.StudentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeRegisterStudentResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeUpdateStudentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	var req request.StudentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	req.Id = id
	return req, nil
}

func encodeUpdateStudentResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeDeleteStudentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	return id, nil
}

func encodeDeleteStudentResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeGetStudentByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	return id, nil
}

func encodeGetStudentByIdResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeRegisterTeacherRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req request.TeacherRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeRegisterTeacherResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeUpdateTeacherRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	var req request.TeacherRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	req.Id = id
	return req, nil
}

func encodeUpdateTeacherResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeDeleteTeacherRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	return id, nil
}

func encodeDeleteTeacherResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeGetTeacherByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	return id, nil
}

func encodeGetTeacherByIdResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeCreateSubjectRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req dto.SubjectDto
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeCreateSubjectResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeUpdateSubjectRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	var req dto.SubjectDto
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	req.Id = id
	return req, nil
}

func encodeUpdateSubjectResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeDeleteSubjectByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	return id, nil
}

func encodeDeleteSubjectResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeGetSubjectByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	return id, nil
}

func encodeGetSubjectByIdResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeGetSubjectListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	major := r.URL.Query().Get("major")
	params := dto.GetSubjectsParamDTO{
		PaginationParams: dto.PaginationParams{},
		Major:            "",
	}
	if major != "" {
		params.Major = major
	}
	if limit != "" {
		var err error
		params.Limit, err = strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}
	}
	if offset != "" {
		var err error
		params.Offset, err = strconv.Atoi(offset)
		if err != nil {
			return nil, err
		}
	}
	return params, nil
}

func encodeGetSubjectListResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeCreateCourseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req request.CourseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeCreateCourseResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

func decodeGetCourseByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	return id, nil
}

func encodeGetCourseByIdResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeUpdateCourseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	var req request.CourseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	req.Id = id
	return req, nil
}

func encodeUpdateCourseResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeDeleteCourseByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	return id, nil
}

func encodeDeleteCourseByIdResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeRegisterStudentToCourseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req request.CourseRegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeRegisterStudentToCourseResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeUnregisterStudentFromCourseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parts := strings.Split(r.URL.Path, "/")
	courseId := parts[len(parts)-3]
	userId := parts[len(parts)-1]
	return request.CourseRegistrationRequest{
		CourseId:  courseId,
		StudentId: userId,
	}, nil
}

func encodeUnregisterStudentFromCourseResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeAddCourseScheduleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req request.CourseScheduleRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeAddCourseScheduleResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeGetCourseScheduleByCourseIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	courseId := r.URL.Query().Get("courseId")
	return courseId, nil
}

func encodeGetCourseScheduleByCourseIdResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeDeleteCourseScheduleByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	return id, nil
}

func encodeDeleteCourseScheduleByIdResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeGetCoursesByUserIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	semester := r.URL.Query().Get("semester")
	s, err := strconv.Atoi(semester)
	if err != nil {
		return nil, err
	}
	academicYear := r.URL.Query().Get("academicYear")
	userId := r.URL.Query().Get("userId")
	return dto.GetCoursesParams{
		UserId:       userId,
		Semester:     s,
		AcademicYear: academicYear,
	}, nil
}

func encodeGetCoursesByUserIdResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func NewHttpServer(db *sqlx.DB, redisClient *redis2.Client) *gin.Engine {
	teacherRepo := postgres.NewTeacherRepo(db)
	userRepo := postgres.NewUserRepo(db)
	studentRepo := postgres.NewStudentRepo(db)
	transactionManager := repo.NewTransactionManager(db)
	subjectRepo := postgres.NewSubjectRepo(db)
	courseRepo := postgres.NewCourseRepo(db)

	teacherCache := redis.NewTeacherCache(redisClient)
	studentCache := redis.NewStudentCache(redisClient)

	jwtUtils := utils.NewJwtUtils()
	authMiddleware := middleware.NewAuthMiddleware(jwtUtils)

	authService := service.NewAuthService(userRepo, jwtUtils)
	teacherService := service.NewTeacherService(userRepo, teacherRepo, transactionManager, teacherCache, authMiddleware)
	studentService := service.NewStudentService(studentRepo, userRepo, transactionManager, studentCache, authMiddleware)
	subjectService := service.NewSubjectService(subjectRepo, authMiddleware)
	courseService := service.NewCourseService(courseRepo, transactionManager, authMiddleware, userRepo)

	authEndpoint := endpoint.NewAuthEndpoint(authService)
	teacherEndpoint := endpoint.NewTeacherEndpoint(teacherService)
	studentEndpoint := endpoint.NewStudentEndpoint(studentService)
	subjectEndpoint := endpoint.NewSubjectEndpoint(subjectService)
	courseEndpoint := endpoint.NewCourseEndpoint(courseService)

	options := []http2.ServerOption{
		http2.ServerErrorEncoder(encodeError),
	}

	loginHandler := http2.NewServer(
		authEndpoint.Login(),
		decodeLoginRequest,
		encodeLoginResponse,
		options...)

	registerStudentHandler := http2.NewServer(
		studentEndpoint.RegisterStudentEndpoint(),
		decodeRegisterStudentRequest,
		encodeRegisterStudentResponse,
		options...)

	updateStudentHandler := http2.NewServer(
		studentEndpoint.UpdateStudentEndpoint(),
		decodeUpdateStudentRequest,
		encodeUpdateStudentResponse,
		options...)

	deleteStudentHandler := http2.NewServer(
		studentEndpoint.DeleteStudentByIdEndpoint(),
		decodeDeleteStudentRequest,
		encodeDeleteStudentResponse,
		options...)

	getStudentByIdHandler := http2.NewServer(
		studentEndpoint.GetStudentByIdEndpoint(),
		decodeGetStudentByIdRequest,
		encodeGetStudentByIdResponse,
		options...)

	registerTeacherHandler := http2.NewServer(
		teacherEndpoint.RegisterTeacherEndpoint(),
		decodeRegisterTeacherRequest,
		encodeRegisterTeacherResponse,
		options...)

	updateTeacherHandler := http2.NewServer(
		teacherEndpoint.UpdateTeacherEndpoint(),
		decodeUpdateTeacherRequest,
		encodeUpdateTeacherResponse,
		options...)

	deleteTeacherHandler := http2.NewServer(
		teacherEndpoint.DeleteTeacherByIdEndpoint(),
		decodeDeleteTeacherRequest,
		encodeDeleteTeacherResponse,
		options...)

	getTeacherByIdHandler := http2.NewServer(
		teacherEndpoint.GetTeacherByIdEndpoint(),
		decodeGetTeacherByIdRequest,
		encodeGetTeacherByIdResponse,
		options...)

	createSubjectHandler := http2.NewServer(
		subjectEndpoint.CreateSubjectEndpoint(),
		decodeCreateSubjectRequest,
		encodeCreateSubjectResponse,
		options...)

	updateSubjectHandler := http2.NewServer(
		subjectEndpoint.UpdateSubjectEndpoint(),
		decodeUpdateSubjectRequest,
		encodeUpdateSubjectResponse,
		options...)

	deleteSubjectByIdHandler := http2.NewServer(
		subjectEndpoint.DeleteSubjectByIdEndpoint(),
		decodeDeleteSubjectByIdRequest,
		encodeDeleteSubjectResponse,
		options...)

	getSubjectByIdHandler := http2.NewServer(
		subjectEndpoint.GetSubjectByIdEndpoint(),
		decodeGetSubjectByIdRequest,
		encodeGetSubjectByIdResponse,
		options...)

	getSubjectListHandler := http2.NewServer(
		subjectEndpoint.GetSubjectListEndpoint(),
		decodeGetSubjectListRequest,
		encodeGetSubjectListResponse,
		options...)

	createCourseHandler := http2.NewServer(
		courseEndpoint.CreateCourse(),
		decodeCreateCourseRequest,
		encodeCreateCourseResponse,
		options...)

	getCourseByIdHandler := http2.NewServer(
		courseEndpoint.GetCourseById(),
		decodeGetCourseByIdRequest,
		encodeGetCourseByIdResponse,
		options...)

	updateCourseHandler := http2.NewServer(
		courseEndpoint.UpdateCourse(),
		decodeUpdateCourseRequest,
		encodeUpdateCourseResponse,
		options...)

	deleteCourseByIdHandler := http2.NewServer(
		courseEndpoint.DeleteCourseById(),
		decodeDeleteCourseByIdRequest,
		encodeDeleteCourseByIdResponse,
		options...)

	registerStudentToCourseHandler := http2.NewServer(
		courseEndpoint.RegisterStudentToCourse(),
		decodeRegisterStudentToCourseRequest,
		encodeRegisterStudentToCourseResponse,
		options...)

	unregisterStudentFromCourseHandler := http2.NewServer(
		courseEndpoint.UnregisterStudentFromCourse(),
		decodeUnregisterStudentFromCourseRequest,
		encodeUnregisterStudentFromCourseResponse,
		options...)

	addCourseScheduleHandler := http2.NewServer(
		courseEndpoint.AddCourseSchedule(),
		decodeAddCourseScheduleRequest,
		encodeAddCourseScheduleResponse,
		options...)

	getCourseSchedulesByCourseIdHandler := http2.NewServer(
		courseEndpoint.GetCourseSchedulesByCourseId(),
		decodeGetCourseScheduleByCourseIdRequest,
		encodeGetCourseScheduleByCourseIdResponse,
		options...)

	deleteCourseScheduleByIdHandler := http2.NewServer(
		courseEndpoint.DeleteCourseScheduleById(),
		decodeDeleteCourseScheduleByIdRequest,
		encodeDeleteCourseScheduleByIdResponse,
		options...)

	getCourseByUserIdHandler := http2.NewServer(
		courseEndpoint.GetCoursesByUserId(),
		decodeGetCoursesByUserIdRequest,
		encodeGetCoursesByUserIdResponse,
		options...)

	r := gin.Default()

	authRoute := r.Group("/auth")
	authRoute.POST("/login", gin.WrapH(loginHandler))

	studentRoute := r.Group("/student")
	studentRoute.POST("/register", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(registerStudentHandler))
	studentRoute.PATCH("/update/:id", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(updateStudentHandler))
	studentRoute.DELETE("/:id", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(deleteStudentHandler))
	studentRoute.GET("/:id", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(getStudentByIdHandler))

	teacherRoute := r.Group("/teacher")
	teacherRoute.POST("/register", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(registerTeacherHandler))
	teacherRoute.PATCH("/update/:id", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(updateTeacherHandler))
	teacherRoute.DELETE("/:id", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(deleteTeacherHandler))
	teacherRoute.GET("/:id", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(getTeacherByIdHandler))

	subjectRoute := r.Group("/subject")
	subjectRoute.POST("/create", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(createSubjectHandler))
	subjectRoute.PATCH("/update/:id", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(updateSubjectHandler))
	subjectRoute.DELETE("/:id", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(deleteSubjectByIdHandler))
	subjectRoute.GET("/:id", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(getSubjectByIdHandler))
	subjectRoute.GET("", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(getSubjectListHandler))

	courseRoute := r.Group("/course")
	courseRoute.POST("/create", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(createCourseHandler))
	courseRoute.GET("/:id", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(getCourseByIdHandler))
	courseRoute.PATCH("/:id", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(updateCourseHandler))
	courseRoute.DELETE("/:id", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(deleteCourseByIdHandler))
	courseRoute.GET("", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(getCourseByUserIdHandler))
	courseRoute.POST("/register", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(registerStudentToCourseHandler))
	courseRoute.DELETE("/unregister/:courseId/student/:studentId", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(unregisterStudentFromCourseHandler))
	courseRoute.POST("/schedule", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(addCourseScheduleHandler))
	courseRoute.GET("/schedule", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(getCourseSchedulesByCourseIdHandler))
	courseRoute.DELETE("/schedule/:id", authMiddleware.ValidateAndExtractJwt(), gin.WrapH(deleteCourseScheduleByIdHandler))
	return r
}
