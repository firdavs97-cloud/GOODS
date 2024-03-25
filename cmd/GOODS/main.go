package main

import (
	"goods/api/server"
	"goods/pkg/cache"
	"goods/pkg/postgres/db"
)

func main() {
	// Read config

	// Start services
	db.Connect(5432, "localhost", "myuser", "mysecretpassword", "mydatabase")
	cache.Connect("localhost:6379", "")
	cache.FlushAll()

	// Start app
	server.Run()
}
