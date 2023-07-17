package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

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

func (s *Storage) CreateStudent(ms Student) (int, error) {
	q := "INSERT INTO students (name, stream, email_id,grade) VALUE (?,?,?,?)"
	row, err := s.db.ExecContext(context.Background(), q, ms.Name, ms.Stream, ms.Email_id, ms.Grade)
	if err != nil {
		return 0, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Storage) GetStudents(sf StudentFilter) ([]Student, error) {
	var students []Student
	res, err := s.db.Query(`Select s.student_id, s.stream, s.name, s.email_id, s.grade, IFNULL(address,'') FROM students s` + generateConditions(sf))
	if err != nil {
		return students, err
	}

	for res.Next() {
		var student Student
		if err := res.Scan(&student.StudentId, &student.Name, &student.Stream, &student.Email_id, &student.Grade, &student.Address); err != nil {
			return students, err
		}

		students = append(students, student)
	}

	return students, nil
}

func (s *Storage) GetStudent(id int) (Student, error) {
	var student Student
	q := fmt.Sprintf(`Select s.student_id, s.name, s.stream, s.email_id, s.grade, IFNULL(address,'') FROM students s where s.student_id = %d`, id)

	res := s.db.QueryRow(q).Scan(&student.StudentId, &student.Name, &student.Stream, &student.Email_id, &student.Grade, &student.Address)
	if res == sql.ErrNoRows {
		return student, sql.ErrNoRows
	}

	return student, nil
}

func generateConditions(i interface{}) string {
	var (
		where  string
		exprs  []string
		sort   string
		sortby string
		page   int
		offset int
	)

	v := reflect.ValueOf(i)
	t := reflect.TypeOf(i)
	sortorder := "ASC"
	limit := 10

	for i := 0; i < v.NumField(); i++ {
		switch t.Field(i).Name {
		case "SortBy":
			if v.Field(i).Interface() != "" {
				if v.Field(i).Interface() == "grade" {
					sortby = "grade"
				}
			}
		case "IsAsc":
			if v.Field(i).Interface() == false {
				sortorder = "DESC"
			}
		case "Limit":
			if v.Field(i).Interface().(int) > 0 {
				page = v.Field(i).Interface().(int)
			}
		}
	}

	if exprs != nil {
		where = fmt.Sprintf(" WHERE %s", strings.Join(exprs, " AND "))
	}

	if sortby != "" {
		sort = fmt.Sprintf(" ORDER BY %s %s", sortby, sortorder)
	}

	if page > 0 {
		offset = (page - 1) * limit
	}

	return fmt.Sprintf("%s%s LIMIT %d OFFSET %d", where, sort, limit, offset)
}

func (s *Storage) UpdateStudent(ms Student) (int, error) {
	var q string
	if ms.Email_id == "" {
		q = fmt.Sprintf("UPDATE students SET address = '%s' WHERE student_id = %d", ms.Address, ms.StudentId)
	} else if ms.Address == "" {
		q = fmt.Sprintf("UPDATE students SET email_id = '%s' WHERE student_id = %d", ms.Email_id, ms.StudentId)

	} else {
		q = fmt.Sprintf("UPDATE students SET address = '%s', email_id = '%s' WHERE student_id = %d", ms.Address, ms.Email_id, ms.StudentId)
	}
	log.Print(q)
	row, err := s.db.Exec(q)
	if err != nil {
		return 0, err
	}

	rows, err := row.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rows), nil
}
