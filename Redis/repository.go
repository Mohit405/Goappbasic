package Redis

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/mohit405/config"
	"github.com/mohit405/mysql"
)

type Storage struct {
	db *redis.Client
}

func MakeConnection(cfgdb config.Redis) (*Storage, error) {
	var err error

	s := new(Storage)

	opts := &redis.Options{Addr: cfgdb.Addr, DB: cfgdb.Db}

	s.db = redis.NewClient(opts)

	_, err = s.db.Ping(context.Background()).Result()
	if err != nil {
		return s, err
	}

	log.Println("Redis connected", cfgdb.Addr, cfgdb.Db)
	return s, nil
}

func (s *Storage) SetData(data []mysql.Student) {
	st, _ := json.Marshal(data)

	s.db.Set(context.Background(), "Students", st, config.RedisTTLOneHour)
}

func (s *Storage) SetStudentData(data mysql.Student) {
	st, _ := json.Marshal(data)
	s.db.Set(context.Background(), "Student:"+strconv.Itoa(data.StudentId), st, config.RedisTTLOneHour)
	// log.Println("Student:" + strconv.Itoa(data.StudentId))
}

func (s *Storage) ReadAllData() ([]mysql.Student, error) {
	var data []mysql.Student
	res := s.db.Get(context.Background(), "Students")
	json.Unmarshal([]byte(res.Val()), &data)

	if res.Val() == "" {
		return data, errors.New(config.RedisErrKeyDoesNotExist)
	}

	return data, nil
}

func (s *Storage) ReadStudentData(id int) (mysql.Student, error) {
	var data mysql.Student

	res := s.db.Get(context.Background(), "Student:"+strconv.Itoa(id))
	// log.Println("Student:" + strconv.Itoa(id))
	json.Unmarshal([]byte(res.Val()), &data)

	if res.Val() == "" {
		return data, errors.New(config.RedisErrKeyDoesNotExist)
	}

	return data, nil
}
