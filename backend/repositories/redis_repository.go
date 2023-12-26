// MIT License

// Copyright (c) The RAI Authors

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package repositories

import (
	"context"
	"time"

	"github.com/retail-ai-inc/bean/dbdrivers"
	"github.com/retail-ai-inc/bean/trace"
)

type RedisRepository interface {
	Keys(c context.Context, tenantID uint64, key string) ([]string, error)
	HSet(c context.Context, tenantID uint64, key string, field string, data interface{}, ttl time.Duration) error
	HGet(c context.Context, tenantID uint64, key string, field string) (string, error)
	HGets(c context.Context, tenantID uint64, keysWithFields map[string]string) (map[string]string, error)
}

type redisRepository struct {
	clients     map[uint64]*dbdrivers.RedisDBConn
	cachePrefix string
}

func NewRedisRepository(clients map[uint64]*dbdrivers.RedisDBConn, cachePrefix string) *redisRepository {
	return &redisRepository{clients, cachePrefix}
}

func (r *redisRepository) Keys(c context.Context, tenantID uint64, key string) ([]string, error) {
	finish := trace.Start(c, "db")
	defer finish()

	prefixKey := r.cachePrefix + "_" + key
	return dbdrivers.RedisGetKeys(c, r.clients[tenantID], prefixKey)
}

func (r *redisRepository) HGet(c context.Context, tenantID uint64, key, field string) (string, error) {
	finish := trace.Start(c, "db")
	defer finish()

	prefixKey := r.cachePrefix + "_" + key
	return dbdrivers.RedisHGet(c, r.clients[tenantID], prefixKey, field)
}

func (r *redisRepository) HGets(c context.Context, tenantID uint64, keysWithFields map[string]string) (map[string]string, error) {
	finish := trace.Start(c, "db")
	defer finish()

	var mappedKeyFieldValues = make(map[string]string)

	for key, field := range keysWithFields {
		prefixKey := r.cachePrefix + "_" + key
		mappedKeyFieldValues[prefixKey] = field
	}

	return dbdrivers.RedisHgets(c, r.clients[tenantID], mappedKeyFieldValues)
}

func (r *redisRepository) HSet(c context.Context, tenantID uint64, key string, field string, data interface{}, ttl time.Duration) error {
	finish := trace.Start(c, "db")
	defer finish()

	prefixKey := r.cachePrefix + "_" + key
	return dbdrivers.RedisHSet(c, r.clients[tenantID], prefixKey, field, data, ttl)
}
