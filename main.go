package main

import (
	"andre/belajar-golang-restful-api/app"
	"andre/belajar-golang-restful-api/controller"
	"andre/belajar-golang-restful-api/helper"
	"andre/belajar-golang-restful-api/middleware"
	"andre/belajar-golang-restful-api/repository"
	"andre/belajar-golang-restful-api/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := app.NewDB()
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoyController := controller.NewCategoryController(categoryService)

	router := app.NewRouter(categoyController)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: middleware.NewAuthMiddleware(router),
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
