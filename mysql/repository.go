package mysql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mohit405/config"
	_ "github.com/upper/db/v4/adapter/mysql"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(cfgdb config.Mysql) (*Storage, error) {
	var err error

	s := new(Storage)
	s.db, err = sql.Open(cfgdb.Dialect, cfgdb.DSN)
	if err != nil {
		log.Println("Error while connecting to mysql")
		return s, err
	}

	err = s.db.Ping()
	if err != nil {
		log.Println("connection is not alive")
		return s, err
	}

	return s, nil
}

func (s *Storage) GetStudents() ([]Student, error) {
	var students []Student
	res, err := s.db.Query(`Select * from students`)
	if err != nil {
		return students, err
	}

	for res.Next() {
		var student Student
		if err := res.Scan(&student.StudentId, &student.Name, &student.Stream); err != nil {
			return students, err
		}

		students = append(students, student)
	}

	return students, nil
}

func (s *Storage) GetStudent(id int) (Student, error) {
	var student Student
	q := fmt.Sprintf(`Select * from students where student_id = %d`, id)

	err := s.db.QueryRow(q).Scan(&student.StudentId, &student.Name, &student.Stream)
	if err != nil {
		return student, err
	}

	return student, nil
}
