package test

import (
	"mail-service/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "securepassword"
	hash, err := utils.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestCheckPassword(t *testing.T) {
	password := "securepassword"
	hash, _ := utils.HashPassword(password)

	assert.True(t, utils.CheckPasswordHash(password, hash))
	assert.False(t, utils.CheckPasswordHash("wrongpassword", hash))
}
