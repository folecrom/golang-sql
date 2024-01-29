package main

import (
	"github.com/redis/go-redis/v9"
)

func newRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Adresse du serveur Redis
		Password: "",               // Pas de mot de passe
		DB:       0,                // Base de données par défaut
	})
	return client
}
