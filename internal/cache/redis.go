package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gitlab.gnous.eu/ada/status/internal/models"
)

var errPing = errors.New("connection to redis does not work")

var db *redis.Client

type Config struct {
	Enabled  bool
	Address  string
	Db       int
	User     string
	Password string
}

func (c Config) Connect() *redis.Client {
	db = redis.NewClient(&redis.Options{
		Addr:     c.Address,
		Username: c.User,
		Password: c.Password,
		DB:       c.Db,
	})

	if db == nil {
		logrus.Fatal("failed to connect to redis")
	}

	return db
}

func Ping() error {
	ctx := context.Background()

	status := db.Ping(ctx)
	if status.String() != "ping: PONG" {
		return errPing
	}

	return nil
}

func SetCacheResult(status models.Status) {
	ctx := context.Background()

	db.HSet(ctx, status.Name, "success", status.Success)
	db.HSet(ctx, status.Name, "description", status.Description)
	db.HSet(ctx, status.Name, "error", status.Description)
	db.HSet(ctx, status.Name, "target", status.Target)
}

func GetCacheResult(name string) (models.Status, error) {
	ctx := context.Background()

	var status models.Status
	var err error

	status.Name = name
	status.Success, err = db.HGet(ctx, name, "success").Bool()
	if err != nil {
		return models.Status{}, err
	}

	status.Description = db.HGet(ctx, name, "description").Val()
	status.Error = db.HGet(ctx, name, "error").Val()
	status.Target = db.HGet(ctx, name, "target").Val()

	db.Expire(ctx, name, 1*time.Minute)

	return status, nil
}

func KeyExist(key string) bool {
	ctx := context.Background()

	return db.Exists(ctx, key).Val() == 1
}
