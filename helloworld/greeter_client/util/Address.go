package util

import (
	"log"
	"os"
)

func GetGrcpAddress() string {
	address := os.Getenv("SERVER_HOST")
	if address == "" {
		address = "localhost:50052"
	}
	log.Println(address)
	return address
}

