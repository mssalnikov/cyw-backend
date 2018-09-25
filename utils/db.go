package utils

import "database/sql"
import "github.com/go-redis/redis"

var (
	DBCon *sql.DB

	RedisCon *redis.Client
)
