// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"andre/belajar-golang-restful-api/app"
	"andre/belajar-golang-restful-api/controller"
	"andre/belajar-golang-restful-api/middleware"
	"andre/belajar-golang-restful-api/repository"
	"andre/belajar-golang-restful-api/service"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"net/http"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

// Injectors from injector.go:

func InitializedServer() *http.Server {
	categoryRepositoryImplement := repository.NewCategoryRepository()
	db := app.NewDB()
	validate := validator.New()
	categoryServiceImplement := service.NewCategoryService(categoryRepositoryImplement, db, validate)
	categoryControllerImplement := controller.NewCategoryController(categoryServiceImplement)
	router := app.NewRouter(categoryControllerImplement)
	authMiddleware := middleware.NewAuthMiddleware(router)
	server := NewServer(authMiddleware)
	return server
}

// injector.go:

var categorySet = wire.NewSet(repository.NewCategoryRepository, wire.Bind(new(repository.CategoryRepository), new(*repository.CategoryRepositoryImplement)), service.NewCategoryService, wire.Bind(new(service.CategoryService), new(*service.CategoryServiceImplement)), controller.NewCategoryController, wire.Bind(new(controller.CategoryController), new(*controller.CategoryControllerImplement)))
