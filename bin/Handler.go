package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/eben-hk/confide"
	"github.com/gorilla/schema"
	"github.com/julienschmidt/httprouter"
	"github.com/mohit405/mysql"
)

func (app *application) getStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		fmt.Fprintf(w, err.Error(), "String to int conversion Error!")
		http.NotFound(w, r)
		return
	}

	data, err := app.GetStudentData(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)

	if err != nil {
		fmt.Fprint(w, "failed to marshal the data!")
		return
	}
}

func (app *application) getStudents(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	msf := new(mysql.StudentFilter)

	if err := r.ParseForm(); err != nil {
		confide.JSON(w, confide.Payload{Code: confide.FCodeBadRequest, Message: err.Error()})
		return
	}

	if err := schema.NewDecoder().Decode(msf, r.Form); err != nil {
		confide.JSON(w, confide.Payload{Code: confide.FCodeBadRequest, Message: err.Error()})
		return
	}

	data, err := app.GetStudentsData(*msf)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) registerStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ms mysql.Student
	// data, err := ioutil.ReadAll(r.Body)
	err := app.readJSON(w, r, &ms)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// json.Unmarshal(data, &ms)

	err = app.RegisterStudent(ms)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
}

func (app *application) updateStudentInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ms mysql.Student
	data, err := ioutil.ReadAll(r.Body)
	id := ps.ByName("id")
	if err != nil {
		fmt.Fprintf(w, err.Error(), "fail to read the data!")
		return
	}

	json.Unmarshal(data, &ms)
	ms.StudentId, err = strconv.Atoi(id)
	if err != nil {
		fmt.Fprintf(w, err.Error(), "enter a integer value in range")
		return
	}

	if ms.Email_id == "" && ms.Address == "" {
		fmt.Fprintf(w, "add one field atleast!")
		return
	}
	_, err = app.sqlconn.UpdateStudent(ms)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
}
