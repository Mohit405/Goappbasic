package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) getStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		log.Print("String to int conversion Error!")
		http.NotFound(w, r)
		return
	}

	data, err := app.GetStudentData(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data1, _ := json.Marshal(data)

	w.Write(data1)
}

func (app *application) getStudents(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	data, err := app.GetStudentsData()
	if err != nil {
		log.Print("no data")
		http.NotFound(w, r)
		return
	}

	data1, _ := json.Marshal(data)

	w.Write(data1)
}
