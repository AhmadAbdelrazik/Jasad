package main

import "log"

func main() {
	db, err := NewMySQLServer()
	if err != nil {
		log.Fatal(err)
	}
	

	server := NewAPIServer(":3000", db)
	server.InfoLog.Println("Server is listening on port 3000")
	server.Run()
}