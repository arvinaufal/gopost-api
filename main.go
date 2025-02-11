package main

import (
	"log"
	"gopost-be/config"
	"gopost-be/routes"
)

func main()  {
	config.ConnectSupabaseDB()

	r := routes.SetupRouter()

	log.Println("Server running on port 8080")
	r.Run(":8080")
}