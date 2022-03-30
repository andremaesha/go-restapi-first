package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTryFirst(t *testing.T) {
	name := "hello"

	assert.Equal(t, "hello", name, "name must be 'hello'.")
}
