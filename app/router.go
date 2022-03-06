package app

import (
	"andre/belajar-golang-restful-api/controller"
	"andre/belajar-golang-restful-api/exception"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(categoyController controller.CategoryController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/categories", categoyController.FindAll)
	router.GET("/api/categories/:categoryId", categoyController.FindById)
	router.POST("/api/categories", categoyController.Create)
	router.PUT("/api/categories/:categoryId", categoyController.Update)
	router.DELETE("/api/categories/:categoryId", categoyController.Delete)

	router.PanicHandler = exception.ErrorHandler

	return router
}
