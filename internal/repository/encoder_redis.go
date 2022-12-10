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

func (r *EncoderRedis) checkTryLimits(counter *int) error {
	if *counter > TRIES_LIMIT {
		numOfKeys := r.rdb.DBSize(context.TODO())
		if numOfKeys.Err() != nil {
			return numOfKeys.Err()
		}
		if numOfKeys.Val() > MAX_NUMBER_OF_ROWS*2 {
			return ErrOutOfUniqueValues
		} else {
			*counter = 0
		}
	}

	return nil
}

func (r *EncoderRedis) SaveBase63(url string, id int64, counter *int) error {
	if err := r.checkTryLimits(counter); err != nil {
		return err
	}

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

	txf := func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(context.TODO(), func(pipe redis.Pipeliner) error {
			pipe.Set(context.TODO(), idString, url, 0)
			return nil
		})
		return err
	}

	// add new id key
	if err := r.rdb.Watch(context.TODO(), txf, idString); err != nil {
		if err == redis.TxFailedErr {
			return ErrDuplicateId
		} else {
			return err
		}
	}

	// add new url key
	res := r.rdb.Set(context.TODO(), url, idString, 0)
	if res.Err() != nil {
		r.rdb.Del(context.TODO(), idString)
		return res.Err()
	}

	return nil
}

func (r *EncoderRedis) IfExistsBase63(url string) (int64, error) {
	result := r.rdb.Get(context.TODO(), url)
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return -1, nil
		}
		return -1, result.Err()
	}

	id, err := strconv.ParseInt(result.Val(), 10, 64)
	if err != nil {
		return -1, err
	}

	return id, nil
}
