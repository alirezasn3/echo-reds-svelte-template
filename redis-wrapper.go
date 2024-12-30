package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type RedisWrapper struct {
	client *redis.Client
}

func CreateRedisClient(url string) (*RedisWrapper, error) {
	// create redis client
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(options)

	// check redis connection
	if client.Ping(context.Background()).Val() == "PONG" {
		return &RedisWrapper{client: client}, nil
	} else {
		return nil, errors.New("failed to connect to redis, did not receive PONG response")
	}
}

func (rw *RedisWrapper) KeyExists(key string) (bool, error) {
	exists, err := rw.client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if key exists in redis: %s", err.Error())
	}
	return exists == 1, nil
}

func (rw *RedisWrapper) CreateUser(user User) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %s", err.Error())
	}
	user.ID = id.String()
	user.Password, err = HashPassword(user.Password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %s", err.Error())
	}
	err = rw.client.HSet(context.Background(), user.PhoneNumber, user).Err()
	if err != nil {
		return "", fmt.Errorf("failed to add user [%s] to redis: %s", user.PhoneNumber, err.Error())
	}
	return id.String(), nil
}

func (rw *RedisWrapper) FindUser(PhoneNumber string) (*User, error) {
	if ok, err := rw.KeyExists(PhoneNumber); err != nil || !ok {
		return nil, fmt.Errorf("user not found")
	}
	var user User
	err := rw.client.HGetAll(context.Background(), PhoneNumber).Scan(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %s", err.Error())
	}
	return &user, nil
}
