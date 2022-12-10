package database

import (
	"errors"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	POSTGRES = "postgres"
	REDIS    = "redis"
)

type Database struct {
	DBType string
	Sql    *sqlx.DB
	Redis  *redis.Client
}

func NewDatabase(dbType string) (*Database, error) {
	switch dbType {
	case POSTGRES:
		logrus.Println("Connecting to Postgres Database...")
		return NewPostgresDB(PostgresConfig{
			Host:     viper.GetString("postgres.host"),
			Port:     viper.GetString("postgres.port"),
			Username: viper.GetString("postgres.username"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   viper.GetString("postgres.dbname"),
			SSLMode:  viper.GetString("postgres.sslmode"),
		})
	case REDIS:
		logrus.Println("Connecting to Redis Database...")
		return NewRedisDB(RedisConfig{
			Addr:     viper.GetString("redis.addr"),
			Password: os.Getenv("REDIS_PASSWORD"),
		})
	}

	return nil, errors.New("unknown database type")
}

func (d *Database) Close() error {
	if d.Sql != nil {
		if err := d.Sql.Close(); err != nil {
			return err
		}
	}

	if d.Redis != nil {
		if err := d.Redis.Close(); err != nil {
			return err
		}
	}

	return nil
}
