package test

import (
	"andre/belajar-golang-restful-api/simple"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleServiceError(t *testing.T) {
	simpleService, err := simple.InitializedService(true)

	assert.NotNil(t, err)
	assert.Nil(t, simpleService)
}

func TestSimpleServiceSuccess(t *testing.T) {
	simpleService, err := simple.InitializedService(false)

	assert.Nil(t, err)
	assert.NotNil(t, simpleService)
}

func TestDatabase(t *testing.T) {
	database := simple.InitializedDatabaseRepository()

	assert.Equal(t, "MONGODB", database.DatabaseMongoDB.Name)
	assert.Equal(t, "POSTGRESQL", database.DatabasePostgreSQL.Name)
}

func TestBindingInterface(t *testing.T) {
	data := simple.InitializedHelloService()

	andre := data.Hello("andre")

	assert.Equal(t, "Hello andre", andre)
}

func TestPersonBindingValue(t *testing.T) {
	andre := simple.InitializedPersonUsingValue()

	fmt.Println(andre.Name)
}

// struct field provider
func TestStructFieldProvider(t *testing.T) {
	data := simple.InitializedConfiguration()

	fmt.Println(&data.Name)
}
