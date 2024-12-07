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

type TeacherCache interface {
	GetTeacherInfoById(ctx context.Context, id string) (model.Teacher, error)
	SaveTeacherInfo(ctx context.Context, teacher model.Teacher)
	DeleteTeacherInfoById(ctx context.Context, id string)
}

type teacherCache struct {
	cacheTime   time.Duration
	redisClient *redis.Client
}

func getTeacherKey(id string) string {
	return fmt.Sprintf("teacher#%s", id)
}

func (t *teacherCache) GetTeacherInfoById(ctx context.Context, id string) (model.Teacher, error) {
	teacherKey := getTeacherKey(id)
	res, err := t.redisClient.Get(ctx, teacherKey).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Println("Teacher cache err, get teacher info :", err)
		}
		return model.Teacher{}, err
	}
	var teacher model.Teacher
	err = json.Unmarshal([]byte(res), &teacher)
	if err != nil {
		log.Println("Teacher cache err, get teacher info :", err)
		return model.Teacher{}, err
	}
	return teacher, nil
}

func (t *teacherCache) SaveTeacherInfo(ctx context.Context, teacher model.Teacher) {
	teacherKey := getTeacherKey(teacher.Id)
	teacherBytes, err := json.Marshal(teacher)
	if err != nil {
		log.Println("Teacher cache err, save teacher info :", err)
	} else {
		_, err = t.redisClient.Set(ctx, teacherKey, string(teacherBytes), t.cacheTime).Result()
		if err != nil {
			log.Println("Teacher cache err, save teacher info :", err)
		}
	}

}

func (t *teacherCache) DeleteTeacherInfoById(ctx context.Context, id string) {
	teacherKey := getTeacherKey(id)
	_, err := t.redisClient.Del(ctx, teacherKey).Result()
	if err != nil {
		log.Println("Teacher cache err, delete teacher info :", err)
	}
}

func NewTeacherCache(redisClient *redis.Client) TeacherCache {
	return &teacherCache{
		cacheTime:   15 * time.Minute,
		redisClient: redisClient,
	}
}
