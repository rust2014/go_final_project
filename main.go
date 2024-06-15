package main

import (
	"github.com/rust2014/go_final_project/database"
	"github.com/rust2014/go_final_project/server"
)

func main() {
	server.Run()
	database.ConnectDB()
}
