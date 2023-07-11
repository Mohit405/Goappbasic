package main

import (
	"log"

	"github.com/mohit405/config"
	"github.com/mohit405/mysql"
)


func (app *application) GetStudentData(id int) (mysql.Student, error) {
	data, err := app.rediconn.ReadStudentData(id)
	if err != nil {
		if err.Error() == config.RedisErrKeyDoesNotExist {
			err = nil

			data, err = app.sqlconn.GetStudent(id)
			if err != nil {
				return data, err
			}

			log.Println("Read from Mysql")
			
			//set the data in the redis..
			app.rediconn.SetStudentData(data)
		}
	}

	return data, nil
}

func (app *application) GetStudentsData() ([]mysql.Student, error) {
	data, err := app.rediconn.ReadAllData()
	if err != nil {
		if err.Error() == config.RedisErrKeyDoesNotExist {
			err = nil

			data, err = app.sqlconn.GetStudents()
			if err != nil {
				return data, err
			}
			log.Print("Read from Mysql")

			//set the data in the redis..
			app.rediconn.SetData(data)
		}
	}

	return data, nil
}
