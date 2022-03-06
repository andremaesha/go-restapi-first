package test

import (
	"andre/belajar-golang-restful-api/app"
	"andre/belajar-golang-restful-api/controller"
	"andre/belajar-golang-restful-api/helper"
	"andre/belajar-golang-restful-api/middleware"
	"andre/belajar-golang-restful-api/model/domain"
	"andre/belajar-golang-restful-api/repository"
	"andre/belajar-golang-restful-api/service"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root:andre@tcp(localhost:3306)/belajar_golang_restful_api_test")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)
	router := app.NewRouter(categoryController)

	return middleware.NewAuthMiddleware(router)
}

func truncateCategory(db *sql.DB) {
	db.Exec("TRUNCATE category")
}

func TestCreateCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	reqBody := strings.NewReader(`{"name": "Gadget"}`)
	req := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", reqBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()

	assert.Equal(t, 200, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	var resBody map[string]interface{}
	json.Unmarshal(body, &resBody)

	assert.Equal(t, 200, int(resBody["code"].(float64)))
	assert.Equal(t, "OK", resBody["status"])
	assert.Equal(t, "Gadget", resBody["data"].(map[string]interface{})["name"])
}

func TestCreateCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	reqBody := strings.NewReader(`{"name": ""}`)
	req := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", reqBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()
	assert.Equal(t, 400, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	var resBody map[string]interface{}
	json.Unmarshal(body, &resBody)

	assert.Equal(t, 400, int(resBody["code"].(float64)))
	assert.Equal(t, "bad request", resBody["status"])
}

func TestUpdateCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gedget",
	})
	tx.Commit()

	router := setupRouter(db)

	reqBody := strings.NewReader(`{"name": "change"}`)
	req := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), reqBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()

	body, _ := io.ReadAll(res.Body)
	var resBody map[string]interface{}
	json.Unmarshal(body, &resBody)

	assert.Equal(t, 200, int(resBody["code"].(float64)))
	assert.Equal(t, "OK", resBody["status"])
	assert.Equal(t, 1, int(resBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "change", resBody["data"].(map[string]interface{})["name"])
}

func TestUpdateCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gedget",
	})
	tx.Commit()

	router := setupRouter(db)

	reqBody := strings.NewReader(`{"name": ""}`)
	req := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), reqBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()

	body, _ := io.ReadAll(res.Body)
	var resBody map[string]interface{}
	json.Unmarshal(body, &resBody)

	assert.Equal(t, 400, int(resBody["code"].(float64)))
	assert.Equal(t, "bad request", resBody["status"])
}

func TestGetCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gedget",
	})
	tx.Commit()

	router := setupRouter(db)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	req.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()

	body, _ := io.ReadAll(res.Body)
	var resBody map[string]interface{}
	json.Unmarshal(body, &resBody)

	assert.Equal(t, 200, int(resBody["code"].(float64)))
	assert.Equal(t, "OK", resBody["status"])
	assert.Equal(t, category.Id, int(resBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, category.Name, resBody["data"].(map[string]interface{})["name"])
}

func TestGetCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	router := setupRouter(db)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/404", nil)
	req.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()

	body, _ := io.ReadAll(res.Body)
	var resBody map[string]interface{}
	json.Unmarshal(body, &resBody)

	assert.Equal(t, 404, int(resBody["code"].(float64)))
	assert.Equal(t, "not found", resBody["status"])
}

func TestDeleteCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gedget",
	})
	tx.Commit()

	router := setupRouter(db)

	req := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()

	body, _ := io.ReadAll(res.Body)
	var resBody map[string]interface{}
	json.Unmarshal(body, &resBody)

	assert.Equal(t, 200, int(resBody["code"].(float64)))
	assert.Equal(t, "OK", resBody["status"])
}

func TestDeleteCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	router := setupRouter(db)

	req := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/404", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()

	body, _ := io.ReadAll(res.Body)
	var resBody map[string]interface{}
	json.Unmarshal(body, &resBody)

	assert.Equal(t, 404, int(resBody["code"].(float64)))
	assert.Equal(t, "not found", resBody["status"])
}

func TestListCategoriesSuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category1 := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gedget",
	})
	category2 := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Komputer",
	})
	tx.Commit()

	router := setupRouter(db)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)
	req.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()

	body, _ := io.ReadAll(res.Body)
	var resBody map[string]interface{}
	json.Unmarshal(body, &resBody)

	fmt.Println(resBody)

	var categories = resBody["data"].([]interface{})

	categoryResponse1 := categories[0].(map[string]interface{})
	categoryResponse2 := categories[1].(map[string]interface{})

	assert.Equal(t, category1.Id, int(categoryResponse1["id"].(float64)))
	assert.Equal(t, category1.Name, categoryResponse1["name"])

	assert.Equal(t, category2.Id, int(categoryResponse2["id"].(float64)))
	assert.Equal(t, category2.Name, categoryResponse2["name"])
}

func TestUnauthorized(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	router := setupRouter(db)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)
	req.Header.Add("X-API-Key", "SALAH")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()

	body, _ := io.ReadAll(res.Body)
	var resBody map[string]interface{}
	json.Unmarshal(body, &resBody)

	assert.Equal(t, 401, int(resBody["code"].(float64)))
	assert.Equal(t, "unauthorized", resBody["status"])
}
