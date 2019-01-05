package main

import (
	"os"

	"github.com/cschen13/spotitube/server"
	"github.com/cschen13/spotitube/utils"
	"github.com/go-redis/redis"
)

func main() {
	hostname := os.Getenv("SPOTITUBE_HOST")
	scheme := "https://"
	devPort := getPort() // Development server hosting the React app
	if hostname == "" {
		hostname = "localhost"
		scheme = "http://"
		devPort = ":3000"
	}

	sessionSecret := os.Getenv("SPOTITUBE_SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = utils.GenerateRandStr(64)
	}

	redisURL := os.Getenv("REDIS_URL")
	redisAddress := ":6379"
	redisPassword := ""
	if redisURL != "" {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			panic(err)
		}

		redisAddress = opt.Addr
		redisPassword = opt.Password
	}

	s := server.NewServer(scheme+hostname, getPort(), redisAddress, redisPassword, sessionSecret, 1, devPort)
	s.Start()
}

func getPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	}
	return ":8080"
}
