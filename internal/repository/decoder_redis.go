package repository

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type DecoderRedis struct {
	rdb *redis.Client
}

func NewDecoderRedis(rdb *redis.Client) *DecoderRedis {
	return &DecoderRedis{rdb: rdb}
}

func (r *DecoderRedis) GetBase63(id int64) (string, error) {
	idString := strconv.FormatInt(id, 10)

	res := r.rdb.Get(context.TODO(), idString)
	if res.Err() != nil {
		if res.Err() == redis.Nil {
			return "", ErrNoSuchURL
		}
		return "", res.Err()
	}

	return res.Val(), nil
}
