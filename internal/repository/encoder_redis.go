package repository

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type EncoderRedis struct {
	rdb *redis.Client
}

func NewEncoderRedis(rdb *redis.Client) *EncoderRedis {
	return &EncoderRedis{rdb: rdb}
}

func (r *EncoderRedis) SaveBase63(url string, id int64) error {
	idString := strconv.FormatInt(id, 10)

	// does the key exist
	exists, err := r.rdb.Exists(context.TODO(), idString).Result()
	if err != nil {
		return err
	}

	// if yes, generate a new one
	if exists == 1 {
		return ErrDuplicateId
	}

	// if not, write to db
	res := r.rdb.Set(context.TODO(), idString, url, 0)
	if res.Err() != nil {
		return res.Err()
	}

	// check if no data race were made
	check := r.rdb.Get(context.TODO(), idString)
	if check.Err() != nil {
		return check.Err()
	}

	// if there was a data race, generate a new id
	if check.Val() != url {
		return ErrDuplicateId
	}

	return nil
}
