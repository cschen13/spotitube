package main

import (
	"github.com/cschen13/spotitube/server"
	"github.com/cschen13/spotitube/utils"
	"os"
)

func main() {
	hostname := os.Getenv("SPOTITUBE_HOST")
	scheme := "https://"
	isDev := false
	if hostname == "" {
		hostname = "localhost"
		scheme = "http://"
		isDev = true
	}

	s := server.NewServer(scheme+hostname, getPort(), utils.GenerateRandStr(64), isDev)
	s.Start()
}

func getPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	}
	return ":8080"
}
