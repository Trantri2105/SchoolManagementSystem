package redis

import (
	"SchoolManagement/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type StudentCache interface {
	GetStudentById(ctx context.Context, id string) (model.Student, error)
	DeleteStudentById(ctx context.Context, id string)
	SaveStudent(ctx context.Context, student model.Student)
}

type studentCache struct {
	cacheTime   time.Duration
	redisClient *redis.Client
}

func getStudentKey(id string) string {
	return fmt.Sprintf("student#%s", id)
}

func (s *studentCache) GetStudentById(ctx context.Context, id string) (model.Student, error) {
	studentKey := getStudentKey(id)
	res, err := s.redisClient.Get(ctx, studentKey).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Println("Student cache, get student err :", err)
		}
		return model.Student{}, err
	}
	var student model.Student
	err = json.Unmarshal([]byte(res), &student)
	if err != nil {
		log.Println("Student cache, get student err :", err)
		return model.Student{}, err
	}
	return student, nil
}

func (s *studentCache) DeleteStudentById(ctx context.Context, id string) {
	studentKey := getStudentKey(id)
	_, err := s.redisClient.Del(ctx, studentKey).Result()
	if err != nil {
		log.Println("Student cache, delete student err :", err)
	}
}

func (s *studentCache) SaveStudent(ctx context.Context, student model.Student) {
	studentKey := getStudentKey(student.Id)
	studentBytes, err := json.Marshal(student)
	if err != nil {
		log.Println("Student cache, save student err :", err)
	} else {
		_, err = s.redisClient.Set(ctx, studentKey, studentBytes, s.cacheTime).Result()
		if err != nil {
			log.Println("Student cache, save student err :", err)
		}
	}
}

func NewStudentCache(redisClient *redis.Client) StudentCache {
	return &studentCache{cacheTime: 15 * time.Minute, redisClient: redisClient}
}
