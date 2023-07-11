package main

import (
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	router.GET("/v1/student/:id", middleware(app.getStudent))
	router.GET("/v1/students", middleware(app.getStudents))

	return router
}
