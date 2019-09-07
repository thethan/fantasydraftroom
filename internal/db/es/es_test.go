package es

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestNewESRepository(t *testing.T) {

	plogger := log.NewNopLogger()

	repo, err := NewESRepository(plogger)

	assert.Nil(t, err)

	newErr := repo.NewIndex(context.TODO(), "ethan")
	assert.Nil(t, newErr)

}
