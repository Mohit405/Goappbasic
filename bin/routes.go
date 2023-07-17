package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodAllowedResponse)

	router.GET("/v1/student/:id", middleware(app.getStudent))
	router.GET("/v1/students", middleware(app.getStudents))
	router.POST("/v1/registerStudent", middleware(app.registerStudent))
	router.PATCH("/v1/updateStudentinfo/:id", middleware(app.updateStudentInfo))

	return router
}
