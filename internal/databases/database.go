package databases

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Databases struct {
	postgres *pgxpool.Pool
	redis    *redis.Client
}

func (d *Databases) Redis() *redis.Client {
	return d.redis
}

func (d *Databases) PostgresHealth(ctx context.Context) error {
	if d.postgres == nil {
		return errors.New("postgres connection not initialized")
	}
	return d.postgres.Ping(ctx)
}

func (d *Databases) RedisHealth(ctx context.Context) error {
	if d.redis == nil {
		return errors.New("redis connection not initialized")
	}
	return d.redis.Ping(ctx).Err()
}

func (d *Databases) PG() *pgxpool.Pool {
	return d.postgres
}

func (d *Databases) NewPostgresInit(dsn string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("unable to connect to database:", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("db ping failed:", err)
	}

	d.postgres = pool
}

func (d *Databases) NewRedis(addr, password string, db int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr, // "localhost:6379"
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatal("redis connection failed:", err)
	}

	log.Println("Connected to Redis successfully")

	d.redis = rdb
}
