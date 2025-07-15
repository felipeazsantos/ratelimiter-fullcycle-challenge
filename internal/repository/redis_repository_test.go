package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestAllowToken(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewRedisRepository(db)

	key := "abc123"
	maxRequest := 2
	blockTime := "1s"

	mock.ExpectTxPipeline()
	mock.ExpectIncr(key).SetVal(1)
	mock.ExpectExpire(key, time.Second).SetVal(true)
	mock.ExpectTxPipelineExec()

	allowed, err := repo.AllowToken(context.Background(), key, maxRequest, blockTime)
	assert.NoError(t, err)
	assert.True(t, allowed)
}

func TestAllowIP(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewRedisRepository(db)

	key := "1.2.3.4"
	maxRequest := 2
	blockTime := "1s"

	mock.ExpectTxPipeline()
	mock.ExpectIncr(key).SetVal(1)
	mock.ExpectExpire(key, time.Second).SetVal(true)
	mock.ExpectTxPipelineExec()

	allowed, err := repo.AllowIP(context.Background(), key, maxRequest, blockTime)
	assert.NoError(t, err)
	assert.True(t, allowed)
}

func TestAllowToken_ErrorParseDuration(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewRedisRepository(db)

	key := "abc123"
	maxRequest := 2
	blockTime := "xx"

	mock.ExpectTxPipeline()
	mock.ExpectIncr(key).SetVal(1)

	allowed, err := repo.AllowToken(context.Background(), key, maxRequest, blockTime)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "time: invalid duration")
	assert.False(t, allowed)
}

func TestAllowIP_ErrorParseDuration(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewRedisRepository(db)

	key := "1.2.3.4"
	maxRequest := 2
	blockTime := "xx"

	mock.ExpectTxPipeline()
	mock.ExpectIncr(key).SetVal(1)

	allowed, err := repo.AllowIP(context.Background(), key, maxRequest, blockTime)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "time: invalid duration")
	assert.False(t, allowed)
}

func TestAllowToken_ErrorPipelineExec(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewRedisRepository(db)

	key := "abc123"
	maxRequest := 2
	blockTime := "1s"

	mock.ExpectTxPipeline()
	mock.ExpectIncr(key).SetVal(1)
	mock.ExpectExpire(key, time.Second).SetVal(true)
	mock.ExpectTxPipelineExec().SetErr(errors.New("some redis error"))

	allowed, err := repo.AllowToken(context.Background(), key, maxRequest, blockTime)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "some redis error")
	assert.False(t, allowed)
}

func TestAllowIP_ErrorPipelineExec(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewRedisRepository(db)

	key := "1.2.3.4"
	maxRequest := 2
	blockTime := "1s"

	mock.ExpectTxPipeline()
	mock.ExpectIncr(key).SetVal(1)
	mock.ExpectExpire(key, time.Second).SetVal(true)
	mock.ExpectTxPipelineExec().SetErr(errors.New("some redis error"))

	allowed, err := repo.AllowIP(context.Background(), key, maxRequest, blockTime)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "some redis error")
	assert.False(t, allowed)
}
