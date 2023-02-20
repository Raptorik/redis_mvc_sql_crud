package config

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
)

type MySQLConnection interface {
	GetConnection() (*sql.DB, error)
}

type RedisConnection interface {
	GetConnection() (*redis.Client, error)
}

type mysqlConnection struct {
	Username string
	Password string
	Hostname string
	Port     string
	DBName   string
}

type redisConnection struct {
	Hostname string
	Port     string
	Password string
}

func NewMySQLConnection() MySQLConnection {
	return mysqlConnection{
		Username: "root",
		Password: "password",
		Hostname: "localhost",
		Port:     "3306",
		DBName:   "my-mvc",
	}
}

func (c mysqlConnection) GetConnection() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.Hostname, c.Port, c.DBName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %s", err.Error())
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %s", err.Error())
	}

	return db, nil
}

func NewRedisConnection() RedisConnection {
	return redisConnection{
		Hostname: "localhost",
		Port:     "6379",
		Password: "mypass",
	}
}
func (c redisConnection) GetConnection() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Hostname, c.Port),
		Password: c.Password,
		DB:       1,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("error pinging cache: %s", err.Error())
	}

	return client, nil
}
