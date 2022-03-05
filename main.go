package main

import (
	"andre/belajar-golang-restful-api/app"
	"andre/belajar-golang-restful-api/controller"
	"andre/belajar-golang-restful-api/exception"
	"andre/belajar-golang-restful-api/helper"
	"andre/belajar-golang-restful-api/repository"
	"andre/belajar-golang-restful-api/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

func main() {
	db := app.NewDB()
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoyController := controller.NewCategoryController(categoryService)

	router := httprouter.New()

	router.GET("/api/categories", categoyController.FindAll)
	router.GET("/api/categories/:categoryId", categoyController.FindById)
	router.POST("/api/categories", categoyController.Create)
	router.PUT("/api/categories/:categoryId", categoyController.Update)
	router.DELETE("/api/categories/:categoryId", categoyController.Delete)

	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
